package service

import (
	"context"
	"errors"

	"floway-backend/internal/model"
)

type LandingPageBlockRepository interface {
	ListByLandingPageID(ctx context.Context, landingPageID int64) ([]model.LandingPageBlock, error)
	Create(ctx context.Context, item model.LandingPageBlock) (model.LandingPageBlock, error)
	Update(ctx context.Context, item model.LandingPageBlock) (model.LandingPageBlock, error)
	Delete(ctx context.Context, landingPageID, id int64) error
}

type LandingPageBlockService struct {
	repo LandingPageBlockRepository
}

func NewLandingPageBlockService(repo LandingPageBlockRepository) *LandingPageBlockService {
	return &LandingPageBlockService{repo: repo}
}

func (s *LandingPageBlockService) ListByLandingPageID(ctx context.Context, landingPageID int64) ([]model.LandingPageBlock, error) {
	return s.repo.ListByLandingPageID(ctx, landingPageID)
}

func (s *LandingPageBlockService) Create(ctx context.Context, item model.LandingPageBlock) (model.LandingPageBlock, error) {
	if item.LandingPageID == 0 {
		return model.LandingPageBlock{}, errors.Join(ErrValidation, errors.New("landingPageId is required"))
	}
	return s.repo.Create(ctx, item)
}

func (s *LandingPageBlockService) Update(ctx context.Context, item model.LandingPageBlock) (model.LandingPageBlock, error) {
	if item.ID == 0 {
		return model.LandingPageBlock{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	// Load-bearing, not just documentation: the repository matches on id AND
	// landing_page_id, so a wrong/missing landingPageId here means the update
	// silently (well, loudly — ErrNotFound) touches nothing rather than
	// accidentally hitting a same-id block under a different landing page.
	if item.LandingPageID == 0 {
		return model.LandingPageBlock{}, errors.Join(ErrValidation, errors.New("landingPageId is required"))
	}
	return s.repo.Update(ctx, item)
}

func (s *LandingPageBlockService) Delete(ctx context.Context, landingPageID, id int64) error {
	if landingPageID == 0 || id == 0 {
		return errors.Join(ErrValidation, errors.New("landingPageId and id are required"))
	}
	return s.repo.Delete(ctx, landingPageID, id)
}
