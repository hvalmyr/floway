package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type CustomDisplayStyleRepository struct {
	db *pgxpool.Pool
}

func NewCustomDisplayStyleRepository(db *pgxpool.Pool) *CustomDisplayStyleRepository {
	return &CustomDisplayStyleRepository{db: db}
}

func (r *CustomDisplayStyleRepository) List(ctx context.Context) ([]model.CustomDisplayStyle, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, bg_color, text_color, sort_order, created_at, updated_at
		FROM custom_display_styles
		ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.CustomDisplayStyle{}
	for rows.Next() {
		var item model.CustomDisplayStyle
		if err := rows.Scan(&item.ID, &item.Name, &item.BgColor, &item.TextColor, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *CustomDisplayStyleRepository) Create(ctx context.Context, item model.CustomDisplayStyle) (model.CustomDisplayStyle, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO custom_display_styles (name, bg_color, text_color, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, item.Name, item.BgColor, item.TextColor, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *CustomDisplayStyleRepository) Update(ctx context.Context, item model.CustomDisplayStyle) (model.CustomDisplayStyle, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE custom_display_styles
		SET name = $1, bg_color = $2, text_color = $3, sort_order = $4, updated_at = now()
		WHERE id = $5
		RETURNING updated_at
	`, item.Name, item.BgColor, item.TextColor, item.SortOrder, item.ID).Scan(&item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *CustomDisplayStyleRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM custom_display_styles WHERE id = $1`, id)
	return checkDeleted(tag, err)
}
