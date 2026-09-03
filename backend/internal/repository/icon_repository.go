package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type IconRepository struct {
	db *pgxpool.Pool
}

func NewIconRepository(db *pgxpool.Pool) *IconRepository {
	return &IconRepository{db: db}
}

func (r *IconRepository) List(ctx context.Context) ([]model.Icon, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, svg, created_at
		FROM icons
		ORDER BY name, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Icon{}
	for rows.Next() {
		var item model.Icon
		if err := rows.Scan(&item.ID, &item.Name, &item.SVG, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *IconRepository) Create(ctx context.Context, item model.Icon) (model.Icon, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO icons (name, svg)
		VALUES ($1, $2)
		RETURNING id, created_at
	`, item.Name, item.SVG).Scan(&item.ID, &item.CreatedAt)
	return item, err
}

func (r *IconRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM icons WHERE id = $1`, id)
	return checkDeleted(tag, err)
}
