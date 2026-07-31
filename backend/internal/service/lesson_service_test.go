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

type fakeLessonRepository struct {
	items  []model.Lesson
	nextID int64
	err    error
}

func newFakeLessonRepository() *fakeLessonRepository {
	return &fakeLessonRepository{nextID: 1}
}

func (f *fakeLessonRepository) ListByCourseBlockID(ctx context.Context, courseBlockID int64) ([]model.Lesson, error) {
	if f.err != nil {
		return nil, f.err
	}
	var items []model.Lesson
	for _, item := range f.items {
		if item.CourseBlockID == courseBlockID {
			items = append(items, item)
		}
	}
	return items, nil
}

func (f *fakeLessonRepository) Create(ctx context.Context, item model.Lesson) (model.Lesson, error) {
	if f.err != nil {
		return model.Lesson{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeLessonRepository) Update(ctx context.Context, item model.Lesson) (model.Lesson, error) {
	if f.err != nil {
		return model.Lesson{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.Lesson{}, errors.New("not found")
}

func (f *fakeLessonRepository) Delete(ctx context.Context, id int64) error {
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

func TestLessonService_Create(t *testing.T) {
	t.Run("creates a valid item", func(t *testing.T) {
		repo := newFakeLessonRepository()
		svc := service.NewLessonService(repo)

		item, err := svc.Create(context.Background(), model.Lesson{
			CourseBlockID: 1,
			Number:        1,
			Title:         "  Введение  ",
			Topics:        "topics",
			Outcomes:      "outcomes",
			DurationHours: 2,
		})

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
		assert.Equal(t, "Введение", item.Title)
		assert.Equal(t, int64(1), item.CourseBlockID)
	})

	t.Run("rejects an empty title", func(t *testing.T) {
		repo := newFakeLessonRepository()
		svc := service.NewLessonService(repo)

		_, err := svc.Create(context.Background(), model.Lesson{CourseBlockID: 1, Title: "   "})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})

	t.Run("rejects a zero courseBlockId", func(t *testing.T) {
		repo := newFakeLessonRepository()
		svc := service.NewLessonService(repo)

		_, err := svc.Create(context.Background(), model.Lesson{Title: "title"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})
}

func TestLessonService_Update(t *testing.T) {
	repo := newFakeLessonRepository()
	svc := service.NewLessonService(repo)
	created, err := svc.Create(context.Background(), model.Lesson{CourseBlockID: 1, Title: "title"})
	require.NoError(t, err)

	t.Run("updates an existing item", func(t *testing.T) {
		created.Title = "updated title"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "updated title", updated.Title)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.Lesson{Title: "title"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestLessonService_Delete(t *testing.T) {
	repo := newFakeLessonRepository()
	svc := service.NewLessonService(repo)
	created, err := svc.Create(context.Background(), model.Lesson{CourseBlockID: 1, Title: "title"})
	require.NoError(t, err)

	t.Run("deletes an existing item", func(t *testing.T) {
		require.NoError(t, svc.Delete(context.Background(), created.ID))
		items, err := svc.ListByCourseBlockID(context.Background(), created.CourseBlockID)
		require.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("rejects a zero id", func(t *testing.T) {
		err := svc.Delete(context.Background(), 0)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestLessonService_ListByCourseBlockID(t *testing.T) {
	repo := newFakeLessonRepository()
	svc := service.NewLessonService(repo)
	_, err := svc.Create(context.Background(), model.Lesson{CourseBlockID: 1, Title: "block 1 lesson"})
	require.NoError(t, err)
	_, err = svc.Create(context.Background(), model.Lesson{CourseBlockID: 2, Title: "block 2 lesson"})
	require.NoError(t, err)

	t.Run("filters by course block id", func(t *testing.T) {
		items, err := svc.ListByCourseBlockID(context.Background(), 1)

		require.NoError(t, err)
		require.Len(t, items, 1)
		assert.Equal(t, "block 1 lesson", items[0].Title)
	})

	t.Run("propagates repository error", func(t *testing.T) {
		repo := newFakeLessonRepository()
		repo.err = errors.New("boom")
		svc := service.NewLessonService(repo)

		_, err := svc.ListByCourseBlockID(context.Background(), 1)

		require.Error(t, err)
	})
}
