package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type LessonRepository struct {
	db *pgxpool.Pool
}

func NewLessonRepository(db *pgxpool.Pool) *LessonRepository {
	return &LessonRepository{db: db}
}

const lessonColumns = "id, course_block_id, course_id, name, description, sort_order, created_at, updated_at"

func scanLesson(row pgx.Row) (model.Lesson, error) {
	var item model.Lesson
	err := row.Scan(&item.ID, &item.CourseBlockID, &item.CourseID, &item.Name, &item.Description, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *LessonRepository) ListByCourseBlockID(ctx context.Context, courseBlockID int64) ([]model.Lesson, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+lessonColumns+`
		FROM lessons
		WHERE course_block_id = $1
		ORDER BY sort_order, id
	`, courseBlockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Lesson{}
	for rows.Next() {
		item, err := scanLesson(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// ListByCourseBlockIDs fetches lessons for several blocks in one query —
// used by CourseCatalogService to avoid a per-block round trip when
// assembling a full course page (architecture review finding #2).
func (r *LessonRepository) ListByCourseBlockIDs(ctx context.Context, courseBlockIDs []int64) ([]model.Lesson, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+lessonColumns+`
		FROM lessons
		WHERE course_block_id = ANY($1)
		ORDER BY course_block_id, sort_order, id
	`, courseBlockIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Lesson{}
	for rows.Next() {
		item, err := scanLesson(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// ListByCourseID is the course-without-blocks counterpart of
// ListByCourseBlockID — a course editing its lessons directly, with no
// CourseBlock row involved at all (see model.Lesson's doc comment).
func (r *LessonRepository) ListByCourseID(ctx context.Context, courseID int64) ([]model.Lesson, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+lessonColumns+`
		FROM lessons
		WHERE course_id = $1
		ORDER BY sort_order, id
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Lesson{}
	for rows.Next() {
		item, err := scanLesson(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// Create works for either parent — item.CourseBlockID or item.CourseID is
// set (never both, enforced by the lessons_exactly_one_parent CHECK), the
// other stays nil and goes in as NULL.
func (r *LessonRepository) Create(ctx context.Context, item model.Lesson) (model.Lesson, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO lessons (course_block_id, course_id, name, description, sort_order)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`, item.CourseBlockID, item.CourseID, item.Name, item.Description, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

// UpdateByCourseBlock matches on id AND course_block_id — a URL like
// /course-blocks/999/lessons/5 can't touch (or, previously, silently claim
// in its response to have moved) a lesson that actually belongs to a
// different block (architecture review finding #3).
func (r *LessonRepository) UpdateByCourseBlock(ctx context.Context, item model.Lesson) (model.Lesson, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE lessons
		SET name = $1, description = $2, sort_order = $3, updated_at = now()
		WHERE id = $4 AND course_block_id = $5
		RETURNING updated_at
	`, item.Name, item.Description, item.SortOrder, item.ID, item.CourseBlockID).Scan(&item.UpdatedAt)
	return item, translateNotFound(err)
}

// UpdateByCourse is UpdateByCourseBlock's course-without-blocks counterpart.
func (r *LessonRepository) UpdateByCourse(ctx context.Context, item model.Lesson) (model.Lesson, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE lessons
		SET name = $1, description = $2, sort_order = $3, updated_at = now()
		WHERE id = $4 AND course_id = $5
		RETURNING updated_at
	`, item.Name, item.Description, item.SortOrder, item.ID, item.CourseID).Scan(&item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *LessonRepository) DeleteByCourseBlock(ctx context.Context, courseBlockID, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM lessons WHERE id = $1 AND course_block_id = $2`, id, courseBlockID)
	return checkDeleted(tag, err)
}

func (r *LessonRepository) DeleteByCourse(ctx context.Context, courseID, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM lessons WHERE id = $1 AND course_id = $2`, id, courseID)
	return checkDeleted(tag, err)
}
