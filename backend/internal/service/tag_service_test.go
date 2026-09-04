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
	tags          []model.Tag
	deletedIDs    []int64
	deleteErr     error
	updatedColors map[int64]string
	updateErr     error
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

func (f *fakeTagSearchRepository) UpdateColor(ctx context.Context, id int64, color string) (model.Tag, error) {
	if f.updateErr != nil {
		return model.Tag{}, f.updateErr
	}
	if f.updatedColors == nil {
		f.updatedColors = map[int64]string{}
	}
	f.updatedColors[id] = color
	return model.Tag{ID: id, Color: color}, nil
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

func TestTagService_SetColor(t *testing.T) {
	t.Run("sets the color on the product repository for product", func(t *testing.T) {
		productTags := &fakeTagSearchRepository{}
		typeTags := &fakeTagSearchRepository{}
		svc := service.NewTagService(productTags, typeTags)

		tag, err := svc.SetColor(context.Background(), model.TagTypeProduct, 5, "#a1b2c3")

		require.NoError(t, err)
		assert.Equal(t, "#a1b2c3", tag.Color)
		assert.Equal(t, "#a1b2c3", productTags.updatedColors[5])
		assert.Empty(t, typeTags.updatedColors)
	})

	t.Run("sets the color on the client-type repository for client_type", func(t *testing.T) {
		productTags := &fakeTagSearchRepository{}
		typeTags := &fakeTagSearchRepository{}
		svc := service.NewTagService(productTags, typeTags)

		_, err := svc.SetColor(context.Background(), model.TagTypeClientType, 7, "#a1b2c3")

		require.NoError(t, err)
		assert.Equal(t, "#a1b2c3", typeTags.updatedColors[7])
		assert.Empty(t, productTags.updatedColors)
	})

	t.Run("rejects a zero id", func(t *testing.T) {
		svc := service.NewTagService(&fakeTagSearchRepository{}, &fakeTagSearchRepository{})

		_, err := svc.SetColor(context.Background(), model.TagTypeProduct, 0, "#a1b2c3")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects an invalid tag type", func(t *testing.T) {
		svc := service.NewTagService(&fakeTagSearchRepository{}, &fakeTagSearchRepository{})

		_, err := svc.SetColor(context.Background(), "unknown", 1, "#a1b2c3")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects a malformed color", func(t *testing.T) {
		svc := service.NewTagService(&fakeTagSearchRepository{}, &fakeTagSearchRepository{})

		for _, bad := range []string{"", "red", "#fff", "#gggggg", "a1b2c3", "#a1b2c344"} {
			_, err := svc.SetColor(context.Background(), model.TagTypeProduct, 1, bad)
			require.Errorf(t, err, "expected %q to be rejected", bad)
			assert.ErrorIs(t, err, service.ErrValidation)
		}
	})

	t.Run("propagates a not-found update", func(t *testing.T) {
		productTags := &fakeTagSearchRepository{updateErr: errors.New("not found")}
		svc := service.NewTagService(productTags, &fakeTagSearchRepository{})

		_, err := svc.SetColor(context.Background(), model.TagTypeProduct, 1, "#a1b2c3")

		require.Error(t, err)
	})
}
