package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type CourseBlockRepository struct {
	db *pgxpool.Pool
}

func NewCourseBlockRepository(db *pgxpool.Pool) *CourseBlockRepository {
	return &CourseBlockRepository{db: db}
}

func (r *CourseBlockRepository) ListByCourseID(ctx context.Context, courseID int64) ([]model.CourseBlock, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, course_id, title, lessons_count, hours, price, old_price, sort_order, created_at, updated_at
		FROM course_blocks
		WHERE course_id = $1
		ORDER BY sort_order, id
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.CourseBlock{}
	for rows.Next() {
		var item model.CourseBlock
		if err := rows.Scan(&item.ID, &item.CourseID, &item.Title, &item.LessonsCount, &item.Hours, &item.Price, &item.OldPrice, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *CourseBlockRepository) Create(ctx context.Context, item model.CourseBlock) (model.CourseBlock, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO course_blocks (course_id, title, lessons_count, hours, price, old_price, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`, item.CourseID, item.Title, item.LessonsCount, item.Hours, item.Price, item.OldPrice, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

// Update matches on id AND course_id — a URL like /courses/7/blocks/42
// can't touch a block that actually belongs to a different course
// (architecture review finding #3: the nesting used to be decorative,
// this makes it real).
func (r *CourseBlockRepository) Update(ctx context.Context, item model.CourseBlock) (model.CourseBlock, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE course_blocks
		SET title = $1, lessons_count = $2, hours = $3, price = $4, old_price = $5, sort_order = $6, updated_at = now()
		WHERE id = $7 AND course_id = $8
		RETURNING course_id, updated_at
	`, item.Title, item.LessonsCount, item.Hours, item.Price, item.OldPrice, item.SortOrder, item.ID, item.CourseID).Scan(&item.CourseID, &item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *CourseBlockRepository) Delete(ctx context.Context, courseID, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM course_blocks WHERE id = $1 AND course_id = $2`, id, courseID)
	return checkDeleted(tag, err)
}
