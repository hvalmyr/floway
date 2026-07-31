package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type LeadRepository struct {
	db *pgxpool.Pool
}

func NewLeadRepository(db *pgxpool.Pool) *LeadRepository {
	return &LeadRepository{db: db}
}

func (r *LeadRepository) List(ctx context.Context) ([]model.Lead, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, phone, email, contact_method, source, request_type, related_id, status, created_at
		FROM leads
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Lead
	for rows.Next() {
		var item model.Lead
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Email,
			&item.ContactMethod,
			&item.Source,
			&item.RequestType,
			&item.RelatedID,
			&item.Status,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *LeadRepository) Create(ctx context.Context, item model.Lead) (model.Lead, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO leads (name, phone, email, contact_method, source, request_type, related_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`,
		item.Name,
		item.Phone,
		item.Email,
		item.ContactMethod,
		item.Source,
		item.RequestType,
		item.RelatedID,
		item.Status,
	).Scan(&item.ID, &item.CreatedAt)
	return item, err
}

func (r *LeadRepository) UpdateStatus(ctx context.Context, id int64, status model.LeadStatus) (model.Lead, error) {
	var item model.Lead
	err := r.db.QueryRow(ctx, `
		UPDATE leads
		SET status = $1
		WHERE id = $2
		RETURNING id, name, phone, email, contact_method, source, request_type, related_id, status, created_at
	`, status, id).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Email,
		&item.ContactMethod,
		&item.Source,
		&item.RequestType,
		&item.RelatedID,
		&item.Status,
		&item.CreatedAt,
	)
	return item, err
}

func (r *LeadRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM leads WHERE id = $1`, id)
	return err
}
