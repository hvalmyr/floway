package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type LeadRepository struct {
	db *pgxpool.Pool
}

func NewLeadRepository(db *pgxpool.Pool) *LeadRepository {
	return &LeadRepository{db: db}
}

const leadColumns = "id, name, phone, email, contact_method, source, request_type, related_id, related_slug, status, created_at, client_id, needs_status_review"

// leadColumnsQualified is the same column list, table-qualified for use in
// queries that join leads against other tables (course/masterclass name
// resolution below) — bare column names would otherwise collide with
// identically-named columns on courses/masterclasses (id, name).
const leadColumnsQualified = "l.id, l.name, l.phone, l.email, l.contact_method, l.source, l.request_type, l.related_id, l.related_slug, l.status, l.created_at, l.client_id, l.needs_status_review"

// relatedNameJoin resolves a lead's RelatedSlug to the human-readable
// course/masterclass title it points at, gated by RequestType since the
// same slug column is shared between the two entity types. LEFT JOINs (not
// inner) so a lead with a blank slug, a trial-lesson lead, or a slug whose
// course/masterclass was since deleted/renamed still returns a row — just
// with related_name empty, and the frontend falls back to the raw slug.
const relatedNameJoin = `
	LEFT JOIN courses crs ON l.request_type = 'course' AND l.related_slug <> '' AND crs.slug = l.related_slug
	LEFT JOIN masterclasses mc ON l.request_type = 'masterclass' AND l.related_slug <> '' AND mc.slug = l.related_slug
`

func scanLead(row pgx.Row) (model.Lead, error) {
	var item model.Lead
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Email,
		&item.ContactMethod,
		&item.Source,
		&item.RequestType,
		&item.RelatedID,
		&item.RelatedSlug,
		&item.Status,
		&item.CreatedAt,
		&item.ClientID,
		&item.NeedsStatusReview,
	)
	return item, err
}

func scanLeadWithRelatedName(row pgx.Row) (model.Lead, error) {
	var item model.Lead
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Email,
		&item.ContactMethod,
		&item.Source,
		&item.RequestType,
		&item.RelatedID,
		&item.RelatedSlug,
		&item.Status,
		&item.CreatedAt,
		&item.ClientID,
		&item.NeedsStatusReview,
		&item.RelatedName,
	)
	return item, err
}

// ListByClientID backs the client detail page's request history.
func (r *LeadRepository) ListByClientID(ctx context.Context, clientID int64) ([]model.Lead, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+leadColumnsQualified+`, COALESCE(crs.name, mc.title, '') AS related_name
		FROM leads l
		`+relatedNameJoin+`
		WHERE l.client_id = $1
		ORDER BY l.created_at DESC
	`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Lead{}
	for rows.Next() {
		item, err := scanLeadWithRelatedName(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// ListWithClient is the source for GET /leads: every lead enriched with its
// client profile, both tag lists, the latest comment excerpt, and the next
// open reminder date — everything the request-list card needs in one round
// trip. Tags/latest-comment/next-reminder are pulled via scalar subqueries
// rather than a plain join, since a join against three one-to-many tables
// would fan out each lead row once per tag/comment/reminder.
func (r *LeadRepository) ListWithClient(ctx context.Context) ([]model.LeadListItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			l.id, l.name, l.phone, l.email, l.contact_method, l.source, l.request_type,
			l.related_id, l.related_slug, l.status, l.created_at, l.client_id, l.needs_status_review,
			COALESCE(crs.name, mc.title, '') AS related_name,
			c.id, c.name, c.phone, c.phone_normalized, c.email, c.created_at, c.updated_at,
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
			(SELECT cm.text FROM client_comments cm WHERE cm.client_id = c.id ORDER BY cm.created_at DESC LIMIT 1) AS latest_comment_text,
			(SELECT cm.created_at FROM client_comments cm WHERE cm.client_id = c.id ORDER BY cm.created_at DESC LIMIT 1) AS latest_comment_at,
			(SELECT MIN(rm.remind_at) FROM reminders rm WHERE rm.client_id = c.id AND rm.completed_at IS NULL) AS next_reminder_at
		FROM leads l
		JOIN clients c ON c.id = l.client_id
		`+relatedNameJoin+`
		ORDER BY l.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.LeadListItem{}
	for rows.Next() {
		var item model.LeadListItem
		var productTagsJSON, clientTypeTagsJSON []byte
		// latestCommentText is nullable — a client with no comments yet has
		// no "latest" one — so it scans into a *string, not model.Lead-
		// ListItem's plain string field directly.
		var latestCommentText *string
		if err := rows.Scan(
			&item.ID, &item.Name, &item.Phone, &item.Email, &item.ContactMethod, &item.Source, &item.RequestType,
			&item.RelatedID, &item.RelatedSlug, &item.Status, &item.CreatedAt, &item.ClientID, &item.NeedsStatusReview,
			&item.RelatedName,
			&item.Client.ID, &item.Client.Name, &item.Client.Phone, &item.Client.PhoneNormalized, &item.Client.Email, &item.Client.CreatedAt, &item.Client.UpdatedAt,
			&productTagsJSON, &clientTypeTagsJSON,
			&latestCommentText, &item.LatestCommentAt,
			&item.NextReminderAt,
		); err != nil {
			return nil, err
		}
		if latestCommentText != nil {
			item.LatestCommentText = *latestCommentText
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

func (r *LeadRepository) Create(ctx context.Context, item model.Lead) (model.Lead, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO leads (name, phone, email, contact_method, source, request_type, related_id, related_slug, status, client_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at
	`,
		item.Name,
		item.Phone,
		item.Email,
		item.ContactMethod,
		item.Source,
		item.RequestType,
		item.RelatedID,
		item.RelatedSlug,
		item.Status,
		item.ClientID,
	).Scan(&item.ID, &item.CreatedAt)
	return item, err
}

// UpdateStatus also clears NeedsStatusReview — any explicit status pick,
// whether reclassifying a migrated lead to closed_lost or reaffirming
// closed_won, counts as "reviewed" (see migration 00031).
func (r *LeadRepository) UpdateStatus(ctx context.Context, id int64, status model.LeadStatus) (model.Lead, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE leads
		SET status = $1, needs_status_review = false
		WHERE id = $2
		RETURNING `+leadColumns+`
	`, status, id)
	item, err := scanLead(row)
	return item, translateNotFound(err)
}

// DismissReview clears NeedsStatusReview without changing status — for the
// case where the migrated closed_won default was already correct.
func (r *LeadRepository) DismissReview(ctx context.Context, id int64) (model.Lead, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE leads
		SET needs_status_review = false
		WHERE id = $1
		RETURNING `+leadColumns+`
	`, id)
	item, err := scanLead(row)
	return item, translateNotFound(err)
}

// CountByStatus powers the conversion-rate stat — how many leads currently
// sit in each of the given statuses.
func (r *LeadRepository) CountByStatus(ctx context.Context, statuses ...model.LeadStatus) (map[model.LeadStatus]int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT status, count(*)
		FROM leads
		WHERE status = ANY($1)
		GROUP BY status
	`, statuses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[model.LeadStatus]int, len(statuses))
	for rows.Next() {
		var status model.LeadStatus
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		counts[status] = count
	}
	return counts, rows.Err()
}

func (r *LeadRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM leads WHERE id = $1`, id)
	return checkDeleted(tag, err)
}
