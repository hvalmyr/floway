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

func (r *AdminUserRepository) List(ctx context.Context) ([]model.AdminUser, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, login, password_hash, created_at
		FROM admin_users
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.AdminUser
	for rows.Next() {
		var user model.AdminUser
		if err := rows.Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *AdminUserRepository) FindByLogin(ctx context.Context, login string) (model.AdminUser, error) {
	var user model.AdminUser
	err := r.db.QueryRow(ctx, `
		SELECT id, login, password_hash, created_at
		FROM admin_users
		WHERE login = $1
	`, login).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	return user, err
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

func (r *AdminUserRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM admin_users WHERE id = $1`, id)
	return err
}
