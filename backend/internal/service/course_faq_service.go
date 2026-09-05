package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type CourseFAQRepository interface {
	ListByCourseID(ctx context.Context, courseID int64) ([]model.CourseFAQItem, error)
	Create(ctx context.Context, item model.CourseFAQItem) (model.CourseFAQItem, error)
	Update(ctx context.Context, item model.CourseFAQItem) (model.CourseFAQItem, error)
	Delete(ctx context.Context, courseID, id int64) error
}

// CourseFAQService manages the Q&A items of a single course's FAQ block —
// the block's own title/description/visible flag live on Course itself and
// are edited through CourseService, same as every other course field.
type CourseFAQService struct {
	repo CourseFAQRepository
}

func NewCourseFAQService(repo CourseFAQRepository) *CourseFAQService {
	return &CourseFAQService{repo: repo}
}

func (s *CourseFAQService) ListByCourseID(ctx context.Context, courseID int64) ([]model.CourseFAQItem, error) {
	return s.repo.ListByCourseID(ctx, courseID)
}

func (s *CourseFAQService) Create(ctx context.Context, item model.CourseFAQItem) (model.CourseFAQItem, error) {
	item.Question = strings.TrimSpace(item.Question)
	item.Answer = strings.TrimSpace(item.Answer)
	if item.CourseID == 0 {
		return model.CourseFAQItem{}, errors.Join(ErrValidation, errors.New("courseId is required"))
	}
	if item.Question == "" || item.Answer == "" {
		return model.CourseFAQItem{}, errors.Join(ErrValidation, errors.New("question and answer are required"))
	}

	return s.repo.Create(ctx, item)
}

func (s *CourseFAQService) Update(ctx context.Context, item model.CourseFAQItem) (model.CourseFAQItem, error) {
	item.Question = strings.TrimSpace(item.Question)
	item.Answer = strings.TrimSpace(item.Answer)
	if item.ID == 0 {
		return model.CourseFAQItem{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	// Load-bearing, not just documentation: the repository matches on id AND
	// course_id, so a wrong/missing courseId here means the update silently
	// (well, loudly — ErrNotFound) touches nothing rather than accidentally
	// hitting a same-id item under a different course.
	if item.CourseID == 0 {
		return model.CourseFAQItem{}, errors.Join(ErrValidation, errors.New("courseId is required"))
	}
	if item.Question == "" || item.Answer == "" {
		return model.CourseFAQItem{}, errors.Join(ErrValidation, errors.New("question and answer are required"))
	}

	return s.repo.Update(ctx, item)
}

func (s *CourseFAQService) Delete(ctx context.Context, courseID, id int64) error {
	if courseID == 0 || id == 0 {
		return errors.Join(ErrValidation, errors.New("courseId and id are required"))
	}
	return s.repo.Delete(ctx, courseID, id)
}
