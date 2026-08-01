package service

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"floway-backend/internal/model"
)

const minPasswordLength = 8

type AdminUserRepository interface {
	FindByLogin(ctx context.Context, login string) (model.AdminUser, error)
	FindByID(ctx context.Context, id int64) (model.AdminUser, error)
	Create(ctx context.Context, login, passwordHash string) (model.AdminUser, error)
	IncrementTokenVersion(ctx context.Context, id int64) error
}

type AdminUserService struct {
	repo AdminUserRepository
}

func NewAdminUserService(repo AdminUserRepository) *AdminUserService {
	return &AdminUserService{repo: repo}
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

// CurrentTokenVersion is what requireAdminMiddleware compares an incoming
// JWT's embedded version against, on every admin request.
func (s *AdminUserService) CurrentTokenVersion(ctx context.Context, userID int64) (int, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return user.TokenVersion, nil
}

// Logout invalidates every JWT already issued to this admin — not just the
// one presented at logout time — by bumping their stored token version.
// There's no per-token denylist (no jti store), so this is a coarser but
// simpler guarantee: "log out" means "none of my old sessions work anymore."
func (s *AdminUserService) Logout(ctx context.Context, userID int64) error {
	return s.repo.IncrementTokenVersion(ctx, userID)
}
