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

type fakePageContentRepository struct {
	items map[string]model.PageContent
	err   error
}

func newFakePageContentRepository() *fakePageContentRepository {
	return &fakePageContentRepository{
		items: map[string]model.PageContent{
			"home_hero_title": {Key: "home_hero_title", Label: "Заголовок", Value: "Старое значение"},
			"home_hero_image": {Key: "home_hero_image", Label: "Фото", Type: "image", Value: "/uploads/old-key.png"},
		},
	}
}

func (f *fakePageContentRepository) List(ctx context.Context) ([]model.PageContent, error) {
	if f.err != nil {
		return nil, f.err
	}
	items := make([]model.PageContent, 0, len(f.items))
	for _, item := range f.items {
		items = append(items, item)
	}
	return items, nil
}

func (f *fakePageContentRepository) Update(ctx context.Context, key, value string) (model.PageContent, string, error) {
	if f.err != nil {
		return model.PageContent{}, "", f.err
	}
	item, ok := f.items[key]
	if !ok {
		return model.PageContent{}, "", service.ErrNotFound
	}
	previousValue := item.Value
	item.Value = value
	f.items[key] = item
	return item, previousValue, nil
}

type fakeImageStorage struct {
	deletedKeys []string
	err         error
}

func (f *fakeImageStorage) Delete(ctx context.Context, key string) error {
	if f.err != nil {
		return f.err
	}
	f.deletedKeys = append(f.deletedKeys, key)
	return nil
}

func TestPageContentService_List(t *testing.T) {
	repo := newFakePageContentRepository()
	svc := service.NewPageContentService(repo, &fakeImageStorage{})

	items, err := svc.List(context.Background())

	require.NoError(t, err)
	assert.Len(t, items, 2)
}

func TestPageContentService_Update(t *testing.T) {
	t.Run("updates an existing key", func(t *testing.T) {
		repo := newFakePageContentRepository()
		svc := service.NewPageContentService(repo, &fakeImageStorage{})

		item, err := svc.Update(context.Background(), "home_hero_title", "Новое значение")

		require.NoError(t, err)
		assert.Equal(t, "Новое значение", item.Value)
	})

	t.Run("returns ErrNotFound for an unknown key", func(t *testing.T) {
		repo := newFakePageContentRepository()
		svc := service.NewPageContentService(repo, &fakeImageStorage{})

		_, err := svc.Update(context.Background(), "does-not-exist", "value")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrNotFound)
	})

	t.Run("rejects an empty key", func(t *testing.T) {
		repo := newFakePageContentRepository()
		svc := service.NewPageContentService(repo, &fakeImageStorage{})

		_, err := svc.Update(context.Background(), "   ", "value")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestPageContentService_Update_ImageCleanup(t *testing.T) {
	t.Run("deletes the old object when an image value is replaced", func(t *testing.T) {
		repo := newFakePageContentRepository()
		storage := &fakeImageStorage{}
		svc := service.NewPageContentService(repo, storage)

		_, err := svc.Update(context.Background(), "home_hero_image", "/uploads/new-key.png")

		require.NoError(t, err)
		assert.Equal(t, []string{"old-key.png"}, storage.deletedKeys)
	})

	t.Run("does not delete anything for a text key", func(t *testing.T) {
		repo := newFakePageContentRepository()
		storage := &fakeImageStorage{}
		svc := service.NewPageContentService(repo, storage)

		_, err := svc.Update(context.Background(), "home_hero_title", "/uploads/looks-like-a-key.png")

		require.NoError(t, err)
		assert.Empty(t, storage.deletedKeys)
	})

	t.Run("does not delete when the value did not actually change", func(t *testing.T) {
		repo := newFakePageContentRepository()
		storage := &fakeImageStorage{}
		svc := service.NewPageContentService(repo, storage)

		_, err := svc.Update(context.Background(), "home_hero_image", "/uploads/old-key.png")

		require.NoError(t, err)
		assert.Empty(t, storage.deletedKeys)
	})

	t.Run("does not delete when there was no previous image", func(t *testing.T) {
		repo := newFakePageContentRepository()
		repo.items["home_hero_image"] = model.PageContent{Key: "home_hero_image", Type: "image", Value: ""}
		storage := &fakeImageStorage{}
		svc := service.NewPageContentService(repo, storage)

		_, err := svc.Update(context.Background(), "home_hero_image", "/uploads/new-key.png")

		require.NoError(t, err)
		assert.Empty(t, storage.deletedKeys)
	})

	t.Run("a cleanup failure does not fail the update", func(t *testing.T) {
		repo := newFakePageContentRepository()
		storage := &fakeImageStorage{err: errors.New("garage unreachable")}
		svc := service.NewPageContentService(repo, storage)

		item, err := svc.Update(context.Background(), "home_hero_image", "/uploads/new-key.png")

		require.NoError(t, err)
		assert.Equal(t, "/uploads/new-key.png", item.Value)
	})
}
