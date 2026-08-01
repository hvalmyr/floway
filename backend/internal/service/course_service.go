package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type CourseRepository interface {
	List(ctx context.Context) ([]model.Course, error)
	FindBySlug(ctx context.Context, slug string) (model.Course, error)
	Create(ctx context.Context, item model.Course) (model.Course, error)
	Update(ctx context.Context, item model.Course) (model.Course, error)
	Delete(ctx context.Context, id int64) error
}

type CourseService struct {
	repo CourseRepository
}

func NewCourseService(repo CourseRepository) *CourseService {
	return &CourseService{repo: repo}
}

func (s *CourseService) List(ctx context.Context) ([]model.Course, error) {
	return s.repo.List(ctx)
}

func (s *CourseService) GetBySlug(ctx context.Context, slug string) (model.Course, error) {
	return s.repo.FindBySlug(ctx, slug)
}

func (s *CourseService) Create(ctx context.Context, item model.Course) (model.Course, error) {
	item.Slug = strings.TrimSpace(item.Slug)
	item.Title = strings.TrimSpace(item.Title)
	if item.Slug == "" || item.Title == "" {
		return model.Course{}, errors.Join(ErrValidation, errors.New("slug and title are required"))
	}

	if item.Status == "" {
		item.Status = model.CourseStatusActive
	}
	if item.Status != model.CourseStatusActive && item.Status != model.CourseStatusArchived {
		return model.Course{}, errors.Join(ErrValidation, errors.New("status must be active or archived"))
	}
	// gallery is NOT NULL DEFAULT '{}' in the DB, but a nil Go slice (an
	// omitted or explicit-null "gallery" in the request JSON) is sent as
	// SQL NULL, not "use the column default" — violates the constraint.
	if item.Gallery == nil {
		item.Gallery = []string{}
	}

	return s.repo.Create(ctx, item)
}

func (s *CourseService) Update(ctx context.Context, item model.Course) (model.Course, error) {
	item.Slug = strings.TrimSpace(item.Slug)
	item.Title = strings.TrimSpace(item.Title)
	if item.ID == 0 {
		return model.Course{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if item.Slug == "" || item.Title == "" {
		return model.Course{}, errors.Join(ErrValidation, errors.New("slug and title are required"))
	}

	if item.Status == "" {
		item.Status = model.CourseStatusActive
	}
	if item.Status != model.CourseStatusActive && item.Status != model.CourseStatusArchived {
		return model.Course{}, errors.Join(ErrValidation, errors.New("status must be active or archived"))
	}
	if item.Gallery == nil {
		item.Gallery = []string{}
	}

	return s.repo.Update(ctx, item)
}

func (s *CourseService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
