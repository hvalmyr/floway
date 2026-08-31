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

type fakeCourseRepository struct {
	items  []model.Course
	nextID int64
	err    error
}

func newFakeCourseRepository() *fakeCourseRepository {
	return &fakeCourseRepository{nextID: 1}
}

func (f *fakeCourseRepository) ListBySectionID(ctx context.Context, sectionID int64) ([]model.Course, error) {
	if f.err != nil {
		return nil, f.err
	}
	var items []model.Course
	for _, existing := range f.items {
		if existing.SectionID == sectionID {
			items = append(items, existing)
		}
	}
	return items, nil
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
	return model.Course{}, service.ErrNotFound
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

// Update and Delete only match a row that belongs to the given section —
// mirrors the real repository's WHERE id = $x AND section_id = $y, so a URL
// like /course-sections/1/courses/42 can't touch a course that actually
// belongs to a different section.
func (f *fakeCourseRepository) Update(ctx context.Context, item model.Course) (model.Course, error) {
	if f.err != nil {
		return model.Course{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID && existing.SectionID == item.SectionID {
			item.SectionID = existing.SectionID
			f.items[i] = item
			return item, nil
		}
	}
	return model.Course{}, service.ErrNotFound
}

func (f *fakeCourseRepository) Delete(ctx context.Context, sectionID, id int64) error {
	if f.err != nil {
		return f.err
	}
	for i, existing := range f.items {
		if existing.ID == id && existing.SectionID == sectionID {
			f.items = append(f.items[:i], f.items[i+1:]...)
			return nil
		}
	}
	return service.ErrNotFound
}

func TestCourseService_Create(t *testing.T) {
	t.Run("creates a valid course", func(t *testing.T) {
		repo := newFakeCourseRepository()
		svc := service.NewCourseService(repo)

		course, err := svc.Create(context.Background(), model.Course{
			SectionID: 1,
			Slug:      "  go-basics  ",
			Name:      "  Go Basics  ",
		})

		require.NoError(t, err)
		assert.Equal(t, int64(1), course.ID)
		assert.Equal(t, "go-basics", course.Slug)
		assert.Equal(t, "Go Basics", course.Name)
	})

	t.Run("rejects an empty slug", func(t *testing.T) {
		repo := newFakeCourseRepository()
		svc := service.NewCourseService(repo)

		_, err := svc.Create(context.Background(), model.Course{SectionID: 1, Slug: "   ", Name: "Name"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})

	t.Run("rejects an empty name", func(t *testing.T) {
		repo := newFakeCourseRepository()
		svc := service.NewCourseService(repo)

		_, err := svc.Create(context.Background(), model.Course{SectionID: 1, Slug: "slug", Name: "   "})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects a zero sectionId", func(t *testing.T) {
		repo := newFakeCourseRepository()
		svc := service.NewCourseService(repo)

		_, err := svc.Create(context.Background(), model.Course{Slug: "slug", Name: "Name"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})
}

func TestCourseService_Update(t *testing.T) {
	repo := newFakeCourseRepository()
	svc := service.NewCourseService(repo)
	created, err := svc.Create(context.Background(), model.Course{SectionID: 1, Slug: "slug", Name: "Name"})
	require.NoError(t, err)

	t.Run("updates an existing course", func(t *testing.T) {
		created.Name = "Updated Name"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "Updated Name", updated.Name)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.Course{SectionID: 1, Slug: "slug", Name: "Name"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects a missing sectionId", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.Course{ID: created.ID, Slug: "slug", Name: "Name"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("does not update a course belonging to a different section", func(t *testing.T) {
		other, err := svc.Create(context.Background(), model.Course{SectionID: 2, Slug: "other", Name: "Other"})
		require.NoError(t, err)

		// Same course id as `other`, but wrong sectionId — simulates
		// PUT /course-sections/1/courses/{other.ID}.
		_, err = svc.Update(context.Background(), model.Course{ID: other.ID, SectionID: 1, Slug: "hijacked", Name: "Hijacked"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrNotFound)
	})
}

func TestCourseService_Delete(t *testing.T) {
	repo := newFakeCourseRepository()
	svc := service.NewCourseService(repo)
	created, err := svc.Create(context.Background(), model.Course{SectionID: 1, Slug: "slug", Name: "Name"})
	require.NoError(t, err)

	t.Run("deletes an existing course", func(t *testing.T) {
		require.NoError(t, svc.Delete(context.Background(), created.SectionID, created.ID))
		items, err := svc.ListBySectionID(context.Background(), created.SectionID)
		require.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("rejects a zero id", func(t *testing.T) {
		err := svc.Delete(context.Background(), 1, 0)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestCourseService_GetBySlug(t *testing.T) {
	repo := newFakeCourseRepository()
	svc := service.NewCourseService(repo)
	created, err := svc.Create(context.Background(), model.Course{SectionID: 1, Slug: "slug", Name: "Name"})
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

func TestCourseService_ListBySectionID_PropagatesRepositoryError(t *testing.T) {
	repo := newFakeCourseRepository()
	repo.err = errors.New("boom")
	svc := service.NewCourseService(repo)

	_, err := svc.ListBySectionID(context.Background(), 1)

	require.Error(t, err)
}
