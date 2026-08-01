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

type fakeFeatureRepository struct {
	items  []model.Feature
	nextID int64
	err    error
}

func newFakeFeatureRepository() *fakeFeatureRepository {
	return &fakeFeatureRepository{nextID: 1}
}

func (f *fakeFeatureRepository) List(ctx context.Context) ([]model.Feature, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func (f *fakeFeatureRepository) ListByPage(ctx context.Context, page string) ([]model.Feature, error) {
	if f.err != nil {
		return nil, f.err
	}
	var items []model.Feature
	for _, item := range f.items {
		if item.Page == page {
			items = append(items, item)
		}
	}
	return items, nil
}

func (f *fakeFeatureRepository) Create(ctx context.Context, item model.Feature) (model.Feature, error) {
	if f.err != nil {
		return model.Feature{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeFeatureRepository) Update(ctx context.Context, item model.Feature) (model.Feature, error) {
	if f.err != nil {
		return model.Feature{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.Feature{}, errors.New("not found")
}

func (f *fakeFeatureRepository) Delete(ctx context.Context, id int64) error {
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

func TestFeatureService_Create(t *testing.T) {
	t.Run("creates a valid item", func(t *testing.T) {
		repo := newFakeFeatureRepository()
		svc := service.NewFeatureService(repo)

		item, err := svc.Create(context.Background(), model.Feature{
			Page: "home", Icon: "levels", Title: "заголовок", Description: "описание",
		})

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
	})

	t.Run("rejects an invalid page", func(t *testing.T) {
		repo := newFakeFeatureRepository()
		svc := service.NewFeatureService(repo)

		_, err := svc.Create(context.Background(), model.Feature{
			Page: "blog", Icon: "levels", Title: "t", Description: "d",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects an empty title", func(t *testing.T) {
		repo := newFakeFeatureRepository()
		svc := service.NewFeatureService(repo)

		_, err := svc.Create(context.Background(), model.Feature{
			Page: "home", Icon: "levels", Title: "  ", Description: "d",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestFeatureService_ListByPage(t *testing.T) {
	repo := newFakeFeatureRepository()
	svc := service.NewFeatureService(repo)
	_, err := svc.Create(context.Background(), model.Feature{Page: "home", Icon: "a", Title: "t1", Description: "d"})
	require.NoError(t, err)
	_, err = svc.Create(context.Background(), model.Feature{Page: "masterclasses", Icon: "b", Title: "t2", Description: "d"})
	require.NoError(t, err)

	items, err := svc.ListByPage(context.Background(), "home")

	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "t1", items[0].Title)
}

func TestFeatureService_Update(t *testing.T) {
	repo := newFakeFeatureRepository()
	svc := service.NewFeatureService(repo)
	created, err := svc.Create(context.Background(), model.Feature{Page: "home", Icon: "a", Title: "t", Description: "d"})
	require.NoError(t, err)

	t.Run("updates an existing item", func(t *testing.T) {
		created.Title = "новый заголовок"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "новый заголовок", updated.Title)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.Feature{Page: "home", Icon: "a", Title: "t", Description: "d"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestFeatureService_Delete(t *testing.T) {
	repo := newFakeFeatureRepository()
	svc := service.NewFeatureService(repo)
	created, err := svc.Create(context.Background(), model.Feature{Page: "home", Icon: "a", Title: "t", Description: "d"})
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
