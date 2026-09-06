package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type fakePageFAQRepository struct {
	settings map[string]model.PageFAQSettings
	items    []model.PageFAQItem
	nextID   int64
}

func newFakePageFAQRepository() *fakePageFAQRepository {
	return &fakePageFAQRepository{
		settings: map[string]model.PageFAQSettings{
			"masterclasses":    {Page: "masterclasses"},
			"gift_certificate": {Page: "gift_certificate"},
		},
		nextID: 1,
	}
}

func (f *fakePageFAQRepository) GetSettings(ctx context.Context, page string) (model.PageFAQSettings, error) {
	settings, ok := f.settings[page]
	if !ok {
		return model.PageFAQSettings{}, service.ErrNotFound
	}
	return settings, nil
}

// UpdateSettings only matches a page that already has a seeded row — mirrors
// the real repository's WHERE page = $x (never inserts).
func (f *fakePageFAQRepository) UpdateSettings(ctx context.Context, item model.PageFAQSettings) (model.PageFAQSettings, error) {
	if _, ok := f.settings[item.Page]; !ok {
		return model.PageFAQSettings{}, service.ErrNotFound
	}
	f.settings[item.Page] = item
	return item, nil
}

func (f *fakePageFAQRepository) ListByPage(ctx context.Context, page string) ([]model.PageFAQItem, error) {
	var items []model.PageFAQItem
	for _, item := range f.items {
		if item.Page == page {
			items = append(items, item)
		}
	}
	return items, nil
}

func (f *fakePageFAQRepository) CreateItem(ctx context.Context, item model.PageFAQItem) (model.PageFAQItem, error) {
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

// UpdateItem only matches a row that belongs to the given page — mirrors the
// real repository's WHERE id = $x AND page = $y.
func (f *fakePageFAQRepository) UpdateItem(ctx context.Context, item model.PageFAQItem) (model.PageFAQItem, error) {
	for i, existing := range f.items {
		if existing.ID == item.ID && existing.Page == item.Page {
			f.items[i] = item
			return item, nil
		}
	}
	return model.PageFAQItem{}, service.ErrNotFound
}

func (f *fakePageFAQRepository) DeleteItem(ctx context.Context, page string, id int64) error {
	for i, existing := range f.items {
		if existing.ID == id && existing.Page == page {
			f.items = append(f.items[:i], f.items[i+1:]...)
			return nil
		}
	}
	return service.ErrNotFound
}

func TestPageFAQService_Get(t *testing.T) {
	t.Run("returns settings and items for a valid page", func(t *testing.T) {
		repo := newFakePageFAQRepository()
		svc := service.NewPageFAQService(repo)
		_, err := svc.CreateItem(context.Background(), model.PageFAQItem{Page: "masterclasses", Question: "q", Answer: "a"})
		require.NoError(t, err)

		result, err := svc.Get(context.Background(), "masterclasses")

		require.NoError(t, err)
		assert.Equal(t, "masterclasses", result.Page)
		require.Len(t, result.Items, 1)
		assert.Equal(t, "q", result.Items[0].Question)
	})

	t.Run("rejects an unknown page", func(t *testing.T) {
		repo := newFakePageFAQRepository()
		svc := service.NewPageFAQService(repo)

		_, err := svc.Get(context.Background(), "home")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestPageFAQService_UpdateSettings(t *testing.T) {
	repo := newFakePageFAQRepository()
	svc := service.NewPageFAQService(repo)

	t.Run("updates a valid page", func(t *testing.T) {
		updated, err := svc.UpdateSettings(context.Background(), model.PageFAQSettings{
			Page:        "masterclasses",
			Title:       "  Вопросы и ответы  ",
			Description: "  Текст  ",
			Visible:     true,
		})

		require.NoError(t, err)
		assert.Equal(t, "Вопросы и ответы", updated.Title)
		assert.Equal(t, "Текст", updated.Description)
		assert.True(t, updated.Visible)
	})

	t.Run("rejects an unknown page", func(t *testing.T) {
		_, err := svc.UpdateSettings(context.Background(), model.PageFAQSettings{Page: "home"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestPageFAQService_CreateItem(t *testing.T) {
	t.Run("creates a valid item", func(t *testing.T) {
		repo := newFakePageFAQRepository()
		svc := service.NewPageFAQService(repo)

		item, err := svc.CreateItem(context.Background(), model.PageFAQItem{
			Page:     "gift_certificate",
			Question: "  Сколько действует сертификат?  ",
			Answer:   "  Год ",
		})

		require.NoError(t, err)
		assert.Equal(t, "Сколько действует сертификат?", item.Question)
		assert.Equal(t, "Год", item.Answer)
	})

	t.Run("rejects an unknown page", func(t *testing.T) {
		repo := newFakePageFAQRepository()
		svc := service.NewPageFAQService(repo)

		_, err := svc.CreateItem(context.Background(), model.PageFAQItem{Page: "home", Question: "q", Answer: "a"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects an empty question or answer", func(t *testing.T) {
		repo := newFakePageFAQRepository()
		svc := service.NewPageFAQService(repo)

		_, err := svc.CreateItem(context.Background(), model.PageFAQItem{Page: "masterclasses", Question: " ", Answer: "a"})
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestPageFAQService_UpdateItem(t *testing.T) {
	repo := newFakePageFAQRepository()
	svc := service.NewPageFAQService(repo)
	created, err := svc.CreateItem(context.Background(), model.PageFAQItem{Page: "masterclasses", Question: "q", Answer: "a"})
	require.NoError(t, err)

	t.Run("does not update an item belonging to a different page", func(t *testing.T) {
		other, err := svc.CreateItem(context.Background(), model.PageFAQItem{Page: "gift_certificate", Question: "other page item", Answer: "a"})
		require.NoError(t, err)

		_, err = svc.UpdateItem(context.Background(), model.PageFAQItem{ID: other.ID, Page: "masterclasses", Question: "hijacked", Answer: "a"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrNotFound)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.UpdateItem(context.Background(), model.PageFAQItem{Page: "masterclasses", Question: "q", Answer: "a"})
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("updates an existing item", func(t *testing.T) {
		created.Question = "updated"
		updated, err := svc.UpdateItem(context.Background(), created)
		require.NoError(t, err)
		assert.Equal(t, "updated", updated.Question)
	})
}

func TestPageFAQService_DeleteItem(t *testing.T) {
	repo := newFakePageFAQRepository()
	svc := service.NewPageFAQService(repo)
	created, err := svc.CreateItem(context.Background(), model.PageFAQItem{Page: "masterclasses", Question: "q", Answer: "a"})
	require.NoError(t, err)

	t.Run("does not delete an item belonging to a different page", func(t *testing.T) {
		err := svc.DeleteItem(context.Background(), "gift_certificate", created.ID)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrNotFound)
	})

	t.Run("deletes an existing item", func(t *testing.T) {
		require.NoError(t, svc.DeleteItem(context.Background(), "masterclasses", created.ID))
		items, err := svc.Get(context.Background(), "masterclasses")
		require.NoError(t, err)
		assert.Empty(t, items.Items)
	})
}
