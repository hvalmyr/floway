package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

var validSiteButtonVariants = map[string]bool{
	"primary": true,
	"outline": true,
}

type SiteButtonRepository interface {
	List(ctx context.Context) ([]model.SiteButton, error)
	Update(ctx context.Context, key, text, variant, url string) (model.SiteButton, error)
}

type SiteButtonService struct {
	repo SiteButtonRepository
}

func NewSiteButtonService(repo SiteButtonRepository) *SiteButtonService {
	return &SiteButtonService{repo: repo}
}

func (s *SiteButtonService) List(ctx context.Context) ([]model.SiteButton, error) {
	return s.repo.List(ctx)
}

func (s *SiteButtonService) Update(ctx context.Context, key, text, variant, url string) (model.SiteButton, error) {
	key = strings.TrimSpace(key)
	text = strings.TrimSpace(text)
	url = strings.TrimSpace(url)

	if key == "" {
		return model.SiteButton{}, errors.Join(ErrValidation, errors.New("key is required"))
	}
	if text == "" {
		return model.SiteButton{}, errors.Join(ErrValidation, errors.New("text is required"))
	}
	if url == "" {
		return model.SiteButton{}, errors.Join(ErrValidation, errors.New("url is required"))
	}
	if !validSiteButtonVariants[variant] {
		return model.SiteButton{}, errors.Join(ErrValidation, errors.New("variant must be one of: primary, outline"))
	}

	return s.repo.Update(ctx, key, text, variant, url)
}
