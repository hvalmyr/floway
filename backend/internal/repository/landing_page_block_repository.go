package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type LandingPageBlockRepository struct {
	db *pgxpool.Pool
}

func NewLandingPageBlockRepository(db *pgxpool.Pool) *LandingPageBlockRepository {
	return &LandingPageBlockRepository{db: db}
}

const landingPageBlockColumns = "id, landing_page_id, image, text, sort_order, created_at, updated_at"

func scanLandingPageBlock(row pgx.Row) (model.LandingPageBlock, error) {
	var item model.LandingPageBlock
	err := row.Scan(&item.ID, &item.LandingPageID, &item.Image, &item.Text, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *LandingPageBlockRepository) ListByLandingPageID(ctx context.Context, landingPageID int64) ([]model.LandingPageBlock, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+landingPageBlockColumns+`
		FROM landing_page_blocks
		WHERE landing_page_id = $1
		ORDER BY sort_order, id
	`, landingPageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.LandingPageBlock{}
	for rows.Next() {
		item, err := scanLandingPageBlock(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *LandingPageBlockRepository) Create(ctx context.Context, item model.LandingPageBlock) (model.LandingPageBlock, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO landing_page_blocks (landing_page_id, image, text, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, item.LandingPageID, item.Image, item.Text, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

// Update matches on id AND landing_page_id — a URL like
// /landing-pages/7/blocks/42 can't touch a block that actually belongs to a
// different landing page (mirrors CourseBlockRepository.Update).
func (r *LandingPageBlockRepository) Update(ctx context.Context, item model.LandingPageBlock) (model.LandingPageBlock, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE landing_page_blocks
		SET image = $1, text = $2, sort_order = $3, updated_at = now()
		WHERE id = $4 AND landing_page_id = $5
		RETURNING landing_page_id, updated_at
	`, item.Image, item.Text, item.SortOrder, item.ID, item.LandingPageID).Scan(&item.LandingPageID, &item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *LandingPageBlockRepository) Delete(ctx context.Context, landingPageID, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM landing_page_blocks WHERE id = $1 AND landing_page_id = $2`, id, landingPageID)
	return checkDeleted(tag, err)
}
