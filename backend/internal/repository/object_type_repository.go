package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type ObjectTypeRepository struct {
	db *pgxpool.Pool
}

func NewObjectTypeRepository(db *pgxpool.Pool) *ObjectTypeRepository {
	return &ObjectTypeRepository{db: db}
}

const objectTypeColumns = "id, slug, name, visible, sort_order, created_at, updated_at"

// objectTypeColumnsQualified is the same column list, table-qualified for
// use in queries that join object_types against other tables (see
// LandingPageRepository.ListVisibleWithObjectType) — bare column names would
// otherwise collide with landing_pages' own id/slug/visible/sort_order.
const objectTypeColumnsQualified = "ot.id, ot.slug, ot.name, ot.visible, ot.sort_order, ot.created_at, ot.updated_at"

func scanObjectType(row pgx.Row) (model.ObjectType, error) {
	var item model.ObjectType
	err := row.Scan(&item.ID, &item.Slug, &item.Name, &item.Visible, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

// List returns every object type, hidden ones included — the admin screen
// needs to see and re-show them; the public frontend filters by Visible
// itself (same convention as FeatureRepository.List).
func (r *ObjectTypeRepository) List(ctx context.Context) ([]model.ObjectType, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+objectTypeColumns+`
		FROM object_types
		ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.ObjectType{}
	for rows.Next() {
		item, err := scanObjectType(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *ObjectTypeRepository) Create(ctx context.Context, item model.ObjectType) (model.ObjectType, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO object_types (slug, name, visible, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, item.Slug, item.Name, item.Visible, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *ObjectTypeRepository) Update(ctx context.Context, item model.ObjectType) (model.ObjectType, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE object_types
		SET slug = $1, name = $2, visible = $3, sort_order = $4, updated_at = now()
		WHERE id = $5
		RETURNING updated_at
	`, item.Slug, item.Name, item.Visible, item.SortOrder, item.ID).Scan(&item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *ObjectTypeRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM object_types WHERE id = $1`, id)
	return checkDeleted(tag, err)
}
