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

type fakeCourseBlockRepository struct {
	items  []model.CourseBlock
	nextID int64
	err    error
}

func newFakeCourseBlockRepository() *fakeCourseBlockRepository {
	return &fakeCourseBlockRepository{nextID: 1}
}

func (f *fakeCourseBlockRepository) ListByCourseID(ctx context.Context, courseID int64) ([]model.CourseBlock, error) {
	if f.err != nil {
		return nil, f.err
	}
	var items []model.CourseBlock
	for _, existing := range f.items {
		if existing.CourseID == courseID {
			items = append(items, existing)
		}
	}
	return items, nil
}

func (f *fakeCourseBlockRepository) Create(ctx context.Context, item model.CourseBlock) (model.CourseBlock, error) {
	if f.err != nil {
		return model.CourseBlock{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeCourseBlockRepository) Update(ctx context.Context, item model.CourseBlock) (model.CourseBlock, error) {
	if f.err != nil {
		return model.CourseBlock{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.CourseBlock{}, errors.New("not found")
}

func (f *fakeCourseBlockRepository) Delete(ctx context.Context, id int64) error {
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

func TestCourseBlockService_Create(t *testing.T) {
	t.Run("creates a valid item", func(t *testing.T) {
		repo := newFakeCourseBlockRepository()
		svc := service.NewCourseBlockService(repo)

		item, err := svc.Create(context.Background(), model.CourseBlock{
			CourseID:     1,
			Title:        "  Основы  ",
			LessonsCount: 5,
			Hours:        10,
			Price:        1000,
			SortOrder:    1,
		})

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
		assert.Equal(t, "Основы", item.Title)
		assert.Equal(t, int64(1), item.CourseID)
	})

	t.Run("rejects an empty title", func(t *testing.T) {
		repo := newFakeCourseBlockRepository()
		svc := service.NewCourseBlockService(repo)

		_, err := svc.Create(context.Background(), model.CourseBlock{CourseID: 1, Title: "   "})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})

	t.Run("rejects a zero courseId", func(t *testing.T) {
		repo := newFakeCourseBlockRepository()
		svc := service.NewCourseBlockService(repo)

		_, err := svc.Create(context.Background(), model.CourseBlock{CourseID: 0, Title: "Основы"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})

	t.Run("accepts an oldPrice greater than price", func(t *testing.T) {
		repo := newFakeCourseBlockRepository()
		svc := service.NewCourseBlockService(repo)
		oldPrice := 1200

		item, err := svc.Create(context.Background(), model.CourseBlock{
			CourseID: 1, Title: "Основы", Price: 1000, OldPrice: &oldPrice,
		})

		require.NoError(t, err)
		require.NotNil(t, item.OldPrice)
		assert.Equal(t, 1200, *item.OldPrice)
	})

	t.Run("rejects an oldPrice not greater than price", func(t *testing.T) {
		repo := newFakeCourseBlockRepository()
		svc := service.NewCourseBlockService(repo)
		oldPrice := 900

		_, err := svc.Create(context.Background(), model.CourseBlock{
			CourseID: 1, Title: "Основы", Price: 1000, OldPrice: &oldPrice,
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestCourseBlockService_Update(t *testing.T) {
	repo := newFakeCourseBlockRepository()
	svc := service.NewCourseBlockService(repo)
	created, err := svc.Create(context.Background(), model.CourseBlock{CourseID: 1, Title: "Основы"})
	require.NoError(t, err)

	t.Run("updates an existing item", func(t *testing.T) {
		created.Title = "updated title"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "updated title", updated.Title)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.CourseBlock{Title: "title"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestCourseBlockService_Delete(t *testing.T) {
	repo := newFakeCourseBlockRepository()
	svc := service.NewCourseBlockService(repo)
	created, err := svc.Create(context.Background(), model.CourseBlock{CourseID: 1, Title: "Основы"})
	require.NoError(t, err)

	t.Run("deletes an existing item", func(t *testing.T) {
		require.NoError(t, svc.Delete(context.Background(), created.ID))
		items, err := svc.ListByCourseID(context.Background(), created.CourseID)
		require.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("rejects a zero id", func(t *testing.T) {
		err := svc.Delete(context.Background(), 0)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestCourseBlockService_ListByCourseID(t *testing.T) {
	t.Run("filters by courseId", func(t *testing.T) {
		repo := newFakeCourseBlockRepository()
		svc := service.NewCourseBlockService(repo)

		_, err := svc.Create(context.Background(), model.CourseBlock{CourseID: 1, Title: "Course 1 Block 1"})
		require.NoError(t, err)
		_, err = svc.Create(context.Background(), model.CourseBlock{CourseID: 2, Title: "Course 2 Block 1"})
		require.NoError(t, err)
		_, err = svc.Create(context.Background(), model.CourseBlock{CourseID: 1, Title: "Course 1 Block 2"})
		require.NoError(t, err)

		items, err := svc.ListByCourseID(context.Background(), 1)

		require.NoError(t, err)
		assert.Len(t, items, 2)
		for _, item := range items {
			assert.Equal(t, int64(1), item.CourseID)
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		repo := newFakeCourseBlockRepository()
		repo.err = errors.New("boom")
		svc := service.NewCourseBlockService(repo)

		_, err := svc.ListByCourseID(context.Background(), 1)

		require.Error(t, err)
	})
}
