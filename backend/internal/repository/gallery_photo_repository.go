package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type GalleryPhotoRepository struct {
	db *pgxpool.Pool
}

func NewGalleryPhotoRepository(db *pgxpool.Pool) *GalleryPhotoRepository {
	return &GalleryPhotoRepository{db: db}
}

func (r *GalleryPhotoRepository) List(ctx context.Context) ([]model.GalleryPhoto, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, image, sort_order, created_at, updated_at
		FROM gallery_photos
		ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.GalleryPhoto{}
	for rows.Next() {
		var item model.GalleryPhoto
		if err := rows.Scan(&item.ID, &item.Image, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *GalleryPhotoRepository) Create(ctx context.Context, item model.GalleryPhoto) (model.GalleryPhoto, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO gallery_photos (image, sort_order)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`, item.Image, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *GalleryPhotoRepository) Update(ctx context.Context, item model.GalleryPhoto) (model.GalleryPhoto, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE gallery_photos
		SET image = $1, sort_order = $2, updated_at = now()
		WHERE id = $3
		RETURNING updated_at
	`, item.Image, item.SortOrder, item.ID).Scan(&item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *GalleryPhotoRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM gallery_photos WHERE id = $1`, id)
	return checkDeleted(tag, err)
}
