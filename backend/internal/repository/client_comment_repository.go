package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type ClientCommentRepository struct {
	db *pgxpool.Pool
}

func NewClientCommentRepository(db *pgxpool.Pool) *ClientCommentRepository {
	return &ClientCommentRepository{db: db}
}

func (r *ClientCommentRepository) ListForClient(ctx context.Context, clientID int64) ([]model.ClientComment, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, client_id, text, created_at
		FROM client_comments
		WHERE client_id = $1
		ORDER BY created_at DESC
	`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.ClientComment{}
	for rows.Next() {
		var item model.ClientComment
		if err := rows.Scan(&item.ID, &item.ClientID, &item.Text, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *ClientCommentRepository) Create(ctx context.Context, item model.ClientComment) (model.ClientComment, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO client_comments (client_id, text)
		VALUES ($1, $2)
		RETURNING id, created_at
	`, item.ClientID, item.Text).Scan(&item.ID, &item.CreatedAt)
	return item, err
}
