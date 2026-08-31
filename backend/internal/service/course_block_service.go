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

var validDisplayStyles = map[model.CourseBlockDisplayStyle]bool{
	model.DisplayStyleBlueBeige:  true,
	model.DisplayStyleBrownBeige: true,
	model.DisplayStyleBeigeBlue:  true,
	model.DisplayStyleBeigeBrown: true,
}

// validateDisplayStyle defaults an empty value (a client that doesn't send
// the field yet) to the same default the DB column has, and rejects
// anything outside the 4 combos CourseCard.vue knows how to render.
func validateDisplayStyle(style model.CourseBlockDisplayStyle) (model.CourseBlockDisplayStyle, error) {
	if style == "" {
		return model.DisplayStyleBlueBeige, nil
	}
	if !validDisplayStyles[style] {
		return "", errors.Join(ErrValidation, errors.New("displayStyle must be one of blue-beige, brown-beige, beige-blue, beige-brown"))
	}
	return style, nil
}

// BlockName isn't required — a course with a single, undivided block (the
// common case; see CourseBlock's doc comment on the frontend rendering
// rules) leaves it blank, since the public pages only show a block's name
// as a label when the course actually has more than one block.
func (s *CourseBlockService) Create(ctx context.Context, item model.CourseBlock) (model.CourseBlock, error) {
	item.BlockName = strings.TrimSpace(item.BlockName)
	if item.CourseID == 0 {
		return model.CourseBlock{}, errors.Join(ErrValidation, errors.New("courseId is required"))
	}
	style, err := validateDisplayStyle(item.DisplayStyle)
	if err != nil {
		return model.CourseBlock{}, err
	}
	item.DisplayStyle = style

	return s.repo.Create(ctx, item)
}

func (s *CourseBlockService) Update(ctx context.Context, item model.CourseBlock) (model.CourseBlock, error) {
	item.BlockName = strings.TrimSpace(item.BlockName)
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
	style, err := validateDisplayStyle(item.DisplayStyle)
	if err != nil {
		return model.CourseBlock{}, err
	}
	item.DisplayStyle = style

	return s.repo.Update(ctx, item)
}

func (s *CourseBlockService) Delete(ctx context.Context, courseID, id int64) error {
	if courseID == 0 || id == 0 {
		return errors.Join(ErrValidation, errors.New("courseId and id are required"))
	}
	return s.repo.Delete(ctx, courseID, id)
}
