package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type LessonRepository interface {
	ListByCourseBlockID(ctx context.Context, courseBlockID int64) ([]model.Lesson, error)
	Create(ctx context.Context, item model.Lesson) (model.Lesson, error)
	Update(ctx context.Context, item model.Lesson) (model.Lesson, error)
	Delete(ctx context.Context, id int64) error
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
	item.Title = strings.TrimSpace(item.Title)
	if item.Title == "" {
		return model.Lesson{}, errors.Join(ErrValidation, errors.New("title is required"))
	}
	if item.CourseBlockID == 0 {
		return model.Lesson{}, errors.Join(ErrValidation, errors.New("courseBlockId is required"))
	}

	return s.repo.Create(ctx, item)
}

func (s *LessonService) Update(ctx context.Context, item model.Lesson) (model.Lesson, error) {
	item.Title = strings.TrimSpace(item.Title)
	if item.ID == 0 {
		return model.Lesson{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if item.Title == "" {
		return model.Lesson{}, errors.Join(ErrValidation, errors.New("title is required"))
	}

	return s.repo.Update(ctx, item)
}

func (s *LessonService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
