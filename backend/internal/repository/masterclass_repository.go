package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type MasterclassRepository struct {
	db *pgxpool.Pool
}

func NewMasterclassRepository(db *pgxpool.Pool) *MasterclassRepository {
	return &MasterclassRepository{db: db}
}

func (r *MasterclassRepository) List(ctx context.Context) ([]model.Masterclass, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, slug, title, short_description, full_description, ending_text, duration,
		       price_group, price_individual, price_description, cover_image, status, created_at, updated_at
		FROM masterclasses
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Masterclass{}
	for rows.Next() {
		var item model.Masterclass
		if err := rows.Scan(
			&item.ID, &item.Slug, &item.Title, &item.ShortDesc, &item.FullDesc, &item.EndingText, &item.Duration,
			&item.PriceGroup, &item.PriceIndividual, &item.PriceDescription, &item.CoverImage, &item.Status,
			&item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *MasterclassRepository) FindBySlug(ctx context.Context, slug string) (model.Masterclass, error) {
	var item model.Masterclass
	err := r.db.QueryRow(ctx, `
		SELECT id, slug, title, short_description, full_description, ending_text, duration,
		       price_group, price_individual, price_description, cover_image, status, created_at, updated_at
		FROM masterclasses
		WHERE slug = $1
	`, slug).Scan(
		&item.ID, &item.Slug, &item.Title, &item.ShortDesc, &item.FullDesc, &item.EndingText, &item.Duration,
		&item.PriceGroup, &item.PriceIndividual, &item.PriceDescription, &item.CoverImage, &item.Status,
		&item.CreatedAt, &item.UpdatedAt,
	)
	return item, translateNotFound(err)
}

func (r *MasterclassRepository) Create(ctx context.Context, item model.Masterclass) (model.Masterclass, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO masterclasses (
			slug, title, short_description, full_description, ending_text, duration,
			price_group, price_individual, price_description, cover_image, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`,
		item.Slug, item.Title, item.ShortDesc, item.FullDesc, item.EndingText, item.Duration,
		item.PriceGroup, item.PriceIndividual, item.PriceDescription, item.CoverImage, item.Status,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *MasterclassRepository) Update(ctx context.Context, item model.Masterclass) (model.Masterclass, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE masterclasses
		SET slug = $1, title = $2, short_description = $3, full_description = $4, ending_text = $5, duration = $6,
		    price_group = $7, price_individual = $8, price_description = $9, cover_image = $10, status = $11,
		    updated_at = now()
		WHERE id = $12
		RETURNING updated_at
	`,
		item.Slug, item.Title, item.ShortDesc, item.FullDesc, item.EndingText, item.Duration,
		item.PriceGroup, item.PriceIndividual, item.PriceDescription, item.CoverImage, item.Status, item.ID,
	).Scan(&item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *MasterclassRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM masterclasses WHERE id = $1`, id)
	return checkDeleted(tag, err)
}
