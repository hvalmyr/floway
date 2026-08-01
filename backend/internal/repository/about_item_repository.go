package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type AboutItemRepository struct {
	db *pgxpool.Pool
}

func NewAboutItemRepository(db *pgxpool.Pool) *AboutItemRepository {
	return &AboutItemRepository{db: db}
}

func (r *AboutItemRepository) List(ctx context.Context) ([]model.AboutItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, badge, description, sort_order, created_at, updated_at
		FROM about_items
		ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.AboutItem{}
	for rows.Next() {
		var item model.AboutItem
		if err := rows.Scan(&item.ID, &item.Badge, &item.Description, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *AboutItemRepository) Create(ctx context.Context, item model.AboutItem) (model.AboutItem, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO about_items (badge, description, sort_order)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`, item.Badge, item.Description, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *AboutItemRepository) Update(ctx context.Context, item model.AboutItem) (model.AboutItem, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE about_items
		SET badge = $1, description = $2, sort_order = $3, updated_at = now()
		WHERE id = $4
		RETURNING updated_at
	`, item.Badge, item.Description, item.SortOrder, item.ID).Scan(&item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *AboutItemRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM about_items WHERE id = $1`, id)
	return checkDeleted(tag, err)
}
