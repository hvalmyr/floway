package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type PageContentRepository interface {
	List(ctx context.Context) ([]model.PageContent, error)
	// Update also returns the value the key held before this update, so the
	// caller can clean up an image it's about to become orphaned by (see
	// ImageStorage below) — one round trip, no read/update race.
	Update(ctx context.Context, key, value string) (item model.PageContent, previousValue string, err error)
}

// ImageStorage is the narrow slice of the Garage client PageContentService
// needs: deleting an object that a page_content image field no longer
// references, once it's been replaced by another upload.
type ImageStorage interface {
	Delete(ctx context.Context, key string) error
}

type PageContentService struct {
	repo    PageContentRepository
	storage ImageStorage
}

func NewPageContentService(repo PageContentRepository, storage ImageStorage) *PageContentService {
	return &PageContentService{repo: repo, storage: storage}
}

func (s *PageContentService) List(ctx context.Context) ([]model.PageContent, error) {
	return s.repo.List(ctx)
}

func (s *PageContentService) Update(ctx context.Context, key, value string) (model.PageContent, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return model.PageContent{}, errors.Join(ErrValidation, errors.New("key is required"))
	}

	item, previousValue, err := s.repo.Update(ctx, key, value)
	if err != nil {
		return model.PageContent{}, err
	}

	if item.Type == "image" && previousValue != "" && previousValue != value {
		if oldKey, ok := uploadObjectKey(previousValue); ok {
			// Best-effort: the content update already succeeded, so a failed
			// cleanup (network blip, key already gone) must not fail the
			// request — it just leaves one orphaned object in Garage.
			_ = s.storage.Delete(ctx, oldKey)
		}
	}

	return item, nil
}

// uploadObjectKey extracts the Garage object key from a page_content image
// value, which is always either "" or the "/uploads/{key}" URL that
// uploadHandler.upload returns. Values that don't match this shape (e.g. a
// hand-edited or pre-Garage value) are left alone rather than guessed at.
func uploadObjectKey(value string) (key string, ok bool) {
	const prefix = "/uploads/"
	if !strings.HasPrefix(value, prefix) {
		return "", false
	}
	key = strings.TrimPrefix(value, prefix)
	return key, key != ""
}
