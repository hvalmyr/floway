package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"

	"floway-backend/internal/model"
)

var validBlogPostStatuses = map[model.BlogPostStatus]bool{
	model.BlogPostStatusDraft:     true,
	model.BlogPostStatusPublished: true,
}

type BlogPostRepository interface {
	List(ctx context.Context) ([]model.BlogPost, error)
	ListPublished(ctx context.Context) ([]model.BlogPost, error)
	FindPublishedBySlug(ctx context.Context, slug string) (model.BlogPost, error)
	Create(ctx context.Context, item model.BlogPost) (model.BlogPost, error)
	Update(ctx context.Context, item model.BlogPost) (model.BlogPost, error)
	Delete(ctx context.Context, id int64) error
}

type BlogPostService struct {
	repo BlogPostRepository
}

func NewBlogPostService(repo BlogPostRepository) *BlogPostService {
	return &BlogPostService{repo: repo}
}

func (s *BlogPostService) List(ctx context.Context) ([]model.BlogPost, error) {
	return s.repo.List(ctx)
}

func (s *BlogPostService) ListPublished(ctx context.Context) ([]model.BlogPost, error) {
	return s.repo.ListPublished(ctx)
}

func (s *BlogPostService) GetPublishedBySlug(ctx context.Context, slug string) (model.BlogPost, error) {
	item, err := s.repo.FindPublishedBySlug(ctx, slug)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.BlogPost{}, ErrNotFound
	}
	return item, err
}

func (s *BlogPostService) validate(item *model.BlogPost) error {
	item.Slug = strings.TrimSpace(item.Slug)
	item.Title = strings.TrimSpace(item.Title)
	if item.Slug == "" || item.Title == "" {
		return errors.Join(ErrValidation, errors.New("slug and title are required"))
	}
	if item.Status == "" {
		item.Status = model.BlogPostStatusDraft
	}
	if !validBlogPostStatuses[item.Status] {
		return errors.Join(ErrValidation, errors.New("invalid status"))
	}
	// tags is NOT NULL DEFAULT '{}' in the DB, but a nil Go slice (an omitted
	// or explicit-null "tags" in the request JSON) is sent as SQL NULL, not
	// as "use the column default" — violates the constraint. Normalize here.
	if item.Tags == nil {
		item.Tags = []string{}
	}
	return nil
}

func (s *BlogPostService) Create(ctx context.Context, item model.BlogPost) (model.BlogPost, error) {
	if err := s.validate(&item); err != nil {
		return model.BlogPost{}, err
	}
	return s.repo.Create(ctx, item)
}

func (s *BlogPostService) Update(ctx context.Context, item model.BlogPost) (model.BlogPost, error) {
	if item.ID == 0 {
		return model.BlogPost{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if err := s.validate(&item); err != nil {
		return model.BlogPost{}, err
	}
	return s.repo.Update(ctx, item)
}

func (s *BlogPostService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
