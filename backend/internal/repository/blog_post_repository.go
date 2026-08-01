package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type BlogPostRepository struct {
	db *pgxpool.Pool
}

func NewBlogPostRepository(db *pgxpool.Pool) *BlogPostRepository {
	return &BlogPostRepository{db: db}
}

// blogPostColumns/scanBlogPost keep the SELECT list and its positional Scan
// in exactly one place — List/ListPublished/FindPublishedBySlug all select
// every column, and a column added to one but not updated in the other used
// to be a real risk with no test to catch a scan-order drift (architecture
// review finding #12).
const blogPostColumns = "id, slug, title, cover_image, category, tags, author, published_at, content, status, created_at, updated_at"

func scanBlogPost(row pgx.Row) (model.BlogPost, error) {
	var item model.BlogPost
	err := row.Scan(&item.ID, &item.Slug, &item.Title, &item.CoverImage, &item.Category, &item.Tags, &item.Author, &item.PublishedAt, &item.Content, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *BlogPostRepository) List(ctx context.Context) ([]model.BlogPost, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+blogPostColumns+`
		FROM blog_posts
		ORDER BY published_at DESC NULLS LAST, id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.BlogPost{}
	for rows.Next() {
		item, err := scanBlogPost(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *BlogPostRepository) ListPublished(ctx context.Context) ([]model.BlogPost, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+blogPostColumns+`
		FROM blog_posts
		WHERE status = $1
		ORDER BY published_at DESC NULLS LAST, id DESC
	`, model.BlogPostStatusPublished)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.BlogPost{}
	for rows.Next() {
		item, err := scanBlogPost(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *BlogPostRepository) FindPublishedBySlug(ctx context.Context, slug string) (model.BlogPost, error) {
	row := r.db.QueryRow(ctx, `
		SELECT `+blogPostColumns+`
		FROM blog_posts
		WHERE slug = $1 AND status = $2
	`, slug, model.BlogPostStatusPublished)
	item, err := scanBlogPost(row)
	return item, translateNotFound(err)
}

func (r *BlogPostRepository) Create(ctx context.Context, item model.BlogPost) (model.BlogPost, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO blog_posts (slug, title, cover_image, category, tags, author, published_at, content, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`, item.Slug, item.Title, item.CoverImage, item.Category, item.Tags, item.Author, item.PublishedAt, item.Content, item.Status).
		Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *BlogPostRepository) Update(ctx context.Context, item model.BlogPost) (model.BlogPost, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE blog_posts
		SET slug = $1, title = $2, cover_image = $3, category = $4, tags = $5,
		    author = $6, published_at = $7, content = $8, status = $9, updated_at = now()
		WHERE id = $10
		RETURNING updated_at
	`, item.Slug, item.Title, item.CoverImage, item.Category, item.Tags, item.Author, item.PublishedAt, item.Content, item.Status, item.ID).
		Scan(&item.UpdatedAt)
	return item, translateNotFound(err)
}

func (r *BlogPostRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM blog_posts WHERE id = $1`, id)
	return checkDeleted(tag, err)
}
