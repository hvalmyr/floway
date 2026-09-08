package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type NotificationEmailRepository interface {
	List(ctx context.Context) ([]model.NotificationEmail, error)
	Create(ctx context.Context, item model.NotificationEmail) (model.NotificationEmail, error)
	Delete(ctx context.Context, id int64) error
}

// NotificationEmailService manages the recipient list for new-lead email
// notifications (see internal/notify.EmailNotifier, which reads this list
// at send time rather than once at startup).
type NotificationEmailService struct {
	repo NotificationEmailRepository
}

func NewNotificationEmailService(repo NotificationEmailRepository) *NotificationEmailService {
	return &NotificationEmailService{repo: repo}
}

func (s *NotificationEmailService) List(ctx context.Context) ([]model.NotificationEmail, error) {
	return s.repo.List(ctx)
}

func (s *NotificationEmailService) Create(ctx context.Context, item model.NotificationEmail) (model.NotificationEmail, error) {
	item.Email = strings.TrimSpace(item.Email)
	if item.Email == "" || !strings.Contains(item.Email, "@") {
		return model.NotificationEmail{}, errors.Join(ErrValidation, errors.New("a valid email is required"))
	}

	return s.repo.Create(ctx, item)
}

func (s *NotificationEmailService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
