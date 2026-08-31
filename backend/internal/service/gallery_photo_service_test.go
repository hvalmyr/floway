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

type fakeGalleryPhotoRepository struct {
	items  []model.GalleryPhoto
	nextID int64
	err    error
}

func newFakeGalleryPhotoRepository() *fakeGalleryPhotoRepository {
	return &fakeGalleryPhotoRepository{nextID: 1}
}

func (f *fakeGalleryPhotoRepository) List(ctx context.Context) ([]model.GalleryPhoto, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func (f *fakeGalleryPhotoRepository) Create(ctx context.Context, item model.GalleryPhoto) (model.GalleryPhoto, error) {
	if f.err != nil {
		return model.GalleryPhoto{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeGalleryPhotoRepository) Update(ctx context.Context, item model.GalleryPhoto) (model.GalleryPhoto, error) {
	if f.err != nil {
		return model.GalleryPhoto{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.GalleryPhoto{}, service.ErrNotFound
}

func (f *fakeGalleryPhotoRepository) Delete(ctx context.Context, id int64) error {
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

func TestGalleryPhotoService_Create(t *testing.T) {
	t.Run("creates a valid photo", func(t *testing.T) {
		repo := newFakeGalleryPhotoRepository()
		svc := service.NewGalleryPhotoService(repo)

		item, err := svc.Create(context.Background(), model.GalleryPhoto{Image: "gallery/1.jpg", SortOrder: 1})

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
		assert.Equal(t, "gallery/1.jpg", item.Image)
	})

	t.Run("rejects an empty image", func(t *testing.T) {
		repo := newFakeGalleryPhotoRepository()
		svc := service.NewGalleryPhotoService(repo)

		_, err := svc.Create(context.Background(), model.GalleryPhoto{})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})
}

func TestGalleryPhotoService_Update(t *testing.T) {
	repo := newFakeGalleryPhotoRepository()
	svc := service.NewGalleryPhotoService(repo)
	created, err := svc.Create(context.Background(), model.GalleryPhoto{Image: "gallery/1.jpg"})
	require.NoError(t, err)

	t.Run("updates an existing photo", func(t *testing.T) {
		created.Image = "gallery/updated.jpg"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "gallery/updated.jpg", updated.Image)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.GalleryPhoto{Image: "gallery/1.jpg"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestGalleryPhotoService_Delete(t *testing.T) {
	repo := newFakeGalleryPhotoRepository()
	svc := service.NewGalleryPhotoService(repo)
	created, err := svc.Create(context.Background(), model.GalleryPhoto{Image: "gallery/1.jpg"})
	require.NoError(t, err)

	t.Run("deletes an existing photo", func(t *testing.T) {
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

func TestGalleryPhotoService_List_PropagatesRepositoryError(t *testing.T) {
	repo := newFakeGalleryPhotoRepository()
	repo.err = errors.New("boom")
	svc := service.NewGalleryPhotoService(repo)

	_, err := svc.List(context.Background())

	require.Error(t, err)
}
