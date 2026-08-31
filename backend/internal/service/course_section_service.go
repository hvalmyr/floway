package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type CourseSectionRepository interface {
	List(ctx context.Context) ([]model.CourseSection, error)
	Create(ctx context.Context, item model.CourseSection) (model.CourseSection, error)
	Update(ctx context.Context, item model.CourseSection) (model.CourseSection, error)
	Delete(ctx context.Context, id int64) error
}

type CourseSectionService struct {
	repo CourseSectionRepository
}

func NewCourseSectionService(repo CourseSectionRepository) *CourseSectionService {
	return &CourseSectionService{repo: repo}
}

func (s *CourseSectionService) List(ctx context.Context) ([]model.CourseSection, error) {
	return s.repo.List(ctx)
}

func (s *CourseSectionService) Create(ctx context.Context, item model.CourseSection) (model.CourseSection, error) {
	item.Heading = strings.TrimSpace(item.Heading)
	if item.Heading == "" {
		return model.CourseSection{}, errors.Join(ErrValidation, errors.New("heading is required"))
	}

	return s.repo.Create(ctx, item)
}

func (s *CourseSectionService) Update(ctx context.Context, item model.CourseSection) (model.CourseSection, error) {
	item.Heading = strings.TrimSpace(item.Heading)
	if item.ID == 0 {
		return model.CourseSection{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if item.Heading == "" {
		return model.CourseSection{}, errors.Join(ErrValidation, errors.New("heading is required"))
	}

	return s.repo.Update(ctx, item)
}

func (s *CourseSectionService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
