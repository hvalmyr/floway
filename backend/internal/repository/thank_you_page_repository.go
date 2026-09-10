package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

// ThankYouPageRepository bundles the settings row plus its two child lists
// (photos, FAQ items) for one thank-you page variant — mirrors
// PageFAQRepository bundling settings+items, extended with a second list.
type ThankYouPageRepository struct {
	db *pgxpool.Pool
}

func NewThankYouPageRepository(db *pgxpool.Pool) *ThankYouPageRepository {
	return &ThankYouPageRepository{db: db}
}

const thankYouPageColumns = "variant, title, subtitle, description, show_messengers, show_social_links, show_blog_link, blog_link_text, blog_link_url, show_carousel, show_faq, show_community, community_text, community_url, updated_at"

func scanThankYouPage(row pgx.Row) (model.ThankYouPage, error) {
	var item model.ThankYouPage
	err := row.Scan(&item.Variant, &item.Title, &item.Subtitle, &item.Description, &item.ShowMessengers, &item.ShowSocialLinks, &item.ShowBlogLink, &item.BlogLinkText, &item.BlogLinkURL, &item.ShowCarousel, &item.ShowFAQ, &item.ShowCommunity, &item.CommunityText, &item.CommunityURL, &item.UpdatedAt)
	return item, err
}

func (r *ThankYouPageRepository) GetSettings(ctx context.Context, variant string) (model.ThankYouPage, error) {
	row := r.db.QueryRow(ctx, `SELECT `+thankYouPageColumns+` FROM thank_you_pages WHERE variant = $1`, variant)
	item, err := scanThankYouPage(row)
	return item, translateNotFound(err)
}

// UpdateSettings never inserts — every valid variant already has a row
// seeded by migration 00047, so an unknown variant comes back as
// ErrNotFound rather than silently creating a new settings row for it
// (mirrors PageFAQRepository.UpdateSettings).
func (r *ThankYouPageRepository) UpdateSettings(ctx context.Context, item model.ThankYouPage) (model.ThankYouPage, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE thank_you_pages
		SET title = $1, subtitle = $2, description = $3, show_messengers = $4, show_social_links = $5,
		    show_blog_link = $6, blog_link_text = $7, blog_link_url = $8, show_carousel = $9,
		    show_faq = $10, show_community = $11, community_text = $12, community_url = $13, updated_at = now()
		WHERE variant = $14
		RETURNING `+thankYouPageColumns,
		item.Title, item.Subtitle, item.Description, item.ShowMessengers, item.ShowSocialLinks,
		item.ShowBlogLink, item.BlogLinkText, item.BlogLinkURL, item.ShowCarousel,
		item.ShowFAQ, item.ShowCommunity, item.CommunityText, item.CommunityURL, item.Variant,
	)
	updated, err := scanThankYouPage(row)
	return updated, translateNotFound(err)
}

const thankYouPagePhotoColumns = "id, variant, image, sort_order, created_at, updated_at"

func scanThankYouPagePhoto(row pgx.Row) (model.ThankYouPagePhoto, error) {
	var item model.ThankYouPagePhoto
	err := row.Scan(&item.ID, &item.Variant, &item.Image, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *ThankYouPageRepository) ListPhotos(ctx context.Context, variant string) ([]model.ThankYouPagePhoto, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+thankYouPagePhotoColumns+`
		FROM thank_you_page_photos
		WHERE variant = $1
		ORDER BY sort_order, id
	`, variant)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.ThankYouPagePhoto{}
	for rows.Next() {
		item, err := scanThankYouPagePhoto(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *ThankYouPageRepository) CreatePhoto(ctx context.Context, item model.ThankYouPagePhoto) (model.ThankYouPagePhoto, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO thank_you_page_photos (variant, image, sort_order)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`, item.Variant, item.Image, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

// UpdatePhoto matches on id AND variant — mirrors PageFAQRepository.UpdateItem.
func (r *ThankYouPageRepository) UpdatePhoto(ctx context.Context, item model.ThankYouPagePhoto) (model.ThankYouPagePhoto, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE thank_you_page_photos
		SET image = $1, sort_order = $2, updated_at = now()
		WHERE id = $3 AND variant = $4
		RETURNING variant, updated_at
	`, item.Image, item.SortOrder, item.ID, item.Variant).Scan(&item.Variant, &item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *ThankYouPageRepository) DeletePhoto(ctx context.Context, variant string, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM thank_you_page_photos WHERE id = $1 AND variant = $2`, id, variant)
	return checkDeleted(tag, err)
}

const thankYouPageFAQItemColumns = "id, variant, question, answer, sort_order, created_at, updated_at"

func scanThankYouPageFAQItem(row pgx.Row) (model.ThankYouPageFAQItem, error) {
	var item model.ThankYouPageFAQItem
	err := row.Scan(&item.ID, &item.Variant, &item.Question, &item.Answer, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *ThankYouPageRepository) ListFAQItems(ctx context.Context, variant string) ([]model.ThankYouPageFAQItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+thankYouPageFAQItemColumns+`
		FROM thank_you_page_faq_items
		WHERE variant = $1
		ORDER BY sort_order, id
	`, variant)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.ThankYouPageFAQItem{}
	for rows.Next() {
		item, err := scanThankYouPageFAQItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *ThankYouPageRepository) CreateFAQItem(ctx context.Context, item model.ThankYouPageFAQItem) (model.ThankYouPageFAQItem, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO thank_you_page_faq_items (variant, question, answer, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, item.Variant, item.Question, item.Answer, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *ThankYouPageRepository) UpdateFAQItem(ctx context.Context, item model.ThankYouPageFAQItem) (model.ThankYouPageFAQItem, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE thank_you_page_faq_items
		SET question = $1, answer = $2, sort_order = $3, updated_at = now()
		WHERE id = $4 AND variant = $5
		RETURNING variant, updated_at
	`, item.Question, item.Answer, item.SortOrder, item.ID, item.Variant).Scan(&item.Variant, &item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *ThankYouPageRepository) DeleteFAQItem(ctx context.Context, variant string, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM thank_you_page_faq_items WHERE id = $1 AND variant = $2`, id, variant)
	return checkDeleted(tag, err)
}
