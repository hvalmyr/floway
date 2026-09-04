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

type fakeTagSearchRepository struct {
	tags       []model.Tag
	deletedIDs []int64
	deleteErr  error
}

func (f *fakeTagSearchRepository) Search(ctx context.Context, query string) ([]model.Tag, error) {
	return f.tags, nil
}

func (f *fakeTagSearchRepository) Delete(ctx context.Context, id int64) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}
	f.deletedIDs = append(f.deletedIDs, id)
	return nil
}

func TestTagService_Search(t *testing.T) {
	t.Run("routes product to the product repository", func(t *testing.T) {
		productTags := &fakeTagSearchRepository{tags: []model.Tag{{ID: 1, Name: "Курс"}}}
		typeTags := &fakeTagSearchRepository{tags: []model.Tag{{ID: 2, Name: "Постоянный"}}}
		svc := service.NewTagService(productTags, typeTags)

		tags, err := svc.Search(context.Background(), model.TagTypeProduct, "")

		require.NoError(t, err)
		require.Len(t, tags, 1)
		assert.Equal(t, "Курс", tags[0].Name)
	})

	t.Run("routes client_type to the client-type repository", func(t *testing.T) {
		productTags := &fakeTagSearchRepository{tags: []model.Tag{{ID: 1, Name: "Курс"}}}
		typeTags := &fakeTagSearchRepository{tags: []model.Tag{{ID: 2, Name: "Постоянный"}}}
		svc := service.NewTagService(productTags, typeTags)

		tags, err := svc.Search(context.Background(), model.TagTypeClientType, "")

		require.NoError(t, err)
		require.Len(t, tags, 1)
		assert.Equal(t, "Постоянный", tags[0].Name)
	})

	t.Run("rejects an invalid tag type", func(t *testing.T) {
		svc := service.NewTagService(&fakeTagSearchRepository{}, &fakeTagSearchRepository{})

		_, err := svc.Search(context.Background(), "unknown", "")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestTagService_Delete(t *testing.T) {
	t.Run("deletes from the product repository for product", func(t *testing.T) {
		productTags := &fakeTagSearchRepository{}
		typeTags := &fakeTagSearchRepository{}
		svc := service.NewTagService(productTags, typeTags)

		err := svc.Delete(context.Background(), model.TagTypeProduct, 5)

		require.NoError(t, err)
		assert.Equal(t, []int64{5}, productTags.deletedIDs)
		assert.Empty(t, typeTags.deletedIDs)
	})

	t.Run("deletes from the client-type repository for client_type", func(t *testing.T) {
		productTags := &fakeTagSearchRepository{}
		typeTags := &fakeTagSearchRepository{}
		svc := service.NewTagService(productTags, typeTags)

		err := svc.Delete(context.Background(), model.TagTypeClientType, 7)

		require.NoError(t, err)
		assert.Equal(t, []int64{7}, typeTags.deletedIDs)
		assert.Empty(t, productTags.deletedIDs)
	})

	t.Run("rejects a zero id", func(t *testing.T) {
		svc := service.NewTagService(&fakeTagSearchRepository{}, &fakeTagSearchRepository{})

		err := svc.Delete(context.Background(), model.TagTypeProduct, 0)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects an invalid tag type", func(t *testing.T) {
		svc := service.NewTagService(&fakeTagSearchRepository{}, &fakeTagSearchRepository{})

		err := svc.Delete(context.Background(), "unknown", 1)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("propagates a not-found delete", func(t *testing.T) {
		productTags := &fakeTagSearchRepository{deleteErr: errors.New("not found")}
		svc := service.NewTagService(productTags, &fakeTagSearchRepository{})

		err := svc.Delete(context.Background(), model.TagTypeProduct, 1)

		require.Error(t, err)
	})
}
