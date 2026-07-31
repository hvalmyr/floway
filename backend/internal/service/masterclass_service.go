package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type MasterclassRepository interface {
	List(ctx context.Context) ([]model.Masterclass, error)
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

func (s *MasterclassService) Create(ctx context.Context, item model.Masterclass) (model.Masterclass, error) {
	item.Slug = strings.TrimSpace(item.Slug)
	item.Title = strings.TrimSpace(item.Title)
	if item.Slug == "" || item.Title == "" {
		return model.Masterclass{}, errors.Join(ErrValidation, errors.New("slug and title are required"))
	}

	if item.Status == "" {
		item.Status = model.MasterclassStatusActive
	}
	if item.Status != model.MasterclassStatusActive && item.Status != model.MasterclassStatusArchived {
		return model.Masterclass{}, errors.Join(ErrValidation, errors.New("status must be active or archived"))
	}

	return s.repo.Create(ctx, item)
}

func (s *MasterclassService) Update(ctx context.Context, item model.Masterclass) (model.Masterclass, error) {
	item.Slug = strings.TrimSpace(item.Slug)
	item.Title = strings.TrimSpace(item.Title)
	if item.ID == 0 {
		return model.Masterclass{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if item.Slug == "" || item.Title == "" {
		return model.Masterclass{}, errors.Join(ErrValidation, errors.New("slug and title are required"))
	}

	if item.Status == "" {
		item.Status = model.MasterclassStatusActive
	}
	if item.Status != model.MasterclassStatusActive && item.Status != model.MasterclassStatusArchived {
		return model.Masterclass{}, errors.Join(ErrValidation, errors.New("status must be active or archived"))
	}

	return s.repo.Update(ctx, item)
}

func (s *MasterclassService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
