package service

import (
	"context"
	"errors"

	"floway-backend/internal/model"
)

type GiftCertificateCarouselPhotoRepository interface {
	List(ctx context.Context) ([]model.GiftCertificateCarouselPhoto, error)
	Create(ctx context.Context, item model.GiftCertificateCarouselPhoto) (model.GiftCertificateCarouselPhoto, error)
	Update(ctx context.Context, item model.GiftCertificateCarouselPhoto) (model.GiftCertificateCarouselPhoto, error)
	Delete(ctx context.Context, id int64) error
}

type GiftCertificateCarouselPhotoService struct {
	repo GiftCertificateCarouselPhotoRepository
}

func NewGiftCertificateCarouselPhotoService(repo GiftCertificateCarouselPhotoRepository) *GiftCertificateCarouselPhotoService {
	return &GiftCertificateCarouselPhotoService{repo: repo}
}

func (s *GiftCertificateCarouselPhotoService) List(ctx context.Context) ([]model.GiftCertificateCarouselPhoto, error) {
	return s.repo.List(ctx)
}

func (s *GiftCertificateCarouselPhotoService) Create(ctx context.Context, item model.GiftCertificateCarouselPhoto) (model.GiftCertificateCarouselPhoto, error) {
	if item.Image == "" {
		return model.GiftCertificateCarouselPhoto{}, errors.Join(ErrValidation, errors.New("image is required"))
	}

	return s.repo.Create(ctx, item)
}

func (s *GiftCertificateCarouselPhotoService) Update(ctx context.Context, item model.GiftCertificateCarouselPhoto) (model.GiftCertificateCarouselPhoto, error) {
	if item.ID == 0 {
		return model.GiftCertificateCarouselPhoto{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if item.Image == "" {
		return model.GiftCertificateCarouselPhoto{}, errors.Join(ErrValidation, errors.New("image is required"))
	}

	return s.repo.Update(ctx, item)
}

func (s *GiftCertificateCarouselPhotoService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
