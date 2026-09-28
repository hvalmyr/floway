package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type SiteButtonRepository struct {
	db *pgxpool.Pool
}

func NewSiteButtonRepository(db *pgxpool.Pool) *SiteButtonRepository {
	return &SiteButtonRepository{db: db}
}

func (r *SiteButtonRepository) List(ctx context.Context) ([]model.SiteButton, error) {
	rows, err := r.db.Query(ctx, `
		SELECT key, label, text, variant, url, updated_at
		FROM site_buttons
		ORDER BY key
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.SiteButton{}
	for rows.Next() {
		var item model.SiteButton
		if err := rows.Scan(&item.Key, &item.Label, &item.Text, &item.Variant, &item.URL, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// Update sets text/variant/url for an existing key. Keys are seeded by
// migration 00064, not created through the API — updating an unknown key
// returns apperr.ErrNotFound.
func (r *SiteButtonRepository) Update(ctx context.Context, key, text, variant, url string) (model.SiteButton, error) {
	item := model.SiteButton{Key: key, Text: text, Variant: variant, URL: url}
	err := r.db.QueryRow(ctx, `
		UPDATE site_buttons
		SET text = $1, variant = $2, url = $3, updated_at = now()
		WHERE key = $4
		RETURNING label, updated_at
	`, text, variant, url, key).Scan(&item.Label, &item.UpdatedAt)
	return item, translateNotFound(err)
}
