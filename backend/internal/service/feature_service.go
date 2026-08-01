package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

var validFeaturePages = map[string]bool{
	"home":          true,
	"masterclasses": true,
}

type FeatureRepository interface {
	List(ctx context.Context) ([]model.Feature, error)
	ListByPage(ctx context.Context, page string) ([]model.Feature, error)
	Create(ctx context.Context, item model.Feature) (model.Feature, error)
	Update(ctx context.Context, item model.Feature) (model.Feature, error)
	Delete(ctx context.Context, id int64) error
}

type FeatureService struct {
	repo FeatureRepository
}

func NewFeatureService(repo FeatureRepository) *FeatureService {
	return &FeatureService{repo: repo}
}

func (s *FeatureService) List(ctx context.Context) ([]model.Feature, error) {
	return s.repo.List(ctx)
}

func (s *FeatureService) ListByPage(ctx context.Context, page string) ([]model.Feature, error) {
	return s.repo.ListByPage(ctx, page)
}

func (s *FeatureService) validate(item model.Feature) (model.Feature, error) {
	item.Page = strings.TrimSpace(item.Page)
	item.Icon = strings.TrimSpace(item.Icon)
	item.Title = strings.TrimSpace(item.Title)
	item.Description = strings.TrimSpace(item.Description)

	if !validFeaturePages[item.Page] {
		return item, errors.Join(ErrValidation, errors.New("page must be one of: home, masterclasses"))
	}
	if item.Icon == "" || item.Title == "" || item.Description == "" {
		return item, errors.Join(ErrValidation, errors.New("icon, title and description are required"))
	}
	return item, nil
}

func (s *FeatureService) Create(ctx context.Context, item model.Feature) (model.Feature, error) {
	item, err := s.validate(item)
	if err != nil {
		return model.Feature{}, err
	}
	return s.repo.Create(ctx, item)
}

func (s *FeatureService) Update(ctx context.Context, item model.Feature) (model.Feature, error) {
	if item.ID == 0 {
		return model.Feature{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	item, err := s.validate(item)
	if err != nil {
		return model.Feature{}, err
	}
	return s.repo.Update(ctx, item)
}

func (s *FeatureService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
