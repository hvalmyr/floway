package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type LandingPageFAQRepository struct {
	db *pgxpool.Pool
}

func NewLandingPageFAQRepository(db *pgxpool.Pool) *LandingPageFAQRepository {
	return &LandingPageFAQRepository{db: db}
}

const landingPageFAQColumns = "id, landing_page_id, question, answer, sort_order, created_at, updated_at"

func scanLandingPageFAQItem(row pgx.Row) (model.LandingPageFAQItem, error) {
	var item model.LandingPageFAQItem
	err := row.Scan(&item.ID, &item.LandingPageID, &item.Question, &item.Answer, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *LandingPageFAQRepository) ListByLandingPageID(ctx context.Context, landingPageID int64) ([]model.LandingPageFAQItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+landingPageFAQColumns+`
		FROM landing_page_faq_items
		WHERE landing_page_id = $1
		ORDER BY sort_order, id
	`, landingPageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.LandingPageFAQItem{}
	for rows.Next() {
		item, err := scanLandingPageFAQItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *LandingPageFAQRepository) Create(ctx context.Context, item model.LandingPageFAQItem) (model.LandingPageFAQItem, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO landing_page_faq_items (landing_page_id, question, answer, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, item.LandingPageID, item.Question, item.Answer, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

// Update matches on id AND landing_page_id — mirrors CourseFAQRepository.Update.
func (r *LandingPageFAQRepository) Update(ctx context.Context, item model.LandingPageFAQItem) (model.LandingPageFAQItem, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE landing_page_faq_items
		SET question = $1, answer = $2, sort_order = $3, updated_at = now()
		WHERE id = $4 AND landing_page_id = $5
		RETURNING landing_page_id, updated_at
	`, item.Question, item.Answer, item.SortOrder, item.ID, item.LandingPageID).Scan(&item.LandingPageID, &item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *LandingPageFAQRepository) Delete(ctx context.Context, landingPageID, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM landing_page_faq_items WHERE id = $1 AND landing_page_id = $2`, id, landingPageID)
	return checkDeleted(tag, err)
}
