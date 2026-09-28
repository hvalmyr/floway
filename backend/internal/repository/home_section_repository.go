package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type HomeSectionRepository struct {
	db *pgxpool.Pool
}

func NewHomeSectionRepository(db *pgxpool.Pool) *HomeSectionRepository {
	return &HomeSectionRepository{db: db}
}

func (r *HomeSectionRepository) List(ctx context.Context) ([]model.HomeSection, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, key, visible, sort_order, updated_at
		FROM home_sections
		ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.HomeSection{}
	for rows.Next() {
		var item model.HomeSection
		if err := rows.Scan(&item.ID, &item.Key, &item.Visible, &item.SortOrder, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// Update sets visible/sort_order for an existing row. Key is seeded by
// migration 00063, not created through the API — updating an unknown id
// returns apperr.ErrNotFound.
func (r *HomeSectionRepository) Update(ctx context.Context, item model.HomeSection) (model.HomeSection, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE home_sections
		SET visible = $1, sort_order = $2, updated_at = now()
		WHERE id = $3
		RETURNING key, updated_at
	`, item.Visible, item.SortOrder, item.ID).Scan(&item.Key, &item.UpdatedAt)
	return item, translateNotFound(err)
}
