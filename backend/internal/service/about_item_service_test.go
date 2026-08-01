package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type fakeAboutItemRepository struct {
	items  []model.AboutItem
	nextID int64
	err    error
}

func newFakeAboutItemRepository() *fakeAboutItemRepository {
	return &fakeAboutItemRepository{nextID: 1}
}

func (f *fakeAboutItemRepository) List(ctx context.Context) ([]model.AboutItem, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func (f *fakeAboutItemRepository) Create(ctx context.Context, item model.AboutItem) (model.AboutItem, error) {
	if f.err != nil {
		return model.AboutItem{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeAboutItemRepository) Update(ctx context.Context, item model.AboutItem) (model.AboutItem, error) {
	if f.err != nil {
		return model.AboutItem{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.AboutItem{}, errors.New("not found")
}

func (f *fakeAboutItemRepository) Delete(ctx context.Context, id int64) error {
	if f.err != nil {
		return f.err
	}
	for i, existing := range f.items {
		if existing.ID == id {
			f.items = append(f.items[:i], f.items[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func TestAboutItemService_Create(t *testing.T) {
	t.Run("creates a valid item", func(t *testing.T) {
		repo := newFakeAboutItemRepository()
		svc := service.NewAboutItemService(repo)

		item, err := svc.Create(context.Background(), model.AboutItem{Badge: "  бейдж  ", Description: "  текст  "})

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
		assert.Equal(t, "бейдж", item.Badge)
		assert.Equal(t, "текст", item.Description)
	})

	t.Run("rejects an empty badge", func(t *testing.T) {
		repo := newFakeAboutItemRepository()
		svc := service.NewAboutItemService(repo)

		_, err := svc.Create(context.Background(), model.AboutItem{Badge: "  ", Description: "текст"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})
}

func TestAboutItemService_Update(t *testing.T) {
	repo := newFakeAboutItemRepository()
	svc := service.NewAboutItemService(repo)
	created, err := svc.Create(context.Background(), model.AboutItem{Badge: "бейдж", Description: "текст"})
	require.NoError(t, err)

	t.Run("updates an existing item", func(t *testing.T) {
		created.Badge = "новый бейдж"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "новый бейдж", updated.Badge)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.AboutItem{Badge: "b", Description: "d"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestAboutItemService_Delete(t *testing.T) {
	repo := newFakeAboutItemRepository()
	svc := service.NewAboutItemService(repo)
	created, err := svc.Create(context.Background(), model.AboutItem{Badge: "b", Description: "d"})
	require.NoError(t, err)

	t.Run("deletes an existing item", func(t *testing.T) {
		require.NoError(t, svc.Delete(context.Background(), created.ID))
		items, err := svc.List(context.Background())
		require.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("rejects a zero id", func(t *testing.T) {
		err := svc.Delete(context.Background(), 0)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}
