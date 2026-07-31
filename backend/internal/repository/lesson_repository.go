package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type LessonRepository struct {
	db *pgxpool.Pool
}

func NewLessonRepository(db *pgxpool.Pool) *LessonRepository {
	return &LessonRepository{db: db}
}

func (r *LessonRepository) ListByCourseBlockID(ctx context.Context, courseBlockID int64) ([]model.Lesson, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, course_block_id, number, title, topics, outcomes, duration_hours, created_at, updated_at
		FROM lessons
		WHERE course_block_id = $1
		ORDER BY number, id
	`, courseBlockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Lesson{}
	for rows.Next() {
		var item model.Lesson
		if err := rows.Scan(&item.ID, &item.CourseBlockID, &item.Number, &item.Title, &item.Topics, &item.Outcomes, &item.DurationHours, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *LessonRepository) Create(ctx context.Context, item model.Lesson) (model.Lesson, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO lessons (course_block_id, number, title, topics, outcomes, duration_hours)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`, item.CourseBlockID, item.Number, item.Title, item.Topics, item.Outcomes, item.DurationHours).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *LessonRepository) Update(ctx context.Context, item model.Lesson) (model.Lesson, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE lessons
		SET number = $1, title = $2, topics = $3, outcomes = $4, duration_hours = $5, updated_at = now()
		WHERE id = $6
		RETURNING updated_at
	`, item.Number, item.Title, item.Topics, item.Outcomes, item.DurationHours, item.ID).Scan(&item.UpdatedAt)
	return item, err
}

func (r *LessonRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM lessons WHERE id = $1`, id)
	return err
}
