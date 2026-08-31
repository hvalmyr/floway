package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type MasterclassRepository struct {
	db *pgxpool.Pool
}

func NewMasterclassRepository(db *pgxpool.Pool) *MasterclassRepository {
	return &MasterclassRepository{db: db}
}

// masterclassColumns/scanMasterclass — see blog_post_repository.go's
// identical comment (architecture review finding #12).
const masterclassColumns = "id, slug, title, description, description2, ending_text, duration, price, cover_image, status, created_at, updated_at"

func scanMasterclass(row pgx.Row) (model.Masterclass, error) {
	var item model.Masterclass
	err := row.Scan(&item.ID, &item.Slug, &item.Title, &item.Description, &item.Description2, &item.EndingText, &item.Duration, &item.Price, &item.CoverImage, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *MasterclassRepository) List(ctx context.Context) ([]model.Masterclass, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+masterclassColumns+`
		FROM masterclasses
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Masterclass{}
	for rows.Next() {
		item, err := scanMasterclass(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *MasterclassRepository) FindBySlug(ctx context.Context, slug string) (model.Masterclass, error) {
	row := r.db.QueryRow(ctx, `
		SELECT `+masterclassColumns+`
		FROM masterclasses
		WHERE slug = $1
	`, slug)
	item, err := scanMasterclass(row)
	return item, translateNotFound(err)
}

func (r *MasterclassRepository) Create(ctx context.Context, item model.Masterclass) (model.Masterclass, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO masterclasses (
			slug, title, description, description2, ending_text, duration, price, cover_image, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`,
		item.Slug, item.Title, item.Description, item.Description2, item.EndingText, item.Duration,
		item.Price, item.CoverImage, item.Status,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *MasterclassRepository) Update(ctx context.Context, item model.Masterclass) (model.Masterclass, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE masterclasses
		SET slug = $1, title = $2, description = $3, description2 = $4, ending_text = $5, duration = $6,
		    price = $7, cover_image = $8, status = $9,
		    updated_at = now()
		WHERE id = $10
		RETURNING updated_at
	`,
		item.Slug, item.Title, item.Description, item.Description2, item.EndingText, item.Duration,
		item.Price, item.CoverImage, item.Status, item.ID,
	).Scan(&item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *MasterclassRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM masterclasses WHERE id = $1`, id)
	return checkDeleted(tag, err)
}
