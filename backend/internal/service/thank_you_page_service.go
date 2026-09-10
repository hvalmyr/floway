package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

// validThankYouPageVariants mirrors model.LeadRequestType's three values —
// the thank-you page a lead is sent to after ApplyForm.vue submits is keyed
// by the same "context" the lead itself records.
var validThankYouPageVariants = map[string]bool{
	"course":       true,
	"masterclass":  true,
	"trial_lesson": true,
}

func validateThankYouPageVariant(variant string) error {
	if !validThankYouPageVariants[variant] {
		return errors.Join(ErrValidation, errors.New("variant must be one of: course, masterclass, trial_lesson"))
	}
	return nil
}

type ThankYouPageRepository interface {
	GetSettings(ctx context.Context, variant string) (model.ThankYouPage, error)
	UpdateSettings(ctx context.Context, item model.ThankYouPage) (model.ThankYouPage, error)
	ListPhotos(ctx context.Context, variant string) ([]model.ThankYouPagePhoto, error)
	CreatePhoto(ctx context.Context, item model.ThankYouPagePhoto) (model.ThankYouPagePhoto, error)
	UpdatePhoto(ctx context.Context, item model.ThankYouPagePhoto) (model.ThankYouPagePhoto, error)
	DeletePhoto(ctx context.Context, variant string, id int64) error
	ListFAQItems(ctx context.Context, variant string) ([]model.ThankYouPageFAQItem, error)
	CreateFAQItem(ctx context.Context, item model.ThankYouPageFAQItem) (model.ThankYouPageFAQItem, error)
	UpdateFAQItem(ctx context.Context, item model.ThankYouPageFAQItem) (model.ThankYouPageFAQItem, error)
	DeleteFAQItem(ctx context.Context, variant string, id int64) error
}

// ThankYouPageService manages the content shown after a lead form
// submission — one settings row plus an optional photo carousel and
// mini-FAQ per LeadRequestType variant.
type ThankYouPageService struct {
	repo ThankYouPageRepository
}

func NewThankYouPageService(repo ThankYouPageRepository) *ThankYouPageService {
	return &ThankYouPageService{repo: repo}
}

func (s *ThankYouPageService) Get(ctx context.Context, variant string) (model.ThankYouPageFull, error) {
	if err := validateThankYouPageVariant(variant); err != nil {
		return model.ThankYouPageFull{}, err
	}
	settings, err := s.repo.GetSettings(ctx, variant)
	if err != nil {
		return model.ThankYouPageFull{}, err
	}
	photos, err := s.repo.ListPhotos(ctx, variant)
	if err != nil {
		return model.ThankYouPageFull{}, err
	}
	faqItems, err := s.repo.ListFAQItems(ctx, variant)
	if err != nil {
		return model.ThankYouPageFull{}, err
	}
	return model.ThankYouPageFull{ThankYouPage: settings, Photos: photos, FAQItems: faqItems}, nil
}

func (s *ThankYouPageService) UpdateSettings(ctx context.Context, item model.ThankYouPage) (model.ThankYouPage, error) {
	if err := validateThankYouPageVariant(item.Variant); err != nil {
		return model.ThankYouPage{}, err
	}
	item.Title = strings.TrimSpace(item.Title)
	item.Subtitle = strings.TrimSpace(item.Subtitle)
	item.Description = strings.TrimSpace(item.Description)
	item.HeroImage = strings.TrimSpace(item.HeroImage)
	item.BlogLinkText = strings.TrimSpace(item.BlogLinkText)
	item.BlogLinkURL = strings.TrimSpace(item.BlogLinkURL)
	item.CommunityText = strings.TrimSpace(item.CommunityText)
	item.CommunityURL = strings.TrimSpace(item.CommunityURL)
	return s.repo.UpdateSettings(ctx, item)
}

func (s *ThankYouPageService) CreatePhoto(ctx context.Context, item model.ThankYouPagePhoto) (model.ThankYouPagePhoto, error) {
	if err := validateThankYouPageVariant(item.Variant); err != nil {
		return model.ThankYouPagePhoto{}, err
	}
	if item.Image == "" {
		return model.ThankYouPagePhoto{}, errors.Join(ErrValidation, errors.New("image is required"))
	}
	return s.repo.CreatePhoto(ctx, item)
}

func (s *ThankYouPageService) UpdatePhoto(ctx context.Context, item model.ThankYouPagePhoto) (model.ThankYouPagePhoto, error) {
	if item.ID == 0 {
		return model.ThankYouPagePhoto{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if err := validateThankYouPageVariant(item.Variant); err != nil {
		return model.ThankYouPagePhoto{}, err
	}
	if item.Image == "" {
		return model.ThankYouPagePhoto{}, errors.Join(ErrValidation, errors.New("image is required"))
	}
	return s.repo.UpdatePhoto(ctx, item)
}

func (s *ThankYouPageService) DeletePhoto(ctx context.Context, variant string, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	if err := validateThankYouPageVariant(variant); err != nil {
		return err
	}
	return s.repo.DeletePhoto(ctx, variant, id)
}

func (s *ThankYouPageService) CreateFAQItem(ctx context.Context, item model.ThankYouPageFAQItem) (model.ThankYouPageFAQItem, error) {
	if err := validateThankYouPageVariant(item.Variant); err != nil {
		return model.ThankYouPageFAQItem{}, err
	}
	item.Question = strings.TrimSpace(item.Question)
	item.Answer = strings.TrimSpace(item.Answer)
	if item.Question == "" || item.Answer == "" {
		return model.ThankYouPageFAQItem{}, errors.Join(ErrValidation, errors.New("question and answer are required"))
	}
	return s.repo.CreateFAQItem(ctx, item)
}

func (s *ThankYouPageService) UpdateFAQItem(ctx context.Context, item model.ThankYouPageFAQItem) (model.ThankYouPageFAQItem, error) {
	if item.ID == 0 {
		return model.ThankYouPageFAQItem{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if err := validateThankYouPageVariant(item.Variant); err != nil {
		return model.ThankYouPageFAQItem{}, err
	}
	item.Question = strings.TrimSpace(item.Question)
	item.Answer = strings.TrimSpace(item.Answer)
	if item.Question == "" || item.Answer == "" {
		return model.ThankYouPageFAQItem{}, errors.Join(ErrValidation, errors.New("question and answer are required"))
	}
	return s.repo.UpdateFAQItem(ctx, item)
}

func (s *ThankYouPageService) DeleteFAQItem(ctx context.Context, variant string, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	if err := validateThankYouPageVariant(variant); err != nil {
		return err
	}
	return s.repo.DeleteFAQItem(ctx, variant, id)
}
