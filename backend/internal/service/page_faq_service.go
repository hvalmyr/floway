package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

// validFAQPages intentionally excludes "home" — the homepage keeps its own
// flat, unscoped FAQ (model.FAQItem/FAQService), predating any page concept.
var validFAQPages = map[string]bool{
	"masterclasses":    true,
	"gift_certificate": true,
}

func validateFAQPage(page string) error {
	if !validFAQPages[page] {
		return errors.Join(ErrValidation, errors.New("page must be one of: masterclasses, gift_certificate"))
	}
	return nil
}

type PageFAQRepository interface {
	GetSettings(ctx context.Context, page string) (model.PageFAQSettings, error)
	UpdateSettings(ctx context.Context, item model.PageFAQSettings) (model.PageFAQSettings, error)
	ListByPage(ctx context.Context, page string) ([]model.PageFAQItem, error)
	CreateItem(ctx context.Context, item model.PageFAQItem) (model.PageFAQItem, error)
	UpdateItem(ctx context.Context, item model.PageFAQItem) (model.PageFAQItem, error)
	DeleteItem(ctx context.Context, page string, id int64) error
}

// PageFAQService manages the FAQ block for a fixed set of static pages
// (masterclasses, gift certificate) — same shape as a course's own FAQ block
// (CourseFAQService), but keyed by a page identifier instead of a course id,
// since these pages aren't rows in a table with their own CRUD.
type PageFAQService struct {
	repo PageFAQRepository
}

func NewPageFAQService(repo PageFAQRepository) *PageFAQService {
	return &PageFAQService{repo: repo}
}

func (s *PageFAQService) Get(ctx context.Context, page string) (model.PageFAQ, error) {
	if err := validateFAQPage(page); err != nil {
		return model.PageFAQ{}, err
	}
	settings, err := s.repo.GetSettings(ctx, page)
	if err != nil {
		return model.PageFAQ{}, err
	}
	items, err := s.repo.ListByPage(ctx, page)
	if err != nil {
		return model.PageFAQ{}, err
	}
	return model.PageFAQ{PageFAQSettings: settings, Items: items}, nil
}

func (s *PageFAQService) UpdateSettings(ctx context.Context, item model.PageFAQSettings) (model.PageFAQSettings, error) {
	if err := validateFAQPage(item.Page); err != nil {
		return model.PageFAQSettings{}, err
	}
	item.Title = strings.TrimSpace(item.Title)
	item.Description = strings.TrimSpace(item.Description)
	return s.repo.UpdateSettings(ctx, item)
}

func (s *PageFAQService) CreateItem(ctx context.Context, item model.PageFAQItem) (model.PageFAQItem, error) {
	if err := validateFAQPage(item.Page); err != nil {
		return model.PageFAQItem{}, err
	}
	item.Question = strings.TrimSpace(item.Question)
	item.Answer = strings.TrimSpace(item.Answer)
	if item.Question == "" || item.Answer == "" {
		return model.PageFAQItem{}, errors.Join(ErrValidation, errors.New("question and answer are required"))
	}
	return s.repo.CreateItem(ctx, item)
}

func (s *PageFAQService) UpdateItem(ctx context.Context, item model.PageFAQItem) (model.PageFAQItem, error) {
	if item.ID == 0 {
		return model.PageFAQItem{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if err := validateFAQPage(item.Page); err != nil {
		return model.PageFAQItem{}, err
	}
	item.Question = strings.TrimSpace(item.Question)
	item.Answer = strings.TrimSpace(item.Answer)
	if item.Question == "" || item.Answer == "" {
		return model.PageFAQItem{}, errors.Join(ErrValidation, errors.New("question and answer are required"))
	}
	return s.repo.UpdateItem(ctx, item)
}

func (s *PageFAQService) DeleteItem(ctx context.Context, page string, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	if err := validateFAQPage(page); err != nil {
		return err
	}
	return s.repo.DeleteItem(ctx, page, id)
}
