package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type LandingPageRepository struct {
	db *pgxpool.Pool
}

func NewLandingPageRepository(db *pgxpool.Pool) *LandingPageRepository {
	return &LandingPageRepository{db: db}
}

const landingPageColumns = "id, object_type_id, slug, h1, meta_title, meta_description, faq_title, faq_description, faq_visible, visible, sort_order, created_at, updated_at"

const landingPageColumnsQualified = "lp.id, lp.object_type_id, lp.slug, lp.h1, lp.meta_title, lp.meta_description, lp.faq_title, lp.faq_description, lp.faq_visible, lp.visible, lp.sort_order, lp.created_at, lp.updated_at"

func scanLandingPage(row pgx.Row) (model.LandingPage, error) {
	var item model.LandingPage
	err := row.Scan(&item.ID, &item.ObjectTypeID, &item.Slug, &item.H1, &item.MetaTitle, &item.MetaDescription, &item.FAQTitle, &item.FAQDescription, &item.FAQVisible, &item.Visible, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

// List returns every landing page, hidden ones included — the admin screen
// needs to see and re-show/copy/delete them.
func (r *LandingPageRepository) List(ctx context.Context) ([]model.LandingPage, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+landingPageColumns+`
		FROM landing_pages
		ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.LandingPage{}
	for rows.Next() {
		item, err := scanLandingPage(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *LandingPageRepository) FindBySlug(ctx context.Context, slug string) (model.LandingPage, error) {
	row := r.db.QueryRow(ctx, `
		SELECT `+landingPageColumns+`
		FROM landing_pages
		WHERE slug = $1
	`, slug)
	item, err := scanLandingPage(row)
	return item, translateNotFound(err)
}

// ListVisibleWithObjectType is the public homepage/nav aggregation — every
// visible landing page joined with its object type, one query (the row
// count here is tiny, a handful of pages, so this isn't the batched-query
// pattern CourseCatalogService needs for a deeper tree).
func (r *LandingPageRepository) ListVisibleWithObjectType(ctx context.Context) ([]model.LandingPageWithObjectType, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+landingPageColumnsQualified+`, `+objectTypeColumnsQualified+`
		FROM landing_pages lp
		JOIN object_types ot ON ot.id = lp.object_type_id
		WHERE lp.visible = true
		ORDER BY lp.sort_order, lp.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.LandingPageWithObjectType{}
	for rows.Next() {
		var item model.LandingPageWithObjectType
		if err := rows.Scan(
			&item.ID, &item.ObjectTypeID, &item.Slug, &item.H1, &item.MetaTitle, &item.MetaDescription, &item.FAQTitle, &item.FAQDescription, &item.FAQVisible, &item.Visible, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt,
			&item.ObjectType.ID, &item.ObjectType.Slug, &item.ObjectType.Name, &item.ObjectType.Visible, &item.ObjectType.SortOrder, &item.ObjectType.CreatedAt, &item.ObjectType.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *LandingPageRepository) Create(ctx context.Context, item model.LandingPage) (model.LandingPage, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO landing_pages (object_type_id, slug, h1, meta_title, meta_description, faq_title, faq_description, faq_visible, visible, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`,
		item.ObjectTypeID,
		item.Slug,
		item.H1,
		item.MetaTitle,
		item.MetaDescription,
		item.FAQTitle,
		item.FAQDescription,
		item.FAQVisible,
		item.Visible,
		item.SortOrder,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *LandingPageRepository) Update(ctx context.Context, item model.LandingPage) (model.LandingPage, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE landing_pages
		SET object_type_id = $1, slug = $2, h1 = $3, meta_title = $4, meta_description = $5,
		    faq_title = $6, faq_description = $7, faq_visible = $8, visible = $9, sort_order = $10, updated_at = now()
		WHERE id = $11
		RETURNING updated_at
	`,
		item.ObjectTypeID,
		item.Slug,
		item.H1,
		item.MetaTitle,
		item.MetaDescription,
		item.FAQTitle,
		item.FAQDescription,
		item.FAQVisible,
		item.Visible,
		item.SortOrder,
		item.ID,
	).Scan(&item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *LandingPageRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM landing_pages WHERE id = $1`, id)
	return checkDeleted(tag, err)
}
