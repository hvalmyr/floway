package service

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"floway-backend/internal/model"
)

type CustomDisplayStyleRepository interface {
	List(ctx context.Context) ([]model.CustomDisplayStyle, error)
	Create(ctx context.Context, item model.CustomDisplayStyle) (model.CustomDisplayStyle, error)
	Update(ctx context.Context, item model.CustomDisplayStyle) (model.CustomDisplayStyle, error)
	Delete(ctx context.Context, id int64) error
}

type CustomDisplayStyleService struct {
	repo CustomDisplayStyleRepository
}

func NewCustomDisplayStyleService(repo CustomDisplayStyleRepository) *CustomDisplayStyleService {
	return &CustomDisplayStyleService{repo: repo}
}

// hexColorPattern matches a plain 6-digit "#rrggbb" hex color — the format
// an <input type="color"> in the admin UI always sends, so no 3-digit or
// alpha variant needs to be accepted.
var hexColorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

func (s *CustomDisplayStyleService) List(ctx context.Context) ([]model.CustomDisplayStyle, error) {
	return s.repo.List(ctx)
}

func (s *CustomDisplayStyleService) validate(item *model.CustomDisplayStyle) error {
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return errors.Join(ErrValidation, errors.New("name is required"))
	}
	if !hexColorPattern.MatchString(item.BgColor) {
		return errors.Join(ErrValidation, errors.New("bgColor must be a hex color like #82b1cc"))
	}
	if !hexColorPattern.MatchString(item.TextColor) {
		return errors.Join(ErrValidation, errors.New("textColor must be a hex color like #82b1cc"))
	}
	return nil
}

func (s *CustomDisplayStyleService) Create(ctx context.Context, item model.CustomDisplayStyle) (model.CustomDisplayStyle, error) {
	if err := s.validate(&item); err != nil {
		return model.CustomDisplayStyle{}, err
	}
	return s.repo.Create(ctx, item)
}

func (s *CustomDisplayStyleService) Update(ctx context.Context, item model.CustomDisplayStyle) (model.CustomDisplayStyle, error) {
	if item.ID == 0 {
		return model.CustomDisplayStyle{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if err := s.validate(&item); err != nil {
		return model.CustomDisplayStyle{}, err
	}
	return s.repo.Update(ctx, item)
}

func (s *CustomDisplayStyleService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
