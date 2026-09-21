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

type fakeObjectTypeRepository struct {
	items  []model.ObjectType
	nextID int64
	err    error
}

func newFakeObjectTypeRepository() *fakeObjectTypeRepository {
	return &fakeObjectTypeRepository{nextID: 1}
}

func (f *fakeObjectTypeRepository) List(ctx context.Context) ([]model.ObjectType, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func (f *fakeObjectTypeRepository) Create(ctx context.Context, item model.ObjectType) (model.ObjectType, error) {
	if f.err != nil {
		return model.ObjectType{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeObjectTypeRepository) Update(ctx context.Context, item model.ObjectType) (model.ObjectType, error) {
	if f.err != nil {
		return model.ObjectType{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.ObjectType{}, service.ErrNotFound
}

func (f *fakeObjectTypeRepository) Delete(ctx context.Context, id int64) error {
	if f.err != nil {
		return f.err
	}
	for i, existing := range f.items {
		if existing.ID == id {
			f.items = append(f.items[:i], f.items[i+1:]...)
			return nil
		}
	}
	return service.ErrNotFound
}

func TestObjectTypeService_Create(t *testing.T) {
	t.Run("creates a valid item", func(t *testing.T) {
		repo := newFakeObjectTypeRepository()
		svc := service.NewObjectTypeService(repo)

		item, err := svc.Create(context.Background(), model.ObjectType{Slug: "house", Name: "Дом"})

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
	})

	t.Run("rejects an empty slug", func(t *testing.T) {
		repo := newFakeObjectTypeRepository()
		svc := service.NewObjectTypeService(repo)

		_, err := svc.Create(context.Background(), model.ObjectType{Name: "Дом"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects an empty name", func(t *testing.T) {
		repo := newFakeObjectTypeRepository()
		svc := service.NewObjectTypeService(repo)

		_, err := svc.Create(context.Background(), model.ObjectType{Slug: "house", Name: "  "})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestObjectTypeService_Update(t *testing.T) {
	repo := newFakeObjectTypeRepository()
	svc := service.NewObjectTypeService(repo)
	created, err := svc.Create(context.Background(), model.ObjectType{Slug: "house", Name: "Дом"})
	require.NoError(t, err)

	t.Run("updates an existing item", func(t *testing.T) {
		created.Name = "Загородный дом"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "Загородный дом", updated.Name)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.ObjectType{Slug: "house", Name: "Дом"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestObjectTypeService_Delete(t *testing.T) {
	repo := newFakeObjectTypeRepository()
	svc := service.NewObjectTypeService(repo)
	created, err := svc.Create(context.Background(), model.ObjectType{Slug: "house", Name: "Дом"})
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

func TestObjectTypeService_List_PropagatesRepositoryError(t *testing.T) {
	repo := newFakeObjectTypeRepository()
	repo.err = errors.New("boom")
	svc := service.NewObjectTypeService(repo)

	_, err := svc.List(context.Background())

	require.Error(t, err)
}
