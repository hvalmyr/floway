package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type MasterclassRepository interface {
	List(ctx context.Context) ([]model.Masterclass, error)
	FindBySlug(ctx context.Context, slug string) (model.Masterclass, error)
	Create(ctx context.Context, item model.Masterclass) (model.Masterclass, error)
	Update(ctx context.Context, item model.Masterclass) (model.Masterclass, error)
	Delete(ctx context.Context, id int64) error
}

type MasterclassService struct {
	repo MasterclassRepository
}

func NewMasterclassService(repo MasterclassRepository) *MasterclassService {
	return &MasterclassService{repo: repo}
}

func (s *MasterclassService) List(ctx context.Context) ([]model.Masterclass, error) {
	return s.repo.List(ctx)
}

func (s *MasterclassService) GetBySlug(ctx context.Context, slug string) (model.Masterclass, error) {
	return s.repo.FindBySlug(ctx, slug)
}

// requiredFields validates the fields the admin brief marks with "*":
// title, slug, cover image, duration, description and price. description2
// and endingText stay optional — not every masterclass needs a second
// paragraph or a closing note.
func requiredMasterclassFields(item model.Masterclass) error {
	if item.Slug == "" || item.Title == "" || item.CoverImage == "" || item.Duration == "" || item.Description == "" || item.Price == "" {
		return errors.Join(ErrValidation, errors.New("slug, title, coverImage, duration, description and price are required"))
	}
	return nil
}

func (s *MasterclassService) Create(ctx context.Context, item model.Masterclass) (model.Masterclass, error) {
	item.Slug = strings.TrimSpace(item.Slug)
	item.Title = strings.TrimSpace(item.Title)
	if item.Status == "" {
		item.Status = model.MasterclassStatusActive
	}
	if item.Status != model.MasterclassStatusActive && item.Status != model.MasterclassStatusArchived {
		return model.Masterclass{}, errors.Join(ErrValidation, errors.New("status must be active or archived"))
	}
	if err := requiredMasterclassFields(item); err != nil {
		return model.Masterclass{}, err
	}

	return s.repo.Create(ctx, item)
}

func (s *MasterclassService) Update(ctx context.Context, item model.Masterclass) (model.Masterclass, error) {
	item.Slug = strings.TrimSpace(item.Slug)
	item.Title = strings.TrimSpace(item.Title)
	if item.ID == 0 {
		return model.Masterclass{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if item.Status == "" {
		item.Status = model.MasterclassStatusActive
	}
	if item.Status != model.MasterclassStatusActive && item.Status != model.MasterclassStatusArchived {
		return model.Masterclass{}, errors.Join(ErrValidation, errors.New("status must be active or archived"))
	}
	if err := requiredMasterclassFields(item); err != nil {
		return model.Masterclass{}, err
	}

	return s.repo.Update(ctx, item)
}

func (s *MasterclassService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
