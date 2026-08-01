package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type AdminUserRepository struct {
	db *pgxpool.Pool
}

func NewAdminUserRepository(db *pgxpool.Pool) *AdminUserRepository {
	return &AdminUserRepository{db: db}
}

func (r *AdminUserRepository) FindByLogin(ctx context.Context, login string) (model.AdminUser, error) {
	var user model.AdminUser
	err := r.db.QueryRow(ctx, `
		SELECT id, login, password_hash, token_version, created_at
		FROM admin_users
		WHERE login = $1
	`, login).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.TokenVersion, &user.CreatedAt)
	return user, translateNotFound(err)
}

func (r *AdminUserRepository) FindByID(ctx context.Context, id int64) (model.AdminUser, error) {
	var user model.AdminUser
	err := r.db.QueryRow(ctx, `
		SELECT id, login, password_hash, token_version, created_at
		FROM admin_users
		WHERE id = $1
	`, id).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.TokenVersion, &user.CreatedAt)
	return user, translateNotFound(err)
}

func (r *AdminUserRepository) Create(ctx context.Context, login, passwordHash string) (model.AdminUser, error) {
	user := model.AdminUser{Login: login, PasswordHash: passwordHash}
	err := r.db.QueryRow(ctx, `
		INSERT INTO admin_users (login, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at
	`, login, passwordHash).Scan(&user.ID, &user.CreatedAt)
	return user, err
}

// IncrementTokenVersion invalidates every JWT issued to this user before
// now — the caller's requireAdminMiddleware rejects any token whose
// embedded version no longer matches. Used on logout.
func (r *AdminUserRepository) IncrementTokenVersion(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `UPDATE admin_users SET token_version = token_version + 1 WHERE id = $1`, id)
	return checkDeleted(tag, err)
}
