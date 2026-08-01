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
	Delete(ctx context.Context, courseID, id int64) error
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
	// Load-bearing, not just documentation: the repository matches on
	// id AND course_id, so a wrong/missing courseId here means the update
	// silently (well, loudly — ErrNotFound) touches nothing rather than
	// accidentally hitting a same-id block under a different course.
	if item.CourseID == 0 {
		return model.CourseBlock{}, errors.Join(ErrValidation, errors.New("courseId is required"))
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

func (s *CourseBlockService) Delete(ctx context.Context, courseID, id int64) error {
	if courseID == 0 || id == 0 {
		return errors.Join(ErrValidation, errors.New("courseId and id are required"))
	}
	return s.repo.Delete(ctx, courseID, id)
}
