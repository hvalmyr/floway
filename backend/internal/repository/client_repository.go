package repository

import (
	"context"

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
