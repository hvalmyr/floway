package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

var ErrValidation = errors.New("validation error")

type FAQRepository interface {
	List(ctx context.Context) ([]model.FAQItem, error)
	Create(ctx context.Context, item model.FAQItem) (model.FAQItem, error)
	Update(ctx context.Context, item model.FAQItem) (model.FAQItem, error)
	Delete(ctx context.Context, id int64) error
}

type FAQService struct {
	repo FAQRepository
}

func NewFAQService(repo FAQRepository) *FAQService {
	return &FAQService{repo: repo}
}

func (s *FAQService) List(ctx context.Context) ([]model.FAQItem, error) {
	return s.repo.List(ctx)
}

func (s *FAQService) Create(ctx context.Context, question, answer string, sortOrder int) (model.FAQItem, error) {
	question = strings.TrimSpace(question)
	answer = strings.TrimSpace(answer)
	if question == "" || answer == "" {
		return model.FAQItem{}, errors.Join(ErrValidation, errors.New("question and answer are required"))
	}

	return s.repo.Create(ctx, model.FAQItem{
		Question:  question,
		Answer:    answer,
		SortOrder: sortOrder,
	})
}

func (s *FAQService) Update(ctx context.Context, item model.FAQItem) (model.FAQItem, error) {
	item.Question = strings.TrimSpace(item.Question)
	item.Answer = strings.TrimSpace(item.Answer)
	if item.ID == 0 {
		return model.FAQItem{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if item.Question == "" || item.Answer == "" {
		return model.FAQItem{}, errors.Join(ErrValidation, errors.New("question and answer are required"))
	}

	return s.repo.Update(ctx, item)
}

func (s *FAQService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
