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

// Update and Delete only match a row that belongs to the given course —
// mirrors the real repository's WHERE id = $x AND course_id = $y (see
// internal/repository/course_block_repository.go), so a URL like
// /courses/7/blocks/42 can't touch a block that actually belongs to a
// different course (architecture review finding #3).
func (f *fakeCourseBlockRepository) Update(ctx context.Context, item model.CourseBlock) (model.CourseBlock, error) {
	if f.err != nil {
		return model.CourseBlock{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID && existing.CourseID == item.CourseID {
			item.CourseID = existing.CourseID
			f.items[i] = item
			return item, nil
		}
	}
	return model.CourseBlock{}, service.ErrNotFound
}

func (f *fakeCourseBlockRepository) Delete(ctx context.Context, courseID, id int64) error {
	if f.err != nil {
		return f.err
	}
	for i, existing := range f.items {
		if existing.ID == id && existing.CourseID == courseID {
			f.items = append(f.items[:i], f.items[i+1:]...)
			return nil
		}
	}
	return service.ErrNotFound
}

func TestCourseBlockService_Create(t *testing.T) {
	t.Run("creates a valid item", func(t *testing.T) {
		repo := newFakeCourseBlockRepository()
		svc := service.NewCourseBlockService(repo)

		item, err := svc.Create(context.Background(), model.CourseBlock{
			CourseID:    1,
			BlockName:   "  Букеты  ",
			LessonCount: "7 занятий",
			TimeLength:  "30 часов",
			Price:       "38 500 ₽",
			SortOrder:   1,
		})

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
		assert.Equal(t, "Букеты", item.BlockName)
		assert.Equal(t, int64(1), item.CourseID)
	})

	t.Run("allows an empty blockName (undivided single-block course)", func(t *testing.T) {
		repo := newFakeCourseBlockRepository()
		svc := service.NewCourseBlockService(repo)

		item, err := svc.Create(context.Background(), model.CourseBlock{CourseID: 1, LessonCount: "7 занятий"})

		require.NoError(t, err)
		assert.Empty(t, item.BlockName)
	})

	t.Run("rejects a zero courseId", func(t *testing.T) {
		repo := newFakeCourseBlockRepository()
		svc := service.NewCourseBlockService(repo)

		_, err := svc.Create(context.Background(), model.CourseBlock{CourseID: 0, BlockName: "Букеты"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})

	t.Run("defaults displayStyle to blue-beige when empty", func(t *testing.T) {
		repo := newFakeCourseBlockRepository()
		svc := service.NewCourseBlockService(repo)

		item, err := svc.Create(context.Background(), model.CourseBlock{CourseID: 1, BlockName: "Букеты"})

		require.NoError(t, err)
		assert.Equal(t, model.DisplayStyleBlueBeige, item.DisplayStyle)
	})

	t.Run("accepts any of the 4 valid display styles", func(t *testing.T) {
		repo := newFakeCourseBlockRepository()
		svc := service.NewCourseBlockService(repo)

		item, err := svc.Create(context.Background(), model.CourseBlock{
			CourseID:     1,
			BlockName:    "Букеты",
			DisplayStyle: model.DisplayStyleBeigeBrown,
		})

		require.NoError(t, err)
		assert.Equal(t, model.DisplayStyleBeigeBrown, item.DisplayStyle)
	})

	t.Run("rejects an invalid displayStyle", func(t *testing.T) {
		repo := newFakeCourseBlockRepository()
		svc := service.NewCourseBlockService(repo)

		_, err := svc.Create(context.Background(), model.CourseBlock{
			CourseID:     1,
			BlockName:    "Букеты",
			DisplayStyle: "green-purple",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})
}

func TestCourseBlockService_Update(t *testing.T) {
	repo := newFakeCourseBlockRepository()
	svc := service.NewCourseBlockService(repo)
	created, err := svc.Create(context.Background(), model.CourseBlock{CourseID: 1, BlockName: "Основы"})
	require.NoError(t, err)

	t.Run("updates an existing item", func(t *testing.T) {
		created.BlockName = "updated name"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "updated name", updated.BlockName)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.CourseBlock{CourseID: 1, BlockName: "name"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects a missing courseId", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.CourseBlock{ID: created.ID, BlockName: "name"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("does not update a block belonging to a different course", func(t *testing.T) {
		other, err := svc.Create(context.Background(), model.CourseBlock{CourseID: 2, BlockName: "Course 2 block"})
		require.NoError(t, err)

		// Same block id as `created`'s course-2 sibling, but wrong courseId —
		// simulates PUT /courses/1/blocks/{other.ID}.
		_, err = svc.Update(context.Background(), model.CourseBlock{ID: other.ID, CourseID: 1, BlockName: "hijacked"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrNotFound)
	})
}

func TestCourseBlockService_Delete(t *testing.T) {
	repo := newFakeCourseBlockRepository()
	svc := service.NewCourseBlockService(repo)
	created, err := svc.Create(context.Background(), model.CourseBlock{CourseID: 1, BlockName: "Основы"})
	require.NoError(t, err)

	t.Run("deletes an existing item", func(t *testing.T) {
		require.NoError(t, svc.Delete(context.Background(), created.CourseID, created.ID))
		items, err := svc.ListByCourseID(context.Background(), created.CourseID)
		require.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("rejects a zero id", func(t *testing.T) {
		err := svc.Delete(context.Background(), 1, 0)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("does not delete a block belonging to a different course", func(t *testing.T) {
		block, err := svc.Create(context.Background(), model.CourseBlock{CourseID: 2, BlockName: "Course 2 block"})
		require.NoError(t, err)

		err = svc.Delete(context.Background(), 1, block.ID) // wrong courseId (1, not 2)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrNotFound)

		items, err := svc.ListByCourseID(context.Background(), 2)
		require.NoError(t, err)
		assert.Len(t, items, 1, "the block must still exist under its real course")
	})
}

func TestCourseBlockService_ListByCourseID(t *testing.T) {
	t.Run("filters by courseId", func(t *testing.T) {
		repo := newFakeCourseBlockRepository()
		svc := service.NewCourseBlockService(repo)

		_, err := svc.Create(context.Background(), model.CourseBlock{CourseID: 1, BlockName: "Course 1 Block 1"})
		require.NoError(t, err)
		_, err = svc.Create(context.Background(), model.CourseBlock{CourseID: 2, BlockName: "Course 2 Block 1"})
		require.NoError(t, err)
		_, err = svc.Create(context.Background(), model.CourseBlock{CourseID: 1, BlockName: "Course 1 Block 2"})
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
