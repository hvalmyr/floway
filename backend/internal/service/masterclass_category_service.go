package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type MasterclassCategoryRepository interface {
	List(ctx context.Context) ([]model.MasterclassCategory, error)
	Create(ctx context.Context, item model.MasterclassCategory) (model.MasterclassCategory, error)
	Update(ctx context.Context, item model.MasterclassCategory) (model.MasterclassCategory, error)
	Delete(ctx context.Context, id int64) error
}

type MasterclassCategoryService struct {
	repo MasterclassCategoryRepository
}

func NewMasterclassCategoryService(repo MasterclassCategoryRepository) *MasterclassCategoryService {
	return &MasterclassCategoryService{repo: repo}
}

func (s *MasterclassCategoryService) List(ctx context.Context) ([]model.MasterclassCategory, error) {
	return s.repo.List(ctx)
}

func (s *MasterclassCategoryService) Create(ctx context.Context, item model.MasterclassCategory) (model.MasterclassCategory, error) {
	item.Slug = strings.TrimSpace(item.Slug)
	item.Name = strings.TrimSpace(item.Name)
	if item.Slug == "" || item.Name == "" {
		return model.MasterclassCategory{}, errors.Join(ErrValidation, errors.New("slug and name are required"))
	}
	return s.repo.Create(ctx, item)
}

func (s *MasterclassCategoryService) Update(ctx context.Context, item model.MasterclassCategory) (model.MasterclassCategory, error) {
	item.Slug = strings.TrimSpace(item.Slug)
	item.Name = strings.TrimSpace(item.Name)
	if item.ID == 0 {
		return model.MasterclassCategory{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if item.Slug == "" || item.Name == "" {
		return model.MasterclassCategory{}, errors.Join(ErrValidation, errors.New("slug and name are required"))
	}
	return s.repo.Update(ctx, item)
}

func (s *MasterclassCategoryService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
