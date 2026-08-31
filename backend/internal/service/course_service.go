package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type CourseRepository interface {
	ListBySectionID(ctx context.Context, sectionID int64) ([]model.Course, error)
	FindBySlug(ctx context.Context, slug string) (model.Course, error)
	Create(ctx context.Context, item model.Course) (model.Course, error)
	Update(ctx context.Context, item model.Course) (model.Course, error)
	Delete(ctx context.Context, sectionID, id int64) error
}

type CourseService struct {
	repo CourseRepository
}

func NewCourseService(repo CourseRepository) *CourseService {
	return &CourseService{repo: repo}
}

func (s *CourseService) ListBySectionID(ctx context.Context, sectionID int64) ([]model.Course, error) {
	return s.repo.ListBySectionID(ctx, sectionID)
}

func (s *CourseService) GetBySlug(ctx context.Context, slug string) (model.Course, error) {
	return s.repo.FindBySlug(ctx, slug)
}

func (s *CourseService) Create(ctx context.Context, item model.Course) (model.Course, error) {
	item.Slug = strings.TrimSpace(item.Slug)
	item.Name = strings.TrimSpace(item.Name)
	if item.SectionID == 0 {
		return model.Course{}, errors.Join(ErrValidation, errors.New("sectionId is required"))
	}
	if item.Slug == "" || item.Name == "" {
		return model.Course{}, errors.Join(ErrValidation, errors.New("slug and name are required"))
	}
	style, err := validateDisplayStyle(item.DisplayStyle)
	if err != nil {
		return model.Course{}, err
	}
	item.DisplayStyle = style

	return s.repo.Create(ctx, item)
}

func (s *CourseService) Update(ctx context.Context, item model.Course) (model.Course, error) {
	item.Slug = strings.TrimSpace(item.Slug)
	item.Name = strings.TrimSpace(item.Name)
	if item.ID == 0 {
		return model.Course{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	// Load-bearing, not just documentation: the repository matches on
	// id AND section_id, so a wrong/missing sectionId here means the update
	// silently (well, loudly — ErrNotFound) touches nothing rather than
	// accidentally hitting a same-id course under a different section.
	if item.SectionID == 0 {
		return model.Course{}, errors.Join(ErrValidation, errors.New("sectionId is required"))
	}
	if item.Slug == "" || item.Name == "" {
		return model.Course{}, errors.Join(ErrValidation, errors.New("slug and name are required"))
	}
	style, err := validateDisplayStyle(item.DisplayStyle)
	if err != nil {
		return model.Course{}, err
	}
	item.DisplayStyle = style

	return s.repo.Update(ctx, item)
}

func (s *CourseService) Delete(ctx context.Context, sectionID, id int64) error {
	if sectionID == 0 || id == 0 {
		return errors.Join(ErrValidation, errors.New("sectionId and id are required"))
	}
	return s.repo.Delete(ctx, sectionID, id)
}
