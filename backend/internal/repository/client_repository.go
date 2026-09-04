package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type ClientRepository struct {
	db *pgxpool.Pool
}

func NewClientRepository(db *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{db: db}
}

const clientColumns = "id, name, phone, phone_normalized, email, created_at, updated_at"

func scanClient(row pgx.Row) (model.Client, error) {
	var item model.Client
	err := row.Scan(&item.ID, &item.Name, &item.Phone, &item.PhoneNormalized, &item.Email, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

// FindByPhoneOrEmail matches on normalized phone first, then email
// (case-insensitive), returning apperr.ErrNotFound when neither matches —
// this is the dedup check LeadService.Create runs on every new submission.
func (r *ClientRepository) FindByPhoneOrEmail(ctx context.Context, phoneNormalized, email string) (model.Client, error) {
	row := r.db.QueryRow(ctx, `
		SELECT `+clientColumns+`
		FROM clients
		WHERE (phone_normalized <> '' AND phone_normalized = $1)
		   OR (email <> '' AND $2 <> '' AND lower(email) = lower($2))
		ORDER BY id ASC
		LIMIT 1
	`, phoneNormalized, email)
	item, err := scanClient(row)
	return item, translateNotFound(err)
}

func (r *ClientRepository) FindByID(ctx context.Context, id int64) (model.Client, error) {
	row := r.db.QueryRow(ctx, `SELECT `+clientColumns+` FROM clients WHERE id = $1`, id)
	item, err := scanClient(row)
	return item, translateNotFound(err)
}

func (r *ClientRepository) Create(ctx context.Context, item model.Client) (model.Client, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO clients (name, phone, phone_normalized, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, item.Name, item.Phone, item.PhoneNormalized, item.Email).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

// List is the source for GET /clients: every client enriched with both tag
// lists and a summary of their most recent activity, for the client-list
// page's cards. Mirrors LeadRepository.ListWithClient's shape (scalar
// subqueries, not a plain join, to avoid fanning out rows).
func (r *ClientRepository) List(ctx context.Context) ([]model.ClientListItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			`+clientColumns+`,
			COALESCE((
				SELECT jsonb_agg(jsonb_build_object('id', pt.id, 'name', pt.name, 'color', pt.color) ORDER BY pt.name)
				FROM client_product_tags cpt JOIN product_tags pt ON pt.id = cpt.product_tag_id
				WHERE cpt.client_id = c.id
			), '[]') AS product_tags,
			COALESCE((
				SELECT jsonb_agg(jsonb_build_object('id', ct.id, 'name', ct.name, 'color', ct.color) ORDER BY ct.name)
				FROM client_client_type_tags cctt JOIN client_type_tags ct ON ct.id = cctt.client_type_tag_id
				WHERE cctt.client_id = c.id
			), '[]') AS client_type_tags,
			(SELECT count(*) FROM leads l WHERE l.client_id = c.id) AS request_count,
			(SELECT l.status FROM leads l WHERE l.client_id = c.id ORDER BY l.created_at DESC LIMIT 1) AS latest_status,
			(SELECT l.created_at FROM leads l WHERE l.client_id = c.id ORDER BY l.created_at DESC LIMIT 1) AS latest_request_at,
			(SELECT cm.created_at FROM client_comments cm WHERE cm.client_id = c.id ORDER BY cm.created_at DESC LIMIT 1) AS latest_comment_at
		FROM clients c
		ORDER BY c.updated_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.ClientListItem{}
	for rows.Next() {
		var item model.ClientListItem
		var productTagsJSON, clientTypeTagsJSON []byte
		var latestStatus *model.LeadStatus
		if err := rows.Scan(
			&item.ID, &item.Name, &item.Phone, &item.PhoneNormalized, &item.Email, &item.CreatedAt, &item.UpdatedAt,
			&productTagsJSON, &clientTypeTagsJSON,
			&item.RequestCount, &latestStatus, &item.LatestRequestAt, &item.LatestCommentAt,
		); err != nil {
			return nil, err
		}
		if latestStatus != nil {
			item.LatestStatus = *latestStatus
		}
		if err := json.Unmarshal(productTagsJSON, &item.ProductTags); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(clientTypeTagsJSON, &item.ClientTypeTags); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// RefreshContactInfo overwrites the client's profile with the latest
// submitted contact details — the newest submission is more likely to be
// current than whatever's on file. The lead's own historical snapshot is
// untouched by this (see model.Lead).
func (r *ClientRepository) RefreshContactInfo(ctx context.Context, id int64, item model.Client) (model.Client, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE clients
		SET name = $1, phone = $2, phone_normalized = $3, email = $4, updated_at = now()
		WHERE id = $5
		RETURNING `+clientColumns+`
	`, item.Name, item.Phone, item.PhoneNormalized, item.Email, id)
	updated, err := scanClient(row)
	return updated, translateNotFound(err)
}
