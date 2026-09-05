package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type CourseFAQRepository struct {
	db *pgxpool.Pool
}

func NewCourseFAQRepository(db *pgxpool.Pool) *CourseFAQRepository {
	return &CourseFAQRepository{db: db}
}

const courseFAQColumns = "id, course_id, question, answer, sort_order, created_at, updated_at"

func scanCourseFAQItem(row pgx.Row) (model.CourseFAQItem, error) {
	var item model.CourseFAQItem
	err := row.Scan(&item.ID, &item.CourseID, &item.Question, &item.Answer, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *CourseFAQRepository) ListByCourseID(ctx context.Context, courseID int64) ([]model.CourseFAQItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+courseFAQColumns+`
		FROM course_faq_items
		WHERE course_id = $1
		ORDER BY sort_order, id
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.CourseFAQItem{}
	for rows.Next() {
		item, err := scanCourseFAQItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// ListByCourseIDs batches FAQ items for several courses in one query — same
// reasoning as CourseBlockRepository.ListByCourseIDs, kept ready for any
// future aggregate (e.g. a homepage listing) that needs it; GetFullBySlug
// itself only ever asks for one course's items via ListByCourseID.
func (r *CourseFAQRepository) ListByCourseIDs(ctx context.Context, courseIDs []int64) ([]model.CourseFAQItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+courseFAQColumns+`
		FROM course_faq_items
		WHERE course_id = ANY($1)
		ORDER BY course_id, sort_order, id
	`, courseIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.CourseFAQItem{}
	for rows.Next() {
		item, err := scanCourseFAQItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *CourseFAQRepository) Create(ctx context.Context, item model.CourseFAQItem) (model.CourseFAQItem, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO course_faq_items (course_id, question, answer, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, item.CourseID, item.Question, item.Answer, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

// Update matches on id AND course_id — a URL like /courses/7/faq-items/42
// can't touch an item that actually belongs to a different course (mirrors
// CourseBlockRepository.Update).
func (r *CourseFAQRepository) Update(ctx context.Context, item model.CourseFAQItem) (model.CourseFAQItem, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE course_faq_items
		SET question = $1, answer = $2, sort_order = $3, updated_at = now()
		WHERE id = $4 AND course_id = $5
		RETURNING course_id, updated_at
	`, item.Question, item.Answer, item.SortOrder, item.ID, item.CourseID).Scan(&item.CourseID, &item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *CourseFAQRepository) Delete(ctx context.Context, courseID, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM course_faq_items WHERE id = $1 AND course_id = $2`, id, courseID)
	return checkDeleted(tag, err)
}
