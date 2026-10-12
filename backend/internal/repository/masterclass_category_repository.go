package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type MasterclassCategoryRepository struct {
	db *pgxpool.Pool
}

func NewMasterclassCategoryRepository(db *pgxpool.Pool) *MasterclassCategoryRepository {
	return &MasterclassCategoryRepository{db: db}
}

const masterclassCategoryColumns = "id, slug, name, sort_order, created_at, updated_at"

func scanMasterclassCategory(row pgx.Row) (model.MasterclassCategory, error) {
	var item model.MasterclassCategory
	err := row.Scan(&item.ID, &item.Slug, &item.Name, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *MasterclassCategoryRepository) List(ctx context.Context) ([]model.MasterclassCategory, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+masterclassCategoryColumns+`
		FROM masterclass_categories
		ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.MasterclassCategory{}
	for rows.Next() {
		item, err := scanMasterclassCategory(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *MasterclassCategoryRepository) Create(ctx context.Context, item model.MasterclassCategory) (model.MasterclassCategory, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO masterclass_categories (slug, name, sort_order)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`, item.Slug, item.Name, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *MasterclassCategoryRepository) Update(ctx context.Context, item model.MasterclassCategory) (model.MasterclassCategory, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE masterclass_categories
		SET slug = $1, name = $2, sort_order = $3, updated_at = now()
		WHERE id = $4
		RETURNING updated_at
	`, item.Slug, item.Name, item.SortOrder, item.ID).Scan(&item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *MasterclassCategoryRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM masterclass_categories WHERE id = $1`, id)
	return checkDeleted(tag, err)
}
