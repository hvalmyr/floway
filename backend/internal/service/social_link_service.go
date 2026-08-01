package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type SocialLinkRepository interface {
	List(ctx context.Context) ([]model.SocialLink, error)
	Create(ctx context.Context, item model.SocialLink) (model.SocialLink, error)
	Update(ctx context.Context, item model.SocialLink) (model.SocialLink, error)
	Delete(ctx context.Context, id int64) error
}

type SocialLinkService struct {
	repo SocialLinkRepository
}

func NewSocialLinkService(repo SocialLinkRepository) *SocialLinkService {
	return &SocialLinkService{repo: repo}
}

func (s *SocialLinkService) List(ctx context.Context) ([]model.SocialLink, error) {
	return s.repo.List(ctx)
}

func (s *SocialLinkService) Create(ctx context.Context, item model.SocialLink) (model.SocialLink, error) {
	item.Label = strings.TrimSpace(item.Label)
	item.Href = strings.TrimSpace(item.Href)
	if item.Label == "" || item.Href == "" {
		return model.SocialLink{}, errors.Join(ErrValidation, errors.New("label and href are required"))
	}
	return s.repo.Create(ctx, item)
}

func (s *SocialLinkService) Update(ctx context.Context, item model.SocialLink) (model.SocialLink, error) {
	item.Label = strings.TrimSpace(item.Label)
	item.Href = strings.TrimSpace(item.Href)
	if item.ID == 0 {
		return model.SocialLink{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if item.Label == "" || item.Href == "" {
		return model.SocialLink{}, errors.Join(ErrValidation, errors.New("label and href are required"))
	}
	return s.repo.Update(ctx, item)
}

func (s *SocialLinkService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
