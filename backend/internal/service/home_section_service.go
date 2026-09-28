package service

import (
	"context"
	"errors"

	"floway-backend/internal/model"
)

type HomeSectionRepository interface {
	List(ctx context.Context) ([]model.HomeSection, error)
	Update(ctx context.Context, item model.HomeSection) (model.HomeSection, error)
}

type HomeSectionService struct {
	repo HomeSectionRepository
}

func NewHomeSectionService(repo HomeSectionRepository) *HomeSectionService {
	return &HomeSectionService{repo: repo}
}

func (s *HomeSectionService) List(ctx context.Context) ([]model.HomeSection, error) {
	return s.repo.List(ctx)
}

func (s *HomeSectionService) Update(ctx context.Context, item model.HomeSection) (model.HomeSection, error) {
	if item.ID == 0 {
		return model.HomeSection{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Update(ctx, item)
}
