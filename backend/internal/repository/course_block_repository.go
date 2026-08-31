package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type CourseBlockRepository struct {
	db *pgxpool.Pool
}

func NewCourseBlockRepository(db *pgxpool.Pool) *CourseBlockRepository {
	return &CourseBlockRepository{db: db}
}

const courseBlockColumns = "id, course_id, block_name, description, block_cover, lesson_count, time_length, price, display_style, visible, sort_order, created_at, updated_at"

func scanCourseBlock(row pgx.Row) (model.CourseBlock, error) {
	var item model.CourseBlock
	err := row.Scan(&item.ID, &item.CourseID, &item.BlockName, &item.Description, &item.BlockCover, &item.LessonCount, &item.TimeLength, &item.Price, &item.DisplayStyle, &item.Visible, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *CourseBlockRepository) ListByCourseID(ctx context.Context, courseID int64) ([]model.CourseBlock, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+courseBlockColumns+`
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
		item, err := scanCourseBlock(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// ListByCourseIDs batches blocks for several courses in one query — used by
// CourseCatalogService to build homepage course summaries and the course
// detail page without a per-course round trip.
func (r *CourseBlockRepository) ListByCourseIDs(ctx context.Context, courseIDs []int64) ([]model.CourseBlock, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+courseBlockColumns+`
		FROM course_blocks
		WHERE course_id = ANY($1)
		ORDER BY course_id, sort_order, id
	`, courseIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.CourseBlock{}
	for rows.Next() {
		item, err := scanCourseBlock(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *CourseBlockRepository) Create(ctx context.Context, item model.CourseBlock) (model.CourseBlock, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO course_blocks (course_id, block_name, description, block_cover, lesson_count, time_length, price, display_style, visible, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`, item.CourseID, item.BlockName, item.Description, item.BlockCover, item.LessonCount, item.TimeLength, item.Price, item.DisplayStyle, item.Visible, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

// Update matches on id AND course_id — a URL like /courses/7/blocks/42
// can't touch a block that actually belongs to a different course
// (architecture review finding #3).
func (r *CourseBlockRepository) Update(ctx context.Context, item model.CourseBlock) (model.CourseBlock, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE course_blocks
		SET block_name = $1, description = $2, block_cover = $3, lesson_count = $4, time_length = $5, price = $6, display_style = $7, visible = $8, sort_order = $9, updated_at = now()
		WHERE id = $10 AND course_id = $11
		RETURNING course_id, updated_at
	`, item.BlockName, item.Description, item.BlockCover, item.LessonCount, item.TimeLength, item.Price, item.DisplayStyle, item.Visible, item.SortOrder, item.ID, item.CourseID).Scan(&item.CourseID, &item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *CourseBlockRepository) Delete(ctx context.Context, courseID, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM course_blocks WHERE id = $1 AND course_id = $2`, id, courseID)
	return checkDeleted(tag, err)
}
