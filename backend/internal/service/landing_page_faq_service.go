package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type LandingPageFAQRepository interface {
	ListByLandingPageID(ctx context.Context, landingPageID int64) ([]model.LandingPageFAQItem, error)
	Create(ctx context.Context, item model.LandingPageFAQItem) (model.LandingPageFAQItem, error)
	Update(ctx context.Context, item model.LandingPageFAQItem) (model.LandingPageFAQItem, error)
	Delete(ctx context.Context, landingPageID, id int64) error
}

// LandingPageFAQService manages the Q&A items of a single landing page's FAQ
// block — the block's own title/description/visible flag live on
// LandingPage itself and are edited through LandingPageService, same as
// CourseFAQService/Course.
type LandingPageFAQService struct {
	repo LandingPageFAQRepository
}

func NewLandingPageFAQService(repo LandingPageFAQRepository) *LandingPageFAQService {
	return &LandingPageFAQService{repo: repo}
}

func (s *LandingPageFAQService) ListByLandingPageID(ctx context.Context, landingPageID int64) ([]model.LandingPageFAQItem, error) {
	return s.repo.ListByLandingPageID(ctx, landingPageID)
}

func (s *LandingPageFAQService) Create(ctx context.Context, item model.LandingPageFAQItem) (model.LandingPageFAQItem, error) {
	item.Question = strings.TrimSpace(item.Question)
	item.Answer = strings.TrimSpace(item.Answer)
	if item.LandingPageID == 0 {
		return model.LandingPageFAQItem{}, errors.Join(ErrValidation, errors.New("landingPageId is required"))
	}
	if item.Question == "" || item.Answer == "" {
		return model.LandingPageFAQItem{}, errors.Join(ErrValidation, errors.New("question and answer are required"))
	}
	return s.repo.Create(ctx, item)
}

func (s *LandingPageFAQService) Update(ctx context.Context, item model.LandingPageFAQItem) (model.LandingPageFAQItem, error) {
	item.Question = strings.TrimSpace(item.Question)
	item.Answer = strings.TrimSpace(item.Answer)
	if item.ID == 0 {
		return model.LandingPageFAQItem{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if item.LandingPageID == 0 {
		return model.LandingPageFAQItem{}, errors.Join(ErrValidation, errors.New("landingPageId is required"))
	}
	if item.Question == "" || item.Answer == "" {
		return model.LandingPageFAQItem{}, errors.Join(ErrValidation, errors.New("question and answer are required"))
	}
	return s.repo.Update(ctx, item)
}

func (s *LandingPageFAQService) Delete(ctx context.Context, landingPageID, id int64) error {
	if landingPageID == 0 || id == 0 {
		return errors.Join(ErrValidation, errors.New("landingPageId and id are required"))
	}
	return s.repo.Delete(ctx, landingPageID, id)
}
