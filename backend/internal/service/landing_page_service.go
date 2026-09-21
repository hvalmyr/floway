package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type LandingPageRepository interface {
	List(ctx context.Context) ([]model.LandingPage, error)
	FindBySlug(ctx context.Context, slug string) (model.LandingPage, error)
	Create(ctx context.Context, item model.LandingPage) (model.LandingPage, error)
	Update(ctx context.Context, item model.LandingPage) (model.LandingPage, error)
	Delete(ctx context.Context, id int64) error
}

type LandingPageService struct {
	repo LandingPageRepository
}

func NewLandingPageService(repo LandingPageRepository) *LandingPageService {
	return &LandingPageService{repo: repo}
}

func (s *LandingPageService) List(ctx context.Context) ([]model.LandingPage, error) {
	return s.repo.List(ctx)
}

func (s *LandingPageService) Create(ctx context.Context, item model.LandingPage) (model.LandingPage, error) {
	item.Slug = strings.TrimSpace(item.Slug)
	item.H1 = strings.TrimSpace(item.H1)
	if item.ObjectTypeID == 0 {
		return model.LandingPage{}, errors.Join(ErrValidation, errors.New("objectTypeId is required"))
	}
	if item.Slug == "" || item.H1 == "" {
		return model.LandingPage{}, errors.Join(ErrValidation, errors.New("slug and h1 are required"))
	}
	return s.repo.Create(ctx, item)
}

func (s *LandingPageService) Update(ctx context.Context, item model.LandingPage) (model.LandingPage, error) {
	item.Slug = strings.TrimSpace(item.Slug)
	item.H1 = strings.TrimSpace(item.H1)
	if item.ID == 0 {
		return model.LandingPage{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if item.ObjectTypeID == 0 {
		return model.LandingPage{}, errors.Join(ErrValidation, errors.New("objectTypeId is required"))
	}
	if item.Slug == "" || item.H1 == "" {
		return model.LandingPage{}, errors.Join(ErrValidation, errors.New("slug and h1 are required"))
	}
	return s.repo.Update(ctx, item)
}

func (s *LandingPageService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
