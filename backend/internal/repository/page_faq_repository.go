package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type PageFAQRepository struct {
	db *pgxpool.Pool
}

func NewPageFAQRepository(db *pgxpool.Pool) *PageFAQRepository {
	return &PageFAQRepository{db: db}
}

func (r *PageFAQRepository) GetSettings(ctx context.Context, page string) (model.PageFAQSettings, error) {
	row := r.db.QueryRow(ctx, `
		SELECT page, title, description, visible, updated_at
		FROM page_faq_settings
		WHERE page = $1
	`, page)
	var item model.PageFAQSettings
	err := row.Scan(&item.Page, &item.Title, &item.Description, &item.Visible, &item.UpdatedAt)
	return item, translateNotFound(err)
}

// UpdateSettings never inserts — every valid page already has a row seeded
// by migration 00040, so an unknown page comes back as ErrNotFound rather
// than silently creating a new settings row for it.
func (r *PageFAQRepository) UpdateSettings(ctx context.Context, item model.PageFAQSettings) (model.PageFAQSettings, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE page_faq_settings
		SET title = $1, description = $2, visible = $3, updated_at = now()
		WHERE page = $4
		RETURNING updated_at
	`, item.Title, item.Description, item.Visible, item.Page).Scan(&item.UpdatedAt)
	return item, translateNotFound(err)
}

const pageFAQItemColumns = "id, page, question, answer, sort_order, created_at, updated_at"

func scanPageFAQItem(row pgx.Row) (model.PageFAQItem, error) {
	var item model.PageFAQItem
	err := row.Scan(&item.ID, &item.Page, &item.Question, &item.Answer, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *PageFAQRepository) ListByPage(ctx context.Context, page string) ([]model.PageFAQItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+pageFAQItemColumns+`
		FROM page_faq_items
		WHERE page = $1
		ORDER BY sort_order, id
	`, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.PageFAQItem{}
	for rows.Next() {
		item, err := scanPageFAQItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PageFAQRepository) CreateItem(ctx context.Context, item model.PageFAQItem) (model.PageFAQItem, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO page_faq_items (page, question, answer, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, item.Page, item.Question, item.Answer, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

// UpdateItem matches on id AND page — a URL like
// /page-faq/masterclasses/items/42 can't touch an item that actually
// belongs to a different page (mirrors CourseFAQRepository.Update).
func (r *PageFAQRepository) UpdateItem(ctx context.Context, item model.PageFAQItem) (model.PageFAQItem, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE page_faq_items
		SET question = $1, answer = $2, sort_order = $3, updated_at = now()
		WHERE id = $4 AND page = $5
		RETURNING page, updated_at
	`, item.Question, item.Answer, item.SortOrder, item.ID, item.Page).Scan(&item.Page, &item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *PageFAQRepository) DeleteItem(ctx context.Context, page string, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM page_faq_items WHERE id = $1 AND page = $2`, id, page)
	return checkDeleted(tag, err)
}
