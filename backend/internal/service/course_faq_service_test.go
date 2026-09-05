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

type fakeCourseFAQRepository struct {
	items  []model.CourseFAQItem
	nextID int64
	err    error
}

func newFakeCourseFAQRepository() *fakeCourseFAQRepository {
	return &fakeCourseFAQRepository{nextID: 1}
}

func (f *fakeCourseFAQRepository) ListByCourseID(ctx context.Context, courseID int64) ([]model.CourseFAQItem, error) {
	if f.err != nil {
		return nil, f.err
	}
	var items []model.CourseFAQItem
	for _, item := range f.items {
		if item.CourseID == courseID {
			items = append(items, item)
		}
	}
	return items, nil
}

func (f *fakeCourseFAQRepository) Create(ctx context.Context, item model.CourseFAQItem) (model.CourseFAQItem, error) {
	if f.err != nil {
		return model.CourseFAQItem{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

// Update only matches a row that belongs to the given course — mirrors the
// real repository's WHERE id = $x AND course_id = $y.
func (f *fakeCourseFAQRepository) Update(ctx context.Context, item model.CourseFAQItem) (model.CourseFAQItem, error) {
	if f.err != nil {
		return model.CourseFAQItem{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID && existing.CourseID == item.CourseID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.CourseFAQItem{}, service.ErrNotFound
}

func (f *fakeCourseFAQRepository) Delete(ctx context.Context, courseID, id int64) error {
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

func TestCourseFAQService_Create(t *testing.T) {
	t.Run("creates a valid item", func(t *testing.T) {
		repo := newFakeCourseFAQRepository()
		svc := service.NewCourseFAQService(repo)

		item, err := svc.Create(context.Background(), model.CourseFAQItem{
			CourseID: 1,
			Question: "  Сколько длится курс?  ",
			Answer:   "  3 месяца ",
		})

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
		assert.Equal(t, "Сколько длится курс?", item.Question)
		assert.Equal(t, "3 месяца", item.Answer)
	})

	t.Run("rejects a zero courseId", func(t *testing.T) {
		repo := newFakeCourseFAQRepository()
		svc := service.NewCourseFAQService(repo)

		_, err := svc.Create(context.Background(), model.CourseFAQItem{Question: "q", Answer: "a"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})

	t.Run("rejects an empty question or answer", func(t *testing.T) {
		repo := newFakeCourseFAQRepository()
		svc := service.NewCourseFAQService(repo)

		_, err := svc.Create(context.Background(), model.CourseFAQItem{CourseID: 1, Question: "  ", Answer: "a"})
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)

		_, err = svc.Create(context.Background(), model.CourseFAQItem{CourseID: 1, Question: "q", Answer: "  "})
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestCourseFAQService_Update(t *testing.T) {
	repo := newFakeCourseFAQRepository()
	svc := service.NewCourseFAQService(repo)
	created, err := svc.Create(context.Background(), model.CourseFAQItem{CourseID: 1, Question: "q", Answer: "a"})
	require.NoError(t, err)

	t.Run("updates an existing item", func(t *testing.T) {
		created.Question = "updated question"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "updated question", updated.Question)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.CourseFAQItem{CourseID: 1, Question: "q", Answer: "a"})
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects a missing courseId", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.CourseFAQItem{ID: created.ID, Question: "q", Answer: "a"})
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("does not update an item belonging to a different course", func(t *testing.T) {
		other, err := svc.Create(context.Background(), model.CourseFAQItem{CourseID: 2, Question: "other course item", Answer: "a"})
		require.NoError(t, err)

		_, err = svc.Update(context.Background(), model.CourseFAQItem{ID: other.ID, CourseID: 1, Question: "hijacked", Answer: "a"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrNotFound)
	})
}

func TestCourseFAQService_Delete(t *testing.T) {
	repo := newFakeCourseFAQRepository()
	svc := service.NewCourseFAQService(repo)
	created, err := svc.Create(context.Background(), model.CourseFAQItem{CourseID: 1, Question: "q", Answer: "a"})
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

	t.Run("does not delete an item belonging to a different course", func(t *testing.T) {
		item, err := svc.Create(context.Background(), model.CourseFAQItem{CourseID: 2, Question: "other course item", Answer: "a"})
		require.NoError(t, err)

		err = svc.Delete(context.Background(), 1, item.ID) // wrong courseId (1, not 2)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrNotFound)

		items, err := svc.ListByCourseID(context.Background(), 2)
		require.NoError(t, err)
		assert.Len(t, items, 1, "the item must still exist under its real course")
	})
}

func TestCourseFAQService_ListByCourseID_PropagatesRepositoryError(t *testing.T) {
	repo := newFakeCourseFAQRepository()
	repo.err = errors.New("boom")
	svc := service.NewCourseFAQService(repo)

	_, err := svc.ListByCourseID(context.Background(), 1)

	require.Error(t, err)
}
