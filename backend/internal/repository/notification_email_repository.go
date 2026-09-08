package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type NotificationEmailRepository struct {
	db *pgxpool.Pool
}

func NewNotificationEmailRepository(db *pgxpool.Pool) *NotificationEmailRepository {
	return &NotificationEmailRepository{db: db}
}

func (r *NotificationEmailRepository) List(ctx context.Context) ([]model.NotificationEmail, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, email, created_at
		FROM notification_emails
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.NotificationEmail{}
	for rows.Next() {
		var item model.NotificationEmail
		if err := rows.Scan(&item.ID, &item.Email, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *NotificationEmailRepository) Create(ctx context.Context, item model.NotificationEmail) (model.NotificationEmail, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO notification_emails (email)
		VALUES ($1)
		RETURNING id, created_at
	`, item.Email).Scan(&item.ID, &item.CreatedAt)
	return item, err
}

func (r *NotificationEmailRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM notification_emails WHERE id = $1`, id)
	return checkDeleted(tag, err)
}
