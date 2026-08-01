package service_test

import (
	"context"
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

func (f *fakePageContentRepository) Update(ctx context.Context, key, value string) (model.PageContent, error) {
	if f.err != nil {
		return model.PageContent{}, f.err
	}
	item, ok := f.items[key]
	if !ok {
		return model.PageContent{}, service.ErrNotFound
	}
	item.Value = value
	f.items[key] = item
	return item, nil
}

func TestPageContentService_List(t *testing.T) {
	repo := newFakePageContentRepository()
	svc := service.NewPageContentService(repo)

	items, err := svc.List(context.Background())

	require.NoError(t, err)
	assert.Len(t, items, 1)
}

func TestPageContentService_Update(t *testing.T) {
	repo := newFakePageContentRepository()
	svc := service.NewPageContentService(repo)

	t.Run("updates an existing key", func(t *testing.T) {
		item, err := svc.Update(context.Background(), "home_hero_title", "Новое значение")

		require.NoError(t, err)
		assert.Equal(t, "Новое значение", item.Value)
	})

	t.Run("returns ErrNotFound for an unknown key", func(t *testing.T) {
		_, err := svc.Update(context.Background(), "does-not-exist", "value")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrNotFound)
	})

	t.Run("rejects an empty key", func(t *testing.T) {
		_, err := svc.Update(context.Background(), "   ", "value")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}
