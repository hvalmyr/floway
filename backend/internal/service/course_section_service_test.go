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

type fakeCourseSectionRepository struct {
	items  []model.CourseSection
	nextID int64
	err    error
}

func newFakeCourseSectionRepository() *fakeCourseSectionRepository {
	return &fakeCourseSectionRepository{nextID: 1}
}

func (f *fakeCourseSectionRepository) List(ctx context.Context) ([]model.CourseSection, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func (f *fakeCourseSectionRepository) Create(ctx context.Context, item model.CourseSection) (model.CourseSection, error) {
	if f.err != nil {
		return model.CourseSection{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeCourseSectionRepository) Update(ctx context.Context, item model.CourseSection) (model.CourseSection, error) {
	if f.err != nil {
		return model.CourseSection{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.CourseSection{}, service.ErrNotFound
}

func (f *fakeCourseSectionRepository) Delete(ctx context.Context, id int64) error {
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

func TestCourseSectionService_Create(t *testing.T) {
	t.Run("creates a valid section", func(t *testing.T) {
		repo := newFakeCourseSectionRepository()
		svc := service.NewCourseSectionService(repo)

		item, err := svc.Create(context.Background(), model.CourseSection{Heading: "  Курсы  "})

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
		assert.Equal(t, "Курсы", item.Heading)
	})

	t.Run("rejects an empty heading", func(t *testing.T) {
		repo := newFakeCourseSectionRepository()
		svc := service.NewCourseSectionService(repo)

		_, err := svc.Create(context.Background(), model.CourseSection{Heading: "   "})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})
}

func TestCourseSectionService_Update(t *testing.T) {
	repo := newFakeCourseSectionRepository()
	svc := service.NewCourseSectionService(repo)
	created, err := svc.Create(context.Background(), model.CourseSection{Heading: "Курсы"})
	require.NoError(t, err)

	t.Run("updates an existing section", func(t *testing.T) {
		created.Heading = "Обновлённый заголовок"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "Обновлённый заголовок", updated.Heading)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.CourseSection{Heading: "Курсы"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestCourseSectionService_Delete(t *testing.T) {
	repo := newFakeCourseSectionRepository()
	svc := service.NewCourseSectionService(repo)
	created, err := svc.Create(context.Background(), model.CourseSection{Heading: "Курсы"})
	require.NoError(t, err)

	t.Run("deletes an existing section", func(t *testing.T) {
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

func TestCourseSectionService_List_PropagatesRepositoryError(t *testing.T) {
	repo := newFakeCourseSectionRepository()
	repo.err = errors.New("boom")
	svc := service.NewCourseSectionService(repo)

	_, err := svc.List(context.Background())

	require.Error(t, err)
}
