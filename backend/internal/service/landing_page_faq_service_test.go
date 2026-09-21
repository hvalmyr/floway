package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type fakeLandingPageFAQRepository struct {
	items  []model.LandingPageFAQItem
	nextID int64
}

func newFakeLandingPageFAQRepository() *fakeLandingPageFAQRepository {
	return &fakeLandingPageFAQRepository{nextID: 1}
}

func (f *fakeLandingPageFAQRepository) ListByLandingPageID(ctx context.Context, landingPageID int64) ([]model.LandingPageFAQItem, error) {
	var items []model.LandingPageFAQItem
	for _, item := range f.items {
		if item.LandingPageID == landingPageID {
			items = append(items, item)
		}
	}
	return items, nil
}

func (f *fakeLandingPageFAQRepository) Create(ctx context.Context, item model.LandingPageFAQItem) (model.LandingPageFAQItem, error) {
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeLandingPageFAQRepository) Update(ctx context.Context, item model.LandingPageFAQItem) (model.LandingPageFAQItem, error) {
	for i, existing := range f.items {
		if existing.ID == item.ID && existing.LandingPageID == item.LandingPageID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.LandingPageFAQItem{}, service.ErrNotFound
}

func (f *fakeLandingPageFAQRepository) Delete(ctx context.Context, landingPageID, id int64) error {
	for i, existing := range f.items {
		if existing.ID == id && existing.LandingPageID == landingPageID {
			f.items = append(f.items[:i], f.items[i+1:]...)
			return nil
		}
	}
	return service.ErrNotFound
}

func TestLandingPageFAQService_Create(t *testing.T) {
	t.Run("creates a valid item", func(t *testing.T) {
		repo := newFakeLandingPageFAQRepository()
		svc := service.NewLandingPageFAQService(repo)

		item, err := svc.Create(context.Background(), model.LandingPageFAQItem{LandingPageID: 1, Question: "Сколько это стоит?", Answer: "Считаем индивидуально."})

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
	})

	t.Run("rejects a missing landing page id", func(t *testing.T) {
		repo := newFakeLandingPageFAQRepository()
		svc := service.NewLandingPageFAQService(repo)

		_, err := svc.Create(context.Background(), model.LandingPageFAQItem{Question: "q", Answer: "a"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects an empty question or answer", func(t *testing.T) {
		repo := newFakeLandingPageFAQRepository()
		svc := service.NewLandingPageFAQService(repo)

		_, err := svc.Create(context.Background(), model.LandingPageFAQItem{LandingPageID: 1, Question: "  ", Answer: "a"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestLandingPageFAQService_Delete(t *testing.T) {
	repo := newFakeLandingPageFAQRepository()
	svc := service.NewLandingPageFAQService(repo)
	created, err := svc.Create(context.Background(), model.LandingPageFAQItem{LandingPageID: 1, Question: "q", Answer: "a"})
	require.NoError(t, err)

	t.Run("deletes an existing item", func(t *testing.T) {
		require.NoError(t, svc.Delete(context.Background(), created.LandingPageID, created.ID))
	})

	t.Run("rejects zero ids", func(t *testing.T) {
		err := svc.Delete(context.Background(), 0, 0)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}
