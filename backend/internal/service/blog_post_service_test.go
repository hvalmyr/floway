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

type fakeBlogPostRepository struct {
	items  []model.BlogPost
	nextID int64
	err    error
}

func newFakeBlogPostRepository() *fakeBlogPostRepository {
	return &fakeBlogPostRepository{nextID: 1}
}

func (f *fakeBlogPostRepository) List(ctx context.Context) ([]model.BlogPost, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func (f *fakeBlogPostRepository) Create(ctx context.Context, item model.BlogPost) (model.BlogPost, error) {
	if f.err != nil {
		return model.BlogPost{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeBlogPostRepository) Update(ctx context.Context, item model.BlogPost) (model.BlogPost, error) {
	if f.err != nil {
		return model.BlogPost{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.BlogPost{}, errors.New("not found")
}

func (f *fakeBlogPostRepository) Delete(ctx context.Context, id int64) error {
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

func TestBlogPostService_Create(t *testing.T) {
	t.Run("creates a valid item and defaults status to draft", func(t *testing.T) {
		repo := newFakeBlogPostRepository()
		svc := service.NewBlogPostService(repo)

		item, err := svc.Create(context.Background(), model.BlogPost{Slug: "  hello-world  ", Title: "  Hello  "})

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
		assert.Equal(t, "hello-world", item.Slug)
		assert.Equal(t, "Hello", item.Title)
		assert.Equal(t, model.BlogPostStatusDraft, item.Status)
	})

	t.Run("rejects an empty slug", func(t *testing.T) {
		repo := newFakeBlogPostRepository()
		svc := service.NewBlogPostService(repo)

		_, err := svc.Create(context.Background(), model.BlogPost{Slug: "  ", Title: "Hello"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})

	t.Run("rejects an invalid status", func(t *testing.T) {
		repo := newFakeBlogPostRepository()
		svc := service.NewBlogPostService(repo)

		_, err := svc.Create(context.Background(), model.BlogPost{Slug: "a", Title: "b", Status: "archived"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestBlogPostService_Update(t *testing.T) {
	repo := newFakeBlogPostRepository()
	svc := service.NewBlogPostService(repo)
	created, err := svc.Create(context.Background(), model.BlogPost{Slug: "a", Title: "b"})
	require.NoError(t, err)

	t.Run("updates an existing item", func(t *testing.T) {
		created.Status = model.BlogPostStatusPublished
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, model.BlogPostStatusPublished, updated.Status)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.BlogPost{Slug: "a", Title: "b"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestBlogPostService_Delete(t *testing.T) {
	repo := newFakeBlogPostRepository()
	svc := service.NewBlogPostService(repo)
	created, err := svc.Create(context.Background(), model.BlogPost{Slug: "a", Title: "b"})
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

func TestBlogPostService_List_PropagatesRepositoryError(t *testing.T) {
	repo := newFakeBlogPostRepository()
	repo.err = errors.New("boom")
	svc := service.NewBlogPostService(repo)

	_, err := svc.List(context.Background())

	require.Error(t, err)
}
