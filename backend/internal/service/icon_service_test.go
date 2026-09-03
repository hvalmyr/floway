package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type fakeIconRepository struct {
	items  []model.Icon
	nextID int64
}

func newFakeIconRepository() *fakeIconRepository {
	return &fakeIconRepository{nextID: 1}
}

func (f *fakeIconRepository) List(ctx context.Context) ([]model.Icon, error) {
	return f.items, nil
}

func (f *fakeIconRepository) Create(ctx context.Context, item model.Icon) (model.Icon, error) {
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeIconRepository) Delete(ctx context.Context, id int64) error {
	for i, item := range f.items {
		if item.ID == id {
			f.items = append(f.items[:i], f.items[i+1:]...)
			return nil
		}
	}
	return nil
}

func TestIconService_Create_RequiresNameAndSVG(t *testing.T) {
	svc := service.NewIconService(newFakeIconRepository())

	_, err := svc.Create(context.Background(), model.Icon{Name: "", SVG: `<svg></svg>`})
	assert.ErrorIs(t, err, service.ErrValidation)

	_, err = svc.Create(context.Background(), model.Icon{Name: "Gift", SVG: ""})
	assert.ErrorIs(t, err, service.ErrValidation)
}

func TestIconService_Create_RejectsInvalidSVG(t *testing.T) {
	svc := service.NewIconService(newFakeIconRepository())

	_, err := svc.Create(context.Background(), model.Icon{Name: "Broken", SVG: "not xml at all <<<"})
	assert.ErrorIs(t, err, service.ErrValidation)

	_, err = svc.Create(context.Background(), model.Icon{Name: "NoRoot", SVG: `<g><path d="M0 0"/></g>`})
	assert.ErrorIs(t, err, service.ErrValidation)
}

func TestIconService_Create_StripsScriptsAndEventHandlers(t *testing.T) {
	svc := service.NewIconService(newFakeIconRepository())

	dirty := `<svg viewBox="0 0 24 24" onload="alert(1)">` +
		`<script>alert(2)</script>` +
		`<path d="M0 0h24v24H0z" fill="currentColor" onclick="alert(3)"/>` +
		`<circle cx="12" cy="12" r="10"/>` +
		`<use href="javascript:alert(4)"/>` +
		`</svg>`

	item, err := svc.Create(context.Background(), model.Icon{Name: "Evil", SVG: dirty})
	require.NoError(t, err)

	assert.NotContains(t, item.SVG, "script")
	assert.NotContains(t, item.SVG, "onload")
	assert.NotContains(t, item.SVG, "onclick")
	assert.NotContains(t, item.SVG, "javascript:")
	assert.NotContains(t, item.SVG, "alert(")
	// The legitimate shapes survive.
	assert.Contains(t, item.SVG, `<path`)
	assert.Contains(t, item.SVG, `fill="currentColor"`)
	assert.Contains(t, item.SVG, `<circle`)
	assert.Contains(t, item.SVG, `<use`)
}

func TestIconService_Create_KeepsLocalFragmentHref(t *testing.T) {
	svc := service.NewIconService(newFakeIconRepository())

	clean := `<svg viewBox="0 0 24 24">` +
		`<defs><path id="shape" d="M0 0h24v24H0z"/></defs>` +
		`<use href="#shape"/>` +
		`</svg>`

	item, err := svc.Create(context.Background(), model.Icon{Name: "Sprite", SVG: clean})
	require.NoError(t, err)
	assert.Contains(t, item.SVG, `href="#shape"`)
}

func TestIconService_Delete_RequiresID(t *testing.T) {
	svc := service.NewIconService(newFakeIconRepository())
	err := svc.Delete(context.Background(), 0)
	assert.ErrorIs(t, err, service.ErrValidation)
}
