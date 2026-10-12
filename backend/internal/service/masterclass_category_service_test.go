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

type fakeMasterclassCategoryRepository struct {
	items  []model.MasterclassCategory
	nextID int64
	err    error
}

func newFakeMasterclassCategoryRepository() *fakeMasterclassCategoryRepository {
	return &fakeMasterclassCategoryRepository{nextID: 1}
}

func (f *fakeMasterclassCategoryRepository) List(ctx context.Context) ([]model.MasterclassCategory, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func (f *fakeMasterclassCategoryRepository) Create(ctx context.Context, item model.MasterclassCategory) (model.MasterclassCategory, error) {
	if f.err != nil {
		return model.MasterclassCategory{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeMasterclassCategoryRepository) Update(ctx context.Context, item model.MasterclassCategory) (model.MasterclassCategory, error) {
	if f.err != nil {
		return model.MasterclassCategory{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.MasterclassCategory{}, service.ErrNotFound
}

func (f *fakeMasterclassCategoryRepository) Delete(ctx context.Context, id int64) error {
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

func TestMasterclassCategoryService_Create(t *testing.T) {
	t.Run("creates a valid item", func(t *testing.T) {
		repo := newFakeMasterclassCategoryRepository()
		svc := service.NewMasterclassCategoryService(repo)

		item, err := svc.Create(context.Background(), model.MasterclassCategory{Slug: "novogodnie", Name: "Новогодние"})

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
	})

	t.Run("rejects an empty slug", func(t *testing.T) {
		repo := newFakeMasterclassCategoryRepository()
		svc := service.NewMasterclassCategoryService(repo)

		_, err := svc.Create(context.Background(), model.MasterclassCategory{Name: "Новогодние"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects an empty name", func(t *testing.T) {
		repo := newFakeMasterclassCategoryRepository()
		svc := service.NewMasterclassCategoryService(repo)

		_, err := svc.Create(context.Background(), model.MasterclassCategory{Slug: "novogodnie", Name: "  "})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestMasterclassCategoryService_Update(t *testing.T) {
	repo := newFakeMasterclassCategoryRepository()
	svc := service.NewMasterclassCategoryService(repo)
	created, err := svc.Create(context.Background(), model.MasterclassCategory{Slug: "novogodnie", Name: "Новогодние"})
	require.NoError(t, err)

	t.Run("updates an existing item", func(t *testing.T) {
		created.Name = "Зимние"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "Зимние", updated.Name)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.MasterclassCategory{Slug: "novogodnie", Name: "Новогодние"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestMasterclassCategoryService_Delete(t *testing.T) {
	repo := newFakeMasterclassCategoryRepository()
	svc := service.NewMasterclassCategoryService(repo)
	created, err := svc.Create(context.Background(), model.MasterclassCategory{Slug: "novogodnie", Name: "Новогодние"})
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

func TestMasterclassCategoryService_List_PropagatesRepositoryError(t *testing.T) {
	repo := newFakeMasterclassCategoryRepository()
	repo.err = errors.New("boom")
	svc := service.NewMasterclassCategoryService(repo)

	_, err := svc.List(context.Background())

	require.Error(t, err)
}
