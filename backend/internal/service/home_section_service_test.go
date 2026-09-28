package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type fakeHomeSectionRepository struct {
	items []model.HomeSection
}

func (f *fakeHomeSectionRepository) List(ctx context.Context) ([]model.HomeSection, error) {
	return f.items, nil
}

func (f *fakeHomeSectionRepository) Update(ctx context.Context, item model.HomeSection) (model.HomeSection, error) {
	for i, existing := range f.items {
		if existing.ID == item.ID {
			item.Key = existing.Key
			f.items[i] = item
			return item, nil
		}
	}
	return model.HomeSection{}, assert.AnError
}

func TestHomeSectionService_Update_RequiresID(t *testing.T) {
	svc := service.NewHomeSectionService(&fakeHomeSectionRepository{})

	_, err := svc.Update(context.Background(), model.HomeSection{Visible: true})

	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrValidation)
}

func TestHomeSectionService_Update_PersistsVisibilityAndOrder(t *testing.T) {
	repo := &fakeHomeSectionRepository{items: []model.HomeSection{{ID: 1, Key: "faq", Visible: true, SortOrder: 0}}}
	svc := service.NewHomeSectionService(repo)

	updated, err := svc.Update(context.Background(), model.HomeSection{ID: 1, Visible: false, SortOrder: 3})

	require.NoError(t, err)
	assert.False(t, updated.Visible)
	assert.Equal(t, 3, updated.SortOrder)
	assert.Equal(t, "faq", updated.Key)
}
