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

type fakeLandingPageRepository struct {
	items  []model.LandingPage
	nextID int64
	err    error
}

func newFakeLandingPageRepository() *fakeLandingPageRepository {
	return &fakeLandingPageRepository{nextID: 1}
}

func (f *fakeLandingPageRepository) List(ctx context.Context) ([]model.LandingPage, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func (f *fakeLandingPageRepository) FindBySlug(ctx context.Context, slug string) (model.LandingPage, error) {
	if f.err != nil {
		return model.LandingPage{}, f.err
	}
	for _, item := range f.items {
		if item.Slug == slug {
			return item, nil
		}
	}
	return model.LandingPage{}, service.ErrNotFound
}

func (f *fakeLandingPageRepository) Create(ctx context.Context, item model.LandingPage) (model.LandingPage, error) {
	if f.err != nil {
		return model.LandingPage{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeLandingPageRepository) Update(ctx context.Context, item model.LandingPage) (model.LandingPage, error) {
	if f.err != nil {
		return model.LandingPage{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.LandingPage{}, service.ErrNotFound
}

func (f *fakeLandingPageRepository) Delete(ctx context.Context, id int64) error {
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

func validLandingPage() model.LandingPage {
	return model.LandingPage{ObjectTypeID: 1, Slug: "dom", H1: "Оформление дома"}
}

func TestLandingPageService_Create(t *testing.T) {
	t.Run("creates a valid page", func(t *testing.T) {
		repo := newFakeLandingPageRepository()
		svc := service.NewLandingPageService(repo)

		item, err := svc.Create(context.Background(), validLandingPage())

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
	})

	t.Run("rejects a missing object type", func(t *testing.T) {
		repo := newFakeLandingPageRepository()
		svc := service.NewLandingPageService(repo)

		item := validLandingPage()
		item.ObjectTypeID = 0

		_, err := svc.Create(context.Background(), item)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects an empty h1", func(t *testing.T) {
		repo := newFakeLandingPageRepository()
		svc := service.NewLandingPageService(repo)

		item := validLandingPage()
		item.H1 = "  "

		_, err := svc.Create(context.Background(), item)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestLandingPageService_Update(t *testing.T) {
	repo := newFakeLandingPageRepository()
	svc := service.NewLandingPageService(repo)
	created, err := svc.Create(context.Background(), validLandingPage())
	require.NoError(t, err)

	t.Run("updates an existing page", func(t *testing.T) {
		created.H1 = "Новогоднее оформление дома"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "Новогоднее оформление дома", updated.H1)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		item := validLandingPage()
		_, err := svc.Update(context.Background(), item)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestLandingPageService_Delete(t *testing.T) {
	repo := newFakeLandingPageRepository()
	svc := service.NewLandingPageService(repo)
	created, err := svc.Create(context.Background(), validLandingPage())
	require.NoError(t, err)

	t.Run("deletes an existing page", func(t *testing.T) {
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

func TestLandingPageService_List_PropagatesRepositoryError(t *testing.T) {
	repo := newFakeLandingPageRepository()
	repo.err = errors.New("boom")
	svc := service.NewLandingPageService(repo)

	_, err := svc.List(context.Background())

	require.Error(t, err)
}
