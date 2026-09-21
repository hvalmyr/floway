package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type ObjectTypeRepository interface {
	List(ctx context.Context) ([]model.ObjectType, error)
	Create(ctx context.Context, item model.ObjectType) (model.ObjectType, error)
	Update(ctx context.Context, item model.ObjectType) (model.ObjectType, error)
	Delete(ctx context.Context, id int64) error
}

type ObjectTypeService struct {
	repo ObjectTypeRepository
}

func NewObjectTypeService(repo ObjectTypeRepository) *ObjectTypeService {
	return &ObjectTypeService{repo: repo}
}

func (s *ObjectTypeService) List(ctx context.Context) ([]model.ObjectType, error) {
	return s.repo.List(ctx)
}

func (s *ObjectTypeService) Create(ctx context.Context, item model.ObjectType) (model.ObjectType, error) {
	item.Slug = strings.TrimSpace(item.Slug)
	item.Name = strings.TrimSpace(item.Name)
	if item.Slug == "" || item.Name == "" {
		return model.ObjectType{}, errors.Join(ErrValidation, errors.New("slug and name are required"))
	}
	return s.repo.Create(ctx, item)
}

func (s *ObjectTypeService) Update(ctx context.Context, item model.ObjectType) (model.ObjectType, error) {
	item.Slug = strings.TrimSpace(item.Slug)
	item.Name = strings.TrimSpace(item.Name)
	if item.ID == 0 {
		return model.ObjectType{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if item.Slug == "" || item.Name == "" {
		return model.ObjectType{}, errors.Join(ErrValidation, errors.New("slug and name are required"))
	}
	return s.repo.Update(ctx, item)
}

func (s *ObjectTypeService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
