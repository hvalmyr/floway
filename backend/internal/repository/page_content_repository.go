package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type PageContentRepository struct {
	db *pgxpool.Pool
}

func NewPageContentRepository(db *pgxpool.Pool) *PageContentRepository {
	return &PageContentRepository{db: db}
}

func (r *PageContentRepository) List(ctx context.Context) ([]model.PageContent, error) {
	rows, err := r.db.Query(ctx, `
		SELECT key, label, value, updated_at
		FROM page_content
		ORDER BY key
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.PageContent{}
	for rows.Next() {
		var item model.PageContent
		if err := rows.Scan(&item.Key, &item.Label, &item.Value, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// Update sets the value for an existing key. Keys are seeded by migrations,
// not created through the API — updating an unknown key returns pgx.ErrNoRows.
func (r *PageContentRepository) Update(ctx context.Context, key, value string) (model.PageContent, error) {
	var item model.PageContent
	item.Key = key
	item.Value = value
	err := r.db.QueryRow(ctx, `
		UPDATE page_content
		SET value = $1, updated_at = now()
		WHERE key = $2
		RETURNING label, updated_at
	`, value, key).Scan(&item.Label, &item.UpdatedAt)
	return item, err
}
