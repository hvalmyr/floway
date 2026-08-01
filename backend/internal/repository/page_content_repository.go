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
		SELECT key, label, value, type, updated_at
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
		if err := rows.Scan(&item.Key, &item.Label, &item.Value, &item.Type, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// Update sets the value for an existing key. Keys are seeded by migrations,
// not created through the API — updating an unknown key returns apperr.ErrNotFound.
//
// It also returns the value the key held before this update (via a CTE, so
// the read and the write share one snapshot instead of racing) — the caller
// uses it to clean up an image that's about to become orphaned in Garage.
func (r *PageContentRepository) Update(ctx context.Context, key, value string) (item model.PageContent, previousValue string, err error) {
	item.Key = key
	item.Value = value
	err = r.db.QueryRow(ctx, `
		WITH previous AS (
			SELECT value FROM page_content WHERE key = $2
		)
		UPDATE page_content
		SET value = $1, updated_at = now()
		WHERE key = $2
		RETURNING label, type, updated_at, (SELECT value FROM previous)
	`, value, key).Scan(&item.Label, &item.Type, &item.UpdatedAt, &previousValue)
	return item, previousValue, translateNotFound(err)
}
