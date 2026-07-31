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

type fakeMasterclassRepository struct {
	items  []model.Masterclass
	nextID int64
	err    error
}

func newFakeMasterclassRepository() *fakeMasterclassRepository {
	return &fakeMasterclassRepository{nextID: 1}
}

func (f *fakeMasterclassRepository) List(ctx context.Context) ([]model.Masterclass, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func (f *fakeMasterclassRepository) Create(ctx context.Context, item model.Masterclass) (model.Masterclass, error) {
	if f.err != nil {
		return model.Masterclass{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeMasterclassRepository) Update(ctx context.Context, item model.Masterclass) (model.Masterclass, error) {
	if f.err != nil {
		return model.Masterclass{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.Masterclass{}, errors.New("not found")
}

func (f *fakeMasterclassRepository) Delete(ctx context.Context, id int64) error {
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

func TestMasterclassService_Create(t *testing.T) {
	t.Run("creates a valid item", func(t *testing.T) {
		repo := newFakeMasterclassRepository()
		svc := service.NewMasterclassService(repo)

		item, err := svc.Create(context.Background(), model.Masterclass{
			Slug:   "  vocal-basics  ",
			Title:  "  Основы вокала  ",
			Status: model.MasterclassStatusActive,
		})

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
		assert.Equal(t, "vocal-basics", item.Slug)
		assert.Equal(t, "Основы вокала", item.Title)
		assert.Equal(t, model.MasterclassStatusActive, item.Status)
	})

	t.Run("defaults status to active when empty", func(t *testing.T) {
		repo := newFakeMasterclassRepository()
		svc := service.NewMasterclassService(repo)

		item, err := svc.Create(context.Background(), model.Masterclass{
			Slug:  "vocal-basics",
			Title: "Основы вокала",
		})

		require.NoError(t, err)
		assert.Equal(t, model.MasterclassStatusActive, item.Status)
	})

	t.Run("rejects an empty slug", func(t *testing.T) {
		repo := newFakeMasterclassRepository()
		svc := service.NewMasterclassService(repo)

		_, err := svc.Create(context.Background(), model.Masterclass{
			Slug:  "   ",
			Title: "Основы вокала",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})

	t.Run("rejects an empty title", func(t *testing.T) {
		repo := newFakeMasterclassRepository()
		svc := service.NewMasterclassService(repo)

		_, err := svc.Create(context.Background(), model.Masterclass{
			Slug:  "vocal-basics",
			Title: "   ",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects an invalid status", func(t *testing.T) {
		repo := newFakeMasterclassRepository()
		svc := service.NewMasterclassService(repo)

		_, err := svc.Create(context.Background(), model.Masterclass{
			Slug:   "vocal-basics",
			Title:  "Основы вокала",
			Status: model.MasterclassStatus("draft"),
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})
}

func TestMasterclassService_Update(t *testing.T) {
	repo := newFakeMasterclassRepository()
	svc := service.NewMasterclassService(repo)
	created, err := svc.Create(context.Background(), model.Masterclass{
		Slug:  "vocal-basics",
		Title: "Основы вокала",
	})
	require.NoError(t, err)

	t.Run("updates an existing item", func(t *testing.T) {
		created.Title = "Основы вокала 2.0"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "Основы вокала 2.0", updated.Title)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.Masterclass{Slug: "slug", Title: "title"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestMasterclassService_Delete(t *testing.T) {
	repo := newFakeMasterclassRepository()
	svc := service.NewMasterclassService(repo)
	created, err := svc.Create(context.Background(), model.Masterclass{
		Slug:  "vocal-basics",
		Title: "Основы вокала",
	})
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

func TestMasterclassService_List_PropagatesRepositoryError(t *testing.T) {
	repo := newFakeMasterclassRepository()
	repo.err = errors.New("boom")
	svc := service.NewMasterclassService(repo)

	_, err := svc.List(context.Background())

	require.Error(t, err)
}
