package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type fakeCourseRepository struct {
	items  []model.Course
	nextID int64
	err    error
}

func newFakeCourseRepository() *fakeCourseRepository {
	return &fakeCourseRepository{nextID: 1}
}

func (f *fakeCourseRepository) List(ctx context.Context) ([]model.Course, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func (f *fakeCourseRepository) FindBySlug(ctx context.Context, slug string) (model.Course, error) {
	if f.err != nil {
		return model.Course{}, f.err
	}
	for _, existing := range f.items {
		if existing.Slug == slug {
			return existing, nil
		}
	}
	return model.Course{}, pgx.ErrNoRows
}

func (f *fakeCourseRepository) Create(ctx context.Context, item model.Course) (model.Course, error) {
	if f.err != nil {
		return model.Course{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeCourseRepository) Update(ctx context.Context, item model.Course) (model.Course, error) {
	if f.err != nil {
		return model.Course{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.Course{}, errors.New("not found")
}

func (f *fakeCourseRepository) Delete(ctx context.Context, id int64) error {
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

func TestCourseService_Create(t *testing.T) {
	t.Run("creates a valid course", func(t *testing.T) {
		repo := newFakeCourseRepository()
		svc := service.NewCourseService(repo)

		course, err := svc.Create(context.Background(), model.Course{
			Slug:  "  go-basics  ",
			Title: "  Go Basics  ",
		})

		require.NoError(t, err)
		assert.Equal(t, int64(1), course.ID)
		assert.Equal(t, "go-basics", course.Slug)
		assert.Equal(t, "Go Basics", course.Title)
	})

	t.Run("rejects an empty slug", func(t *testing.T) {
		repo := newFakeCourseRepository()
		svc := service.NewCourseService(repo)

		_, err := svc.Create(context.Background(), model.Course{Slug: "   ", Title: "Title"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})

	t.Run("rejects an empty title", func(t *testing.T) {
		repo := newFakeCourseRepository()
		svc := service.NewCourseService(repo)

		_, err := svc.Create(context.Background(), model.Course{Slug: "slug", Title: "   "})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects an invalid status", func(t *testing.T) {
		repo := newFakeCourseRepository()
		svc := service.NewCourseService(repo)

		_, err := svc.Create(context.Background(), model.Course{
			Slug:   "slug",
			Title:  "Title",
			Status: model.CourseStatus("draft"),
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("defaults status to active when empty", func(t *testing.T) {
		repo := newFakeCourseRepository()
		svc := service.NewCourseService(repo)

		course, err := svc.Create(context.Background(), model.Course{Slug: "slug", Title: "Title"})

		require.NoError(t, err)
		assert.Equal(t, model.CourseStatusActive, course.Status)
	})

	t.Run("stores gallery as provided", func(t *testing.T) {
		repo := newFakeCourseRepository()
		svc := service.NewCourseService(repo)

		course, err := svc.Create(context.Background(), model.Course{
			Slug:    "slug",
			Title:   "Title",
			Gallery: []string{"one.jpg", "two.jpg"},
		})

		require.NoError(t, err)
		assert.Equal(t, []string{"one.jpg", "two.jpg"}, course.Gallery)
	})
}

func TestCourseService_Update(t *testing.T) {
	repo := newFakeCourseRepository()
	svc := service.NewCourseService(repo)
	created, err := svc.Create(context.Background(), model.Course{Slug: "slug", Title: "Title"})
	require.NoError(t, err)

	t.Run("updates an existing course", func(t *testing.T) {
		created.Title = "Updated Title"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "Updated Title", updated.Title)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.Course{Slug: "slug", Title: "Title"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestCourseService_Delete(t *testing.T) {
	repo := newFakeCourseRepository()
	svc := service.NewCourseService(repo)
	created, err := svc.Create(context.Background(), model.Course{Slug: "slug", Title: "Title"})
	require.NoError(t, err)

	t.Run("deletes an existing course", func(t *testing.T) {
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

func TestCourseService_GetBySlug(t *testing.T) {
	repo := newFakeCourseRepository()
	svc := service.NewCourseService(repo)
	created, err := svc.Create(context.Background(), model.Course{Slug: "slug", Title: "Title"})
	require.NoError(t, err)

	t.Run("finds an existing course", func(t *testing.T) {
		found, err := svc.GetBySlug(context.Background(), "slug")

		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
	})

	t.Run("returns ErrNotFound for an unknown slug", func(t *testing.T) {
		_, err := svc.GetBySlug(context.Background(), "does-not-exist")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrNotFound)
	})
}

func TestCourseService_List_PropagatesRepositoryError(t *testing.T) {
	repo := newFakeCourseRepository()
	repo.err = errors.New("boom")
	svc := service.NewCourseService(repo)

	_, err := svc.List(context.Background())

	require.Error(t, err)
}
