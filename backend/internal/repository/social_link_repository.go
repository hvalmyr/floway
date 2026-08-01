package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type SocialLinkRepository struct {
	db *pgxpool.Pool
}

func NewSocialLinkRepository(db *pgxpool.Pool) *SocialLinkRepository {
	return &SocialLinkRepository{db: db}
}

func (r *SocialLinkRepository) List(ctx context.Context) ([]model.SocialLink, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, label, href, disclaimer, sort_order, created_at, updated_at
		FROM social_links
		ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.SocialLink{}
	for rows.Next() {
		var item model.SocialLink
		if err := rows.Scan(&item.ID, &item.Label, &item.Href, &item.Disclaimer, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *SocialLinkRepository) Create(ctx context.Context, item model.SocialLink) (model.SocialLink, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO social_links (label, href, disclaimer, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, item.Label, item.Href, item.Disclaimer, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *SocialLinkRepository) Update(ctx context.Context, item model.SocialLink) (model.SocialLink, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE social_links
		SET label = $1, href = $2, disclaimer = $3, sort_order = $4, updated_at = now()
		WHERE id = $5
		RETURNING updated_at
	`, item.Label, item.Href, item.Disclaimer, item.SortOrder, item.ID).Scan(&item.UpdatedAt)
	return item, err
}

func (r *SocialLinkRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM social_links WHERE id = $1`, id)
	return err
}
