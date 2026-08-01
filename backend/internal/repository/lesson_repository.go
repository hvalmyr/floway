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

// ListByCourseBlockIDs fetches lessons for several blocks in one query —
// used by CourseDetailService to avoid a per-block round trip when
// assembling a full course page (architecture review finding #2).
func (r *LessonRepository) ListByCourseBlockIDs(ctx context.Context, courseBlockIDs []int64) ([]model.Lesson, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, course_block_id, number, title, topics, outcomes, duration_hours, created_at, updated_at
		FROM lessons
		WHERE course_block_id = ANY($1)
		ORDER BY course_block_id, number, id
	`, courseBlockIDs)
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

// Update matches on id AND course_block_id — a URL like
// /course-blocks/999/lessons/5 can't touch (or, previously, silently claim
// in its response to have moved) a lesson that actually belongs to a
// different block (architecture review finding #3).
func (r *LessonRepository) Update(ctx context.Context, item model.Lesson) (model.Lesson, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE lessons
		SET number = $1, title = $2, topics = $3, outcomes = $4, duration_hours = $5, updated_at = now()
		WHERE id = $6 AND course_block_id = $7
		RETURNING updated_at
	`, item.Number, item.Title, item.Topics, item.Outcomes, item.DurationHours, item.ID, item.CourseBlockID).Scan(&item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *LessonRepository) Delete(ctx context.Context, courseBlockID, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM lessons WHERE id = $1 AND course_block_id = $2`, id, courseBlockID)
	return checkDeleted(tag, err)
}
