package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type CourseSectionRepository struct {
	db *pgxpool.Pool
}

func NewCourseSectionRepository(db *pgxpool.Pool) *CourseSectionRepository {
	return &CourseSectionRepository{db: db}
}

func (r *CourseSectionRepository) List(ctx context.Context) ([]model.CourseSection, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, heading, description, visible, sort_order, created_at, updated_at
		FROM course_sections
		ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.CourseSection{}
	for rows.Next() {
		var item model.CourseSection
		if err := rows.Scan(&item.ID, &item.Heading, &item.Description, &item.Visible, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *CourseSectionRepository) Create(ctx context.Context, item model.CourseSection) (model.CourseSection, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO course_sections (heading, description, visible, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, item.Heading, item.Description, item.Visible, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *CourseSectionRepository) Update(ctx context.Context, item model.CourseSection) (model.CourseSection, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE course_sections
		SET heading = $1, description = $2, visible = $3, sort_order = $4, updated_at = now()
		WHERE id = $5
		RETURNING updated_at
	`, item.Heading, item.Description, item.Visible, item.SortOrder, item.ID).Scan(&item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *CourseSectionRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM course_sections WHERE id = $1`, id)
	return checkDeleted(tag, err)
}
