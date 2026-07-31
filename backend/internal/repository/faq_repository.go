package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type FAQRepository struct {
	db *pgxpool.Pool
}

func NewFAQRepository(db *pgxpool.Pool) *FAQRepository {
	return &FAQRepository{db: db}
}

func (r *FAQRepository) List(ctx context.Context) ([]model.FAQItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, question, answer, sort_order, created_at, updated_at
		FROM faq_items
		ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.FAQItem{}
	for rows.Next() {
		var item model.FAQItem
		if err := rows.Scan(&item.ID, &item.Question, &item.Answer, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *FAQRepository) Create(ctx context.Context, item model.FAQItem) (model.FAQItem, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO faq_items (question, answer, sort_order)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`, item.Question, item.Answer, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *FAQRepository) Update(ctx context.Context, item model.FAQItem) (model.FAQItem, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE faq_items
		SET question = $1, answer = $2, sort_order = $3, updated_at = now()
		WHERE id = $4
		RETURNING updated_at
	`, item.Question, item.Answer, item.SortOrder, item.ID).Scan(&item.UpdatedAt)
	return item, err
}

func (r *FAQRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM faq_items WHERE id = $1`, id)
	return err
}
