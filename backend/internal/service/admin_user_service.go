package service

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"floway-backend/internal/model"
)

// ErrInvalidCredentials is returned by Authenticate when the login does not
// exist or the password does not match. It intentionally does not reveal
// which of the two is the case.
var ErrInvalidCredentials = errors.New("invalid credentials")

const minPasswordLength = 8

type AdminUserRepository interface {
	List(ctx context.Context) ([]model.AdminUser, error)
	FindByLogin(ctx context.Context, login string) (model.AdminUser, error)
	Create(ctx context.Context, login, passwordHash string) (model.AdminUser, error)
	Delete(ctx context.Context, id int64) error
}

type AdminUserService struct {
	repo AdminUserRepository
}

func NewAdminUserService(repo AdminUserRepository) *AdminUserService {
	return &AdminUserService{repo: repo}
}

func (s *AdminUserService) List(ctx context.Context) ([]model.AdminUser, error) {
	return s.repo.List(ctx)
}

func (s *AdminUserService) Create(ctx context.Context, login, plainPassword string) (model.AdminUser, error) {
	login = strings.TrimSpace(login)
	if login == "" {
		return model.AdminUser{}, errors.Join(ErrValidation, errors.New("login is required"))
	}
	if len(plainPassword) < minPasswordLength {
		return model.AdminUser{}, errors.Join(ErrValidation, errors.New("password must be at least 8 characters"))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return model.AdminUser{}, err
	}

	return s.repo.Create(ctx, login, string(hash))
}

// Authenticate verifies login credentials for the admin panel. On any
// failure (unknown login or wrong password) it returns ErrInvalidCredentials
// without disclosing which part of the credentials was incorrect.
func (s *AdminUserService) Authenticate(ctx context.Context, login, plainPassword string) (model.AdminUser, error) {
	user, err := s.repo.FindByLogin(ctx, login)
	if err != nil {
		return model.AdminUser{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(plainPassword)); err != nil {
		return model.AdminUser{}, ErrInvalidCredentials
	}

	return user, nil
}

func (s *AdminUserService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
