package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type LessonRepository interface {
	ListByCourseBlockID(ctx context.Context, courseBlockID int64) ([]model.Lesson, error)
	ListByCourseID(ctx context.Context, courseID int64) ([]model.Lesson, error)
	Create(ctx context.Context, item model.Lesson) (model.Lesson, error)
	UpdateByCourseBlock(ctx context.Context, item model.Lesson) (model.Lesson, error)
	UpdateByCourse(ctx context.Context, item model.Lesson) (model.Lesson, error)
	DeleteByCourseBlock(ctx context.Context, courseBlockID, id int64) error
	DeleteByCourse(ctx context.Context, courseID, id int64) error
}

type LessonService struct {
	repo LessonRepository
}

func NewLessonService(repo LessonRepository) *LessonService {
	return &LessonService{repo: repo}
}

func (s *LessonService) ListByCourseBlockID(ctx context.Context, courseBlockID int64) ([]model.Lesson, error) {
	return s.repo.ListByCourseBlockID(ctx, courseBlockID)
}

func (s *LessonService) Create(ctx context.Context, item model.Lesson) (model.Lesson, error) {
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return model.Lesson{}, errors.Join(ErrValidation, errors.New("name is required"))
	}
	if item.CourseBlockID == nil || *item.CourseBlockID == 0 {
		return model.Lesson{}, errors.Join(ErrValidation, errors.New("courseBlockId is required"))
	}

	return s.repo.Create(ctx, item)
}

func (s *LessonService) Update(ctx context.Context, item model.Lesson) (model.Lesson, error) {
	item.Name = strings.TrimSpace(item.Name)
	if item.ID == 0 {
		return model.Lesson{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	// Load-bearing: the repository matches on id AND course_block_id, so a
	// wrong/missing courseBlockId here means the update touches nothing
	// (ErrNotFound) instead of silently landing on a same-id lesson under a
	// different block.
	if item.CourseBlockID == nil || *item.CourseBlockID == 0 {
		return model.Lesson{}, errors.Join(ErrValidation, errors.New("courseBlockId is required"))
	}
	if item.Name == "" {
		return model.Lesson{}, errors.Join(ErrValidation, errors.New("name is required"))
	}

	return s.repo.UpdateByCourseBlock(ctx, item)
}

func (s *LessonService) Delete(ctx context.Context, courseBlockID, id int64) error {
	if courseBlockID == 0 || id == 0 {
		return errors.Join(ErrValidation, errors.New("courseBlockId and id are required"))
	}
	return s.repo.DeleteByCourseBlock(ctx, courseBlockID, id)
}

// --- course-without-blocks counterparts (see model.Lesson's doc comment) ---

func (s *LessonService) ListByCourseID(ctx context.Context, courseID int64) ([]model.Lesson, error) {
	return s.repo.ListByCourseID(ctx, courseID)
}

func (s *LessonService) CreateForCourse(ctx context.Context, item model.Lesson) (model.Lesson, error) {
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return model.Lesson{}, errors.Join(ErrValidation, errors.New("name is required"))
	}
	if item.CourseID == nil || *item.CourseID == 0 {
		return model.Lesson{}, errors.Join(ErrValidation, errors.New("courseId is required"))
	}

	return s.repo.Create(ctx, item)
}

func (s *LessonService) UpdateForCourse(ctx context.Context, item model.Lesson) (model.Lesson, error) {
	item.Name = strings.TrimSpace(item.Name)
	if item.ID == 0 {
		return model.Lesson{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if item.CourseID == nil || *item.CourseID == 0 {
		return model.Lesson{}, errors.Join(ErrValidation, errors.New("courseId is required"))
	}
	if item.Name == "" {
		return model.Lesson{}, errors.Join(ErrValidation, errors.New("name is required"))
	}

	return s.repo.UpdateByCourse(ctx, item)
}

func (s *LessonService) DeleteForCourse(ctx context.Context, courseID, id int64) error {
	if courseID == 0 || id == 0 {
		return errors.Join(ErrValidation, errors.New("courseId and id are required"))
	}
	return s.repo.DeleteByCourse(ctx, courseID, id)
}
