package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"floway-backend/internal/model"
)

type ClientLookupRepository interface {
	FindByID(ctx context.Context, id int64) (model.Client, error)
}

type ClientLeadRepository interface {
	ListByClientID(ctx context.Context, clientID int64) ([]model.Lead, error)
}

type TagTypeRepository interface {
	FindOrCreateByName(ctx context.Context, name string) (model.Tag, error)
	SetForClient(ctx context.Context, clientID int64, tagIDs []int64) error
	ListForClient(ctx context.Context, clientID int64) ([]model.Tag, error)
}

type ClientCommentRepository interface {
	ListForClient(ctx context.Context, clientID int64) ([]model.ClientComment, error)
	Create(ctx context.Context, item model.ClientComment) (model.ClientComment, error)
}

type ReminderRepository interface {
	ListForClient(ctx context.Context, clientID int64) ([]model.Reminder, error)
	Create(ctx context.Context, item model.Reminder) (model.Reminder, error)
	Complete(ctx context.Context, id int64) error
}

type ClientService struct {
	clients     ClientLookupRepository
	leads       ClientLeadRepository
	productTags TagTypeRepository
	typeTags    TagTypeRepository
	comments    ClientCommentRepository
	reminders   ReminderRepository
}

func NewClientService(
	clients ClientLookupRepository,
	leads ClientLeadRepository,
	productTags TagTypeRepository,
	typeTags TagTypeRepository,
	comments ClientCommentRepository,
	reminders ReminderRepository,
) *ClientService {
	return &ClientService{
		clients:     clients,
		leads:       leads,
		productTags: productTags,
		typeTags:    typeTags,
		comments:    comments,
		reminders:   reminders,
	}
}

func (s *ClientService) GetDetail(ctx context.Context, id int64) (model.ClientDetail, error) {
	if id == 0 {
		return model.ClientDetail{}, errors.Join(ErrValidation, errors.New("id is required"))
	}

	client, err := s.clients.FindByID(ctx, id)
	if err != nil {
		return model.ClientDetail{}, err
	}
	requests, err := s.leads.ListByClientID(ctx, id)
	if err != nil {
		return model.ClientDetail{}, err
	}
	comments, err := s.comments.ListForClient(ctx, id)
	if err != nil {
		return model.ClientDetail{}, err
	}
	productTags, err := s.productTags.ListForClient(ctx, id)
	if err != nil {
		return model.ClientDetail{}, err
	}
	clientTypeTags, err := s.typeTags.ListForClient(ctx, id)
	if err != nil {
		return model.ClientDetail{}, err
	}
	reminders, err := s.reminders.ListForClient(ctx, id)
	if err != nil {
		return model.ClientDetail{}, err
	}

	return model.ClientDetail{
		Client:         client,
		Requests:       requests,
		Comments:       comments,
		ProductTags:    productTags,
		ClientTypeTags: clientTypeTags,
		Reminders:      reminders,
	}, nil
}

// setTags is shared by SetProductTags/SetClientTypeTags: find-or-create
// each submitted name, then replace the client's full tag set — the
// backend half of "autocomplete, create a new tag on the fly."
func setTags(ctx context.Context, repo TagTypeRepository, clientID int64, names []string) ([]model.Tag, error) {
	if clientID == 0 {
		return nil, errors.Join(ErrValidation, errors.New("client id is required"))
	}

	ids := make([]int64, 0, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		tag, err := repo.FindOrCreateByName(ctx, name)
		if err != nil {
			return nil, err
		}
		ids = append(ids, tag.ID)
	}

	if err := repo.SetForClient(ctx, clientID, ids); err != nil {
		return nil, err
	}
	return repo.ListForClient(ctx, clientID)
}

func (s *ClientService) SetProductTags(ctx context.Context, clientID int64, names []string) ([]model.Tag, error) {
	return setTags(ctx, s.productTags, clientID, names)
}

func (s *ClientService) SetClientTypeTags(ctx context.Context, clientID int64, names []string) ([]model.Tag, error) {
	return setTags(ctx, s.typeTags, clientID, names)
}

func (s *ClientService) AddComment(ctx context.Context, clientID int64, text string) (model.ClientComment, error) {
	if clientID == 0 {
		return model.ClientComment{}, errors.Join(ErrValidation, errors.New("client id is required"))
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return model.ClientComment{}, errors.Join(ErrValidation, errors.New("comment text is required"))
	}
	return s.comments.Create(ctx, model.ClientComment{ClientID: clientID, Text: text})
}

func (s *ClientService) AddReminder(ctx context.Context, clientID int64, days int, note string) (model.Reminder, error) {
	if clientID == 0 {
		return model.Reminder{}, errors.Join(ErrValidation, errors.New("client id is required"))
	}
	if days < 1 {
		return model.Reminder{}, errors.Join(ErrValidation, errors.New("days must be at least 1"))
	}
	remindAt := time.Now().UTC().AddDate(0, 0, days)
	return s.reminders.Create(ctx, model.Reminder{
		ClientID: clientID,
		RemindAt: remindAt,
		Note:     strings.TrimSpace(note),
	})
}

func (s *ClientService) CompleteReminder(ctx context.Context, reminderID int64) error {
	if reminderID == 0 {
		return errors.Join(ErrValidation, errors.New("reminder id is required"))
	}
	return s.reminders.Complete(ctx, reminderID)
}
