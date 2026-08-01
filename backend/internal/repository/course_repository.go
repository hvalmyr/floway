package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type CourseRepository struct {
	db *pgxpool.Pool
}

func NewCourseRepository(db *pgxpool.Pool) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) List(ctx context.Context) ([]model.Course, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, slug, title, short_description, full_description, status, cover_image, gallery, sort_order, created_at, updated_at
		FROM courses
		ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Course{}
	for rows.Next() {
		var item model.Course
		if err := rows.Scan(
			&item.ID,
			&item.Slug,
			&item.Title,
			&item.ShortDesc,
			&item.FullDesc,
			&item.Status,
			&item.CoverImage,
			&item.Gallery,
			&item.SortOrder,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *CourseRepository) FindBySlug(ctx context.Context, slug string) (model.Course, error) {
	var item model.Course
	err := r.db.QueryRow(ctx, `
		SELECT id, slug, title, short_description, full_description, status, cover_image, gallery, sort_order, created_at, updated_at
		FROM courses
		WHERE slug = $1
	`, slug).Scan(
		&item.ID,
		&item.Slug,
		&item.Title,
		&item.ShortDesc,
		&item.FullDesc,
		&item.Status,
		&item.CoverImage,
		&item.Gallery,
		&item.SortOrder,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	return item, err
}

func (r *CourseRepository) Create(ctx context.Context, item model.Course) (model.Course, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO courses (slug, title, short_description, full_description, status, cover_image, gallery, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`,
		item.Slug,
		item.Title,
		item.ShortDesc,
		item.FullDesc,
		item.Status,
		item.CoverImage,
		item.Gallery,
		item.SortOrder,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *CourseRepository) Update(ctx context.Context, item model.Course) (model.Course, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE courses
		SET slug = $1, title = $2, short_description = $3, full_description = $4, status = $5,
			cover_image = $6, gallery = $7, sort_order = $8, updated_at = now()
		WHERE id = $9
		RETURNING updated_at
	`,
		item.Slug,
		item.Title,
		item.ShortDesc,
		item.FullDesc,
		item.Status,
		item.CoverImage,
		item.Gallery,
		item.SortOrder,
		item.ID,
	).Scan(&item.UpdatedAt)
	return item, err
}

func (r *CourseRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM courses WHERE id = $1`, id)
	return err
}
