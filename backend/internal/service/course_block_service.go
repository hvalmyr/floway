package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type CourseBlockRepository interface {
	ListByCourseID(ctx context.Context, courseID int64) ([]model.CourseBlock, error)
	Create(ctx context.Context, item model.CourseBlock) (model.CourseBlock, error)
	Update(ctx context.Context, item model.CourseBlock) (model.CourseBlock, error)
	Delete(ctx context.Context, id int64) error
}

type CourseBlockService struct {
	repo CourseBlockRepository
}

func NewCourseBlockService(repo CourseBlockRepository) *CourseBlockService {
	return &CourseBlockService{repo: repo}
}

func (s *CourseBlockService) ListByCourseID(ctx context.Context, courseID int64) ([]model.CourseBlock, error) {
	return s.repo.ListByCourseID(ctx, courseID)
}

func (s *CourseBlockService) Create(ctx context.Context, item model.CourseBlock) (model.CourseBlock, error) {
	item.Title = strings.TrimSpace(item.Title)
	if item.CourseID == 0 {
		return model.CourseBlock{}, errors.Join(ErrValidation, errors.New("courseId is required"))
	}
	if item.Title == "" {
		return model.CourseBlock{}, errors.Join(ErrValidation, errors.New("title is required"))
	}
	if err := validateOldPrice(item); err != nil {
		return model.CourseBlock{}, err
	}

	return s.repo.Create(ctx, item)
}

func (s *CourseBlockService) Update(ctx context.Context, item model.CourseBlock) (model.CourseBlock, error) {
	item.Title = strings.TrimSpace(item.Title)
	if item.ID == 0 {
		return model.CourseBlock{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if item.Title == "" {
		return model.CourseBlock{}, errors.Join(ErrValidation, errors.New("title is required"))
	}
	if err := validateOldPrice(item); err != nil {
		return model.CourseBlock{}, err
	}

	return s.repo.Update(ctx, item)
}

// validateOldPrice keeps "old price" meaningful as a was/now discount label:
// if set, it has to be strictly higher than the current price.
func validateOldPrice(item model.CourseBlock) error {
	if item.OldPrice != nil && *item.OldPrice <= item.Price {
		return errors.Join(ErrValidation, errors.New("oldPrice must be greater than price"))
	}
	return nil
}

func (s *CourseBlockService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
