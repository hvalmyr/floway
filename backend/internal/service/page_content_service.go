package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type PageContentRepository interface {
	List(ctx context.Context) ([]model.PageContent, error)
	Update(ctx context.Context, key, value string) (model.PageContent, error)
}

type PageContentService struct {
	repo PageContentRepository
}

func NewPageContentService(repo PageContentRepository) *PageContentService {
	return &PageContentService{repo: repo}
}

func (s *PageContentService) List(ctx context.Context) ([]model.PageContent, error) {
	return s.repo.List(ctx)
}

func (s *PageContentService) Update(ctx context.Context, key, value string) (model.PageContent, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return model.PageContent{}, errors.Join(ErrValidation, errors.New("key is required"))
	}

	return s.repo.Update(ctx, key, value)
}
