package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type CourseRepository struct {
	db *pgxpool.Pool
}

func NewCourseRepository(db *pgxpool.Pool) *CourseRepository {
	return &CourseRepository{db: db}
}

const courseColumns = "id, section_id, slug, name, description, cover_image, lesson_count, time_length, price, display_style, visible, sort_order, single_card, created_at, updated_at"

func scanCourse(row pgx.Row) (model.Course, error) {
	var item model.Course
	err := row.Scan(&item.ID, &item.SectionID, &item.Slug, &item.Name, &item.Description, &item.CoverImage, &item.LessonCount, &item.TimeLength, &item.Price, &item.DisplayStyle, &item.Visible, &item.SortOrder, &item.SingleCard, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

// ListBySectionID is used by the admin nested "courses in this section" page.
func (r *CourseRepository) ListBySectionID(ctx context.Context, sectionID int64) ([]model.Course, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+courseColumns+`
		FROM courses
		WHERE section_id = $1
		ORDER BY sort_order, id
	`, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Course{}
	for rows.Next() {
		item, err := scanCourse(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// ListBySectionIDs batches courses for several sections in one query — used
// by CourseCatalogService.ListSections to avoid a per-section round trip
// (same reasoning as LessonRepository.ListByCourseBlockIDs).
func (r *CourseRepository) ListBySectionIDs(ctx context.Context, sectionIDs []int64) ([]model.Course, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+courseColumns+`
		FROM courses
		WHERE section_id = ANY($1)
		ORDER BY section_id, sort_order, id
	`, sectionIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Course{}
	for rows.Next() {
		item, err := scanCourse(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *CourseRepository) FindBySlug(ctx context.Context, slug string) (model.Course, error) {
	row := r.db.QueryRow(ctx, `
		SELECT `+courseColumns+`
		FROM courses
		WHERE slug = $1
	`, slug)
	item, err := scanCourse(row)
	return item, translateNotFound(err)
}

func (r *CourseRepository) Create(ctx context.Context, item model.Course) (model.Course, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO courses (section_id, slug, name, description, cover_image, lesson_count, time_length, price, display_style, visible, sort_order, single_card)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`,
		item.SectionID,
		item.Slug,
		item.Name,
		item.Description,
		item.CoverImage,
		item.LessonCount,
		item.TimeLength,
		item.Price,
		item.DisplayStyle,
		item.Visible,
		item.SortOrder,
		item.SingleCard,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

// Update matches on id AND section_id — a URL like
// /course-sections/1/courses/42 can't touch a course that actually belongs
// to a different section (mirrors CourseBlockRepository.Update).
func (r *CourseRepository) Update(ctx context.Context, item model.Course) (model.Course, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE courses
		SET slug = $1, name = $2, description = $3, cover_image = $4, lesson_count = $5, time_length = $6,
		    price = $7, display_style = $8, visible = $9, sort_order = $10, single_card = $11, updated_at = now()
		WHERE id = $12 AND section_id = $13
		RETURNING section_id, updated_at
	`,
		item.Slug,
		item.Name,
		item.Description,
		item.CoverImage,
		item.LessonCount,
		item.TimeLength,
		item.Price,
		item.DisplayStyle,
		item.Visible,
		item.SortOrder,
		item.SingleCard,
		item.ID,
		item.SectionID,
	).Scan(&item.SectionID, &item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *CourseRepository) Delete(ctx context.Context, sectionID, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM courses WHERE id = $1 AND section_id = $2`, id, sectionID)
	return checkDeleted(tag, err)
}
