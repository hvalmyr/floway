package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type fakeLandingPageBlockRepository struct {
	items  []model.LandingPageBlock
	nextID int64
}

func newFakeLandingPageBlockRepository() *fakeLandingPageBlockRepository {
	return &fakeLandingPageBlockRepository{nextID: 1}
}

func (f *fakeLandingPageBlockRepository) ListByLandingPageID(ctx context.Context, landingPageID int64) ([]model.LandingPageBlock, error) {
	var items []model.LandingPageBlock
	for _, item := range f.items {
		if item.LandingPageID == landingPageID {
			items = append(items, item)
		}
	}
	return items, nil
}

func (f *fakeLandingPageBlockRepository) Create(ctx context.Context, item model.LandingPageBlock) (model.LandingPageBlock, error) {
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeLandingPageBlockRepository) Update(ctx context.Context, item model.LandingPageBlock) (model.LandingPageBlock, error) {
	for i, existing := range f.items {
		if existing.ID == item.ID && existing.LandingPageID == item.LandingPageID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.LandingPageBlock{}, service.ErrNotFound
}

func (f *fakeLandingPageBlockRepository) Delete(ctx context.Context, landingPageID, id int64) error {
	for i, existing := range f.items {
		if existing.ID == id && existing.LandingPageID == landingPageID {
			f.items = append(f.items[:i], f.items[i+1:]...)
			return nil
		}
	}
	return service.ErrNotFound
}

func TestLandingPageBlockService_Create(t *testing.T) {
	t.Run("creates a valid block", func(t *testing.T) {
		repo := newFakeLandingPageBlockRepository()
		svc := service.NewLandingPageBlockService(repo)

		item, err := svc.Create(context.Background(), model.LandingPageBlock{LandingPageID: 1, Image: "block/1.jpg"})

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
	})

	t.Run("rejects a missing landing page id", func(t *testing.T) {
		repo := newFakeLandingPageBlockRepository()
		svc := service.NewLandingPageBlockService(repo)

		_, err := svc.Create(context.Background(), model.LandingPageBlock{Image: "block/1.jpg"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestLandingPageBlockService_Update(t *testing.T) {
	repo := newFakeLandingPageBlockRepository()
	svc := service.NewLandingPageBlockService(repo)
	created, err := svc.Create(context.Background(), model.LandingPageBlock{LandingPageID: 1, Image: "block/1.jpg"})
	require.NoError(t, err)

	t.Run("rejects an update crossing into another landing page", func(t *testing.T) {
		wrong := created
		wrong.LandingPageID = 2

		_, err := svc.Update(context.Background(), wrong)

		require.Error(t, err)
	})

	t.Run("updates within the same landing page", func(t *testing.T) {
		created.Text = "новый текст"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "новый текст", updated.Text)
	})
}

func TestLandingPageBlockService_Delete(t *testing.T) {
	repo := newFakeLandingPageBlockRepository()
	svc := service.NewLandingPageBlockService(repo)
	created, err := svc.Create(context.Background(), model.LandingPageBlock{LandingPageID: 1, Image: "block/1.jpg"})
	require.NoError(t, err)

	t.Run("deletes an existing block", func(t *testing.T) {
		require.NoError(t, svc.Delete(context.Background(), created.LandingPageID, created.ID))
	})

	t.Run("rejects zero ids", func(t *testing.T) {
		err := svc.Delete(context.Background(), 0, 0)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}
