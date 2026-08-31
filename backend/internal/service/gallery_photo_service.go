package service

import (
	"context"
	"errors"

	"floway-backend/internal/model"
)

type GalleryPhotoRepository interface {
	List(ctx context.Context) ([]model.GalleryPhoto, error)
	Create(ctx context.Context, item model.GalleryPhoto) (model.GalleryPhoto, error)
	Update(ctx context.Context, item model.GalleryPhoto) (model.GalleryPhoto, error)
	Delete(ctx context.Context, id int64) error
}

type GalleryPhotoService struct {
	repo GalleryPhotoRepository
}

func NewGalleryPhotoService(repo GalleryPhotoRepository) *GalleryPhotoService {
	return &GalleryPhotoService{repo: repo}
}

func (s *GalleryPhotoService) List(ctx context.Context) ([]model.GalleryPhoto, error) {
	return s.repo.List(ctx)
}

func (s *GalleryPhotoService) Create(ctx context.Context, item model.GalleryPhoto) (model.GalleryPhoto, error) {
	if item.Image == "" {
		return model.GalleryPhoto{}, errors.Join(ErrValidation, errors.New("image is required"))
	}

	return s.repo.Create(ctx, item)
}

func (s *GalleryPhotoService) Update(ctx context.Context, item model.GalleryPhoto) (model.GalleryPhoto, error) {
	if item.ID == 0 {
		return model.GalleryPhoto{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if item.Image == "" {
		return model.GalleryPhoto{}, errors.Join(ErrValidation, errors.New("image is required"))
	}

	return s.repo.Update(ctx, item)
}

func (s *GalleryPhotoService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
