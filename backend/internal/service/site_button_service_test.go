package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type fakeSiteButtonRepository struct {
	items map[string]model.SiteButton
}

func (f *fakeSiteButtonRepository) List(ctx context.Context) ([]model.SiteButton, error) {
	items := []model.SiteButton{}
	for _, item := range f.items {
		items = append(items, item)
	}
	return items, nil
}

func (f *fakeSiteButtonRepository) Update(ctx context.Context, key, text, variant, url string) (model.SiteButton, error) {
	existing, ok := f.items[key]
	if !ok {
		return model.SiteButton{}, assert.AnError
	}
	existing.Text, existing.Variant, existing.URL = text, variant, url
	f.items[key] = existing
	return existing, nil
}

func TestSiteButtonService_Update_RejectsUnknownVariant(t *testing.T) {
	repo := &fakeSiteButtonRepository{items: map[string]model.SiteButton{"home_hero_courses": {Key: "home_hero_courses"}}}
	svc := service.NewSiteButtonService(repo)

	_, err := svc.Update(context.Background(), "home_hero_courses", "Курсы", "danger", "/#courses")

	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrValidation)
}

func TestSiteButtonService_Update_RejectsBlankTextOrURL(t *testing.T) {
	repo := &fakeSiteButtonRepository{items: map[string]model.SiteButton{"home_hero_courses": {Key: "home_hero_courses"}}}
	svc := service.NewSiteButtonService(repo)

	_, err := svc.Update(context.Background(), "home_hero_courses", "  ", "primary", "/#courses")
	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrValidation)

	_, err = svc.Update(context.Background(), "home_hero_courses", "Курсы", "primary", "  ")
	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrValidation)
}

func TestSiteButtonService_Update_PersistsValidChange(t *testing.T) {
	repo := &fakeSiteButtonRepository{items: map[string]model.SiteButton{"home_hero_courses": {Key: "home_hero_courses"}}}
	svc := service.NewSiteButtonService(repo)

	updated, err := svc.Update(context.Background(), "home_hero_courses", "Наши курсы", "outline", "/#courses")

	require.NoError(t, err)
	assert.Equal(t, "Наши курсы", updated.Text)
	assert.Equal(t, "outline", updated.Variant)
	assert.Equal(t, "/#courses", updated.URL)
}
