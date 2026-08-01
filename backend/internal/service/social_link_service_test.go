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

type fakeSocialLinkRepository struct {
	items  []model.SocialLink
	nextID int64
	err    error
}

func newFakeSocialLinkRepository() *fakeSocialLinkRepository {
	return &fakeSocialLinkRepository{nextID: 1}
}

func (f *fakeSocialLinkRepository) List(ctx context.Context) ([]model.SocialLink, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func (f *fakeSocialLinkRepository) Create(ctx context.Context, item model.SocialLink) (model.SocialLink, error) {
	if f.err != nil {
		return model.SocialLink{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeSocialLinkRepository) Update(ctx context.Context, item model.SocialLink) (model.SocialLink, error) {
	if f.err != nil {
		return model.SocialLink{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.SocialLink{}, errors.New("not found")
}

func (f *fakeSocialLinkRepository) Delete(ctx context.Context, id int64) error {
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

func TestSocialLinkService_Create(t *testing.T) {
	t.Run("creates a valid item", func(t *testing.T) {
		repo := newFakeSocialLinkRepository()
		svc := service.NewSocialLinkService(repo)

		item, err := svc.Create(context.Background(), model.SocialLink{Label: "  Telegram  ", Href: "  https://t.me/floway  "})

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
		assert.Equal(t, "Telegram", item.Label)
		assert.Equal(t, "https://t.me/floway", item.Href)
	})

	t.Run("rejects an empty href", func(t *testing.T) {
		repo := newFakeSocialLinkRepository()
		svc := service.NewSocialLinkService(repo)

		_, err := svc.Create(context.Background(), model.SocialLink{Label: "Telegram", Href: "  "})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})
}

func TestSocialLinkService_Update(t *testing.T) {
	repo := newFakeSocialLinkRepository()
	svc := service.NewSocialLinkService(repo)
	created, err := svc.Create(context.Background(), model.SocialLink{Label: "Telegram", Href: "https://t.me/floway"})
	require.NoError(t, err)

	t.Run("updates an existing item", func(t *testing.T) {
		created.Href = "https://t.me/floway2"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "https://t.me/floway2", updated.Href)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.SocialLink{Label: "l", Href: "h"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestSocialLinkService_Delete(t *testing.T) {
	repo := newFakeSocialLinkRepository()
	svc := service.NewSocialLinkService(repo)
	created, err := svc.Create(context.Background(), model.SocialLink{Label: "l", Href: "h"})
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
