package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type IconRepository interface {
	List(ctx context.Context) ([]model.Icon, error)
	Create(ctx context.Context, item model.Icon) (model.Icon, error)
	Delete(ctx context.Context, id int64) error
}

// IconService is the admin-uploaded SVG icon library — see model.Icon's doc
// comment for the "icon:<id>" value convention that features/page_content
// use to reference one.
type IconService struct {
	repo IconRepository
}

func NewIconService(repo IconRepository) *IconService {
	return &IconService{repo: repo}
}

func (s *IconService) List(ctx context.Context) ([]model.Icon, error) {
	return s.repo.List(ctx)
}

// Create sanitizes the uploaded SVG before it ever reaches the database —
// see sanitizeSVG's doc comment for exactly what survives. This is the
// *only* point of entry for icon content, so every later reader (the public
// GET this feeds, the frontend's v-html render) can trust it unconditionally.
func (s *IconService) Create(ctx context.Context, item model.Icon) (model.Icon, error) {
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return model.Icon{}, errors.Join(ErrValidation, errors.New("name is required"))
	}
	if strings.TrimSpace(item.SVG) == "" {
		return model.Icon{}, errors.Join(ErrValidation, errors.New("svg is required"))
	}

	clean, err := sanitizeSVG(item.SVG)
	if err != nil {
		return model.Icon{}, errors.Join(ErrValidation, errors.New("svg is not a valid SVG file"))
	}
	item.SVG = clean

	return s.repo.Create(ctx, item)
}

func (s *IconService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
