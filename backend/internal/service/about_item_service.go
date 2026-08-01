package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type AboutItemRepository interface {
	List(ctx context.Context) ([]model.AboutItem, error)
	Create(ctx context.Context, item model.AboutItem) (model.AboutItem, error)
	Update(ctx context.Context, item model.AboutItem) (model.AboutItem, error)
	Delete(ctx context.Context, id int64) error
}

type AboutItemService struct {
	repo AboutItemRepository
}

func NewAboutItemService(repo AboutItemRepository) *AboutItemService {
	return &AboutItemService{repo: repo}
}

func (s *AboutItemService) List(ctx context.Context) ([]model.AboutItem, error) {
	return s.repo.List(ctx)
}

func (s *AboutItemService) Create(ctx context.Context, item model.AboutItem) (model.AboutItem, error) {
	item.Badge = strings.TrimSpace(item.Badge)
	item.Description = strings.TrimSpace(item.Description)
	if item.Badge == "" || item.Description == "" {
		return model.AboutItem{}, errors.Join(ErrValidation, errors.New("badge and description are required"))
	}
	return s.repo.Create(ctx, item)
}

func (s *AboutItemService) Update(ctx context.Context, item model.AboutItem) (model.AboutItem, error) {
	item.Badge = strings.TrimSpace(item.Badge)
	item.Description = strings.TrimSpace(item.Description)
	if item.ID == 0 {
		return model.AboutItem{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if item.Badge == "" || item.Description == "" {
		return model.AboutItem{}, errors.Join(ErrValidation, errors.New("badge and description are required"))
	}
	return s.repo.Update(ctx, item)
}

func (s *AboutItemService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
