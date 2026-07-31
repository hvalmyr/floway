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

type fakeFAQRepository struct {
	items  []model.FAQItem
	nextID int64
	err    error
}

func newFakeFAQRepository() *fakeFAQRepository {
	return &fakeFAQRepository{nextID: 1}
}

func (f *fakeFAQRepository) List(ctx context.Context) ([]model.FAQItem, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func (f *fakeFAQRepository) Create(ctx context.Context, item model.FAQItem) (model.FAQItem, error) {
	if f.err != nil {
		return model.FAQItem{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeFAQRepository) Update(ctx context.Context, item model.FAQItem) (model.FAQItem, error) {
	if f.err != nil {
		return model.FAQItem{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.FAQItem{}, errors.New("not found")
}

func (f *fakeFAQRepository) Delete(ctx context.Context, id int64) error {
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

func TestFAQService_Create(t *testing.T) {
	t.Run("creates a valid item", func(t *testing.T) {
		repo := newFakeFAQRepository()
		svc := service.NewFAQService(repo)

		item, err := svc.Create(context.Background(), "  Сколько длится курс?  ", "  3 месяца ", 1)

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
		assert.Equal(t, "Сколько длится курс?", item.Question)
		assert.Equal(t, "3 месяца", item.Answer)
	})

	t.Run("rejects an empty question", func(t *testing.T) {
		repo := newFakeFAQRepository()
		svc := service.NewFAQService(repo)

		_, err := svc.Create(context.Background(), "   ", "answer", 0)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})

	t.Run("rejects an empty answer", func(t *testing.T) {
		repo := newFakeFAQRepository()
		svc := service.NewFAQService(repo)

		_, err := svc.Create(context.Background(), "question", "  ", 0)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestFAQService_Update(t *testing.T) {
	repo := newFakeFAQRepository()
	svc := service.NewFAQService(repo)
	created, err := svc.Create(context.Background(), "question", "answer", 0)
	require.NoError(t, err)

	t.Run("updates an existing item", func(t *testing.T) {
		created.Question = "updated question"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "updated question", updated.Question)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.FAQItem{Question: "q", Answer: "a"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestFAQService_Delete(t *testing.T) {
	repo := newFakeFAQRepository()
	svc := service.NewFAQService(repo)
	created, err := svc.Create(context.Background(), "question", "answer", 0)
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

func TestFAQService_List_PropagatesRepositoryError(t *testing.T) {
	repo := newFakeFAQRepository()
	repo.err = errors.New("boom")
	svc := service.NewFAQService(repo)

	_, err := svc.List(context.Background())

	require.Error(t, err)
}
