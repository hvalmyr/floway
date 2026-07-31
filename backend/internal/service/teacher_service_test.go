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

type fakeTeacherRepository struct {
	items  []model.Teacher
	nextID int64
	err    error
}

func newFakeTeacherRepository() *fakeTeacherRepository {
	return &fakeTeacherRepository{nextID: 1}
}

func (f *fakeTeacherRepository) List(ctx context.Context) ([]model.Teacher, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func (f *fakeTeacherRepository) Create(ctx context.Context, item model.Teacher) (model.Teacher, error) {
	if f.err != nil {
		return model.Teacher{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeTeacherRepository) Update(ctx context.Context, item model.Teacher) (model.Teacher, error) {
	if f.err != nil {
		return model.Teacher{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == item.ID {
			f.items[i] = item
			return item, nil
		}
	}
	return model.Teacher{}, errors.New("not found")
}

func (f *fakeTeacherRepository) Delete(ctx context.Context, id int64) error {
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

func TestTeacherService_Create(t *testing.T) {
	t.Run("creates a valid item", func(t *testing.T) {
		repo := newFakeTeacherRepository()
		svc := service.NewTeacherService(repo)

		item, err := svc.Create(context.Background(), model.Teacher{Name: "  Анна Иванова  "})

		require.NoError(t, err)
		assert.Equal(t, int64(1), item.ID)
		assert.Equal(t, "Анна Иванова", item.Name)
	})

	t.Run("rejects an empty name", func(t *testing.T) {
		repo := newFakeTeacherRepository()
		svc := service.NewTeacherService(repo)

		_, err := svc.Create(context.Background(), model.Teacher{Name: "   "})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})
}

func TestTeacherService_Update(t *testing.T) {
	repo := newFakeTeacherRepository()
	svc := service.NewTeacherService(repo)
	created, err := svc.Create(context.Background(), model.Teacher{Name: "Анна"})
	require.NoError(t, err)

	t.Run("updates an existing item", func(t *testing.T) {
		created.Name = "Мария"
		updated, err := svc.Update(context.Background(), created)

		require.NoError(t, err)
		assert.Equal(t, "Мария", updated.Name)
	})

	t.Run("rejects a missing id", func(t *testing.T) {
		_, err := svc.Update(context.Background(), model.Teacher{Name: "Мария"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestTeacherService_Delete(t *testing.T) {
	repo := newFakeTeacherRepository()
	svc := service.NewTeacherService(repo)
	created, err := svc.Create(context.Background(), model.Teacher{Name: "Анна"})
	require.NoError(t, err)

	t.Run("deletes an existing item", func(t *testing.T) {
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

func TestTeacherService_List_PropagatesRepositoryError(t *testing.T) {
	repo := newFakeTeacherRepository()
	repo.err = errors.New("boom")
	svc := service.NewTeacherService(repo)

	_, err := svc.List(context.Background())

	require.Error(t, err)
}
