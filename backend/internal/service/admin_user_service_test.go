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

type fakeAdminUserRepository struct {
	items  []model.AdminUser
	nextID int64
	err    error
}

func newFakeAdminUserRepository() *fakeAdminUserRepository {
	return &fakeAdminUserRepository{nextID: 1}
}

func (f *fakeAdminUserRepository) FindByLogin(ctx context.Context, login string) (model.AdminUser, error) {
	if f.err != nil {
		return model.AdminUser{}, f.err
	}
	for _, existing := range f.items {
		if existing.Login == login {
			return existing, nil
		}
	}
	return model.AdminUser{}, errors.New("not found")
}

func (f *fakeAdminUserRepository) FindByID(ctx context.Context, id int64) (model.AdminUser, error) {
	if f.err != nil {
		return model.AdminUser{}, f.err
	}
	for _, existing := range f.items {
		if existing.ID == id {
			return existing, nil
		}
	}
	return model.AdminUser{}, service.ErrNotFound
}

func (f *fakeAdminUserRepository) Create(ctx context.Context, login, passwordHash string) (model.AdminUser, error) {
	if f.err != nil {
		return model.AdminUser{}, f.err
	}
	user := model.AdminUser{ID: f.nextID, Login: login, PasswordHash: passwordHash}
	f.nextID++
	f.items = append(f.items, user)
	return user, nil
}

func (f *fakeAdminUserRepository) IncrementTokenVersion(ctx context.Context, id int64) error {
	if f.err != nil {
		return f.err
	}
	for i, existing := range f.items {
		if existing.ID == id {
			f.items[i].TokenVersion++
			return nil
		}
	}
	return service.ErrNotFound
}

func TestAdminUserService_Create(t *testing.T) {
	t.Run("creates a valid admin user with a hashed password", func(t *testing.T) {
		repo := newFakeAdminUserRepository()
		svc := service.NewAdminUserService(repo)

		user, err := svc.Create(context.Background(), "  admin  ", "supersecret")

		require.NoError(t, err)
		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, "admin", user.Login)
		assert.NotEqual(t, "supersecret", user.PasswordHash)
		assert.NotEmpty(t, user.PasswordHash)
	})

	t.Run("rejects an empty login", func(t *testing.T) {
		repo := newFakeAdminUserRepository()
		svc := service.NewAdminUserService(repo)

		_, err := svc.Create(context.Background(), "   ", "supersecret")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})

	t.Run("rejects a too short password", func(t *testing.T) {
		repo := newFakeAdminUserRepository()
		svc := service.NewAdminUserService(repo)

		_, err := svc.Create(context.Background(), "admin", "short")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})
}

func TestAdminUserService_Authenticate(t *testing.T) {
	repo := newFakeAdminUserRepository()
	svc := service.NewAdminUserService(repo)
	_, err := svc.Create(context.Background(), "admin", "supersecret")
	require.NoError(t, err)

	t.Run("authenticates with the correct password", func(t *testing.T) {
		user, err := svc.Authenticate(context.Background(), "admin", "supersecret")

		require.NoError(t, err)
		assert.Equal(t, "admin", user.Login)
	})

	t.Run("rejects an incorrect password", func(t *testing.T) {
		_, err := svc.Authenticate(context.Background(), "admin", "wrongpassword")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrInvalidCredentials)
	})

	t.Run("rejects an unknown login", func(t *testing.T) {
		_, err := svc.Authenticate(context.Background(), "nobody", "supersecret")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrInvalidCredentials)
	})
}

func TestAdminUserService_CurrentTokenVersion(t *testing.T) {
	repo := newFakeAdminUserRepository()
	svc := service.NewAdminUserService(repo)
	created, err := svc.Create(context.Background(), "admin", "supersecret")
	require.NoError(t, err)

	t.Run("returns the stored version for an existing user", func(t *testing.T) {
		version, err := svc.CurrentTokenVersion(context.Background(), created.ID)

		require.NoError(t, err)
		assert.Equal(t, 0, version)
	})

	t.Run("propagates the repository error for an unknown user", func(t *testing.T) {
		_, err := svc.CurrentTokenVersion(context.Background(), 999)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrNotFound)
	})
}

func TestAdminUserService_Logout(t *testing.T) {
	repo := newFakeAdminUserRepository()
	svc := service.NewAdminUserService(repo)
	created, err := svc.Create(context.Background(), "admin", "supersecret")
	require.NoError(t, err)

	t.Run("bumps the token version, invalidating tokens issued before it", func(t *testing.T) {
		before, err := svc.CurrentTokenVersion(context.Background(), created.ID)
		require.NoError(t, err)

		require.NoError(t, svc.Logout(context.Background(), created.ID))

		after, err := svc.CurrentTokenVersion(context.Background(), created.ID)
		require.NoError(t, err)
		assert.Equal(t, before+1, after)
	})

	t.Run("propagates the repository error for an unknown user", func(t *testing.T) {
		err := svc.Logout(context.Background(), 999)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrNotFound)
	})
}
