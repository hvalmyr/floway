package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type FeatureRepository struct {
	db *pgxpool.Pool
}

func NewFeatureRepository(db *pgxpool.Pool) *FeatureRepository {
	return &FeatureRepository{db: db}
}

// featureColumns/scanFeature — see blog_post_repository.go's identical
// comment (architecture review finding #12).
const featureColumns = "id, page, icon, title, description, sort_order, created_at, updated_at"

func scanFeature(row pgx.Row) (model.Feature, error) {
	var item model.Feature
	err := row.Scan(&item.ID, &item.Page, &item.Icon, &item.Title, &item.Description, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *FeatureRepository) List(ctx context.Context) ([]model.Feature, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+featureColumns+`
		FROM features
		ORDER BY page, sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Feature{}
	for rows.Next() {
		item, err := scanFeature(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *FeatureRepository) ListByPage(ctx context.Context, page string) ([]model.Feature, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+featureColumns+`
		FROM features
		WHERE page = $1
		ORDER BY sort_order, id
	`, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Feature{}
	for rows.Next() {
		item, err := scanFeature(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *FeatureRepository) Create(ctx context.Context, item model.Feature) (model.Feature, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO features (page, icon, title, description, sort_order)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`, item.Page, item.Icon, item.Title, item.Description, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *FeatureRepository) Update(ctx context.Context, item model.Feature) (model.Feature, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE features
		SET page = $1, icon = $2, title = $3, description = $4, sort_order = $5, updated_at = now()
		WHERE id = $6
		RETURNING updated_at
	`, item.Page, item.Icon, item.Title, item.Description, item.SortOrder, item.ID).Scan(&item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *FeatureRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM features WHERE id = $1`, id)
	return checkDeleted(tag, err)
}
