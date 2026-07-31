package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

var validContactMethods = map[model.ContactMethod]struct{}{
	model.ContactMethodCall:     {},
	model.ContactMethodTelegram: {},
	model.ContactMethodWhatsapp: {},
	model.ContactMethodMax:      {},
}

var validLeadSources = map[model.LeadSource]struct{}{
	model.LeadSourceReferral: {},
	model.LeadSourceAds:      {},
	model.LeadSourceInternet: {},
	model.LeadSourceSocial:   {},
	model.LeadSourceMaps:     {},
}

var validLeadRequestTypes = map[model.LeadRequestType]struct{}{
	model.LeadRequestTypeCourse:      {},
	model.LeadRequestTypeMasterclass: {},
	model.LeadRequestTypeTrialLesson: {},
}

var validLeadStatuses = map[model.LeadStatus]struct{}{
	model.LeadStatusNew:        {},
	model.LeadStatusInProgress: {},
	model.LeadStatusClosed:     {},
}

type LeadRepository interface {
	List(ctx context.Context) ([]model.Lead, error)
	Create(ctx context.Context, item model.Lead) (model.Lead, error)
	UpdateStatus(ctx context.Context, id int64, status model.LeadStatus) (model.Lead, error)
	Delete(ctx context.Context, id int64) error
}

type LeadService struct {
	repo LeadRepository
}

func NewLeadService(repo LeadRepository) *LeadService {
	return &LeadService{repo: repo}
}

func (s *LeadService) List(ctx context.Context) ([]model.Lead, error) {
	return s.repo.List(ctx)
}

func (s *LeadService) Create(ctx context.Context, item model.Lead) (model.Lead, error) {
	item.Name = strings.TrimSpace(item.Name)
	item.Phone = strings.TrimSpace(item.Phone)
	item.Email = strings.TrimSpace(item.Email)

	if item.Name == "" || item.Phone == "" {
		return model.Lead{}, errors.Join(ErrValidation, errors.New("name and phone are required"))
	}

	if _, ok := validContactMethods[item.ContactMethod]; !ok {
		return model.Lead{}, errors.Join(ErrValidation, errors.New("invalid contact method"))
	}
	if _, ok := validLeadSources[item.Source]; !ok {
		return model.Lead{}, errors.Join(ErrValidation, errors.New("invalid source"))
	}
	if _, ok := validLeadRequestTypes[item.RequestType]; !ok {
		return model.Lead{}, errors.Join(ErrValidation, errors.New("invalid request type"))
	}

	// This is a public endpoint: the client never gets to set the status.
	item.Status = model.LeadStatusNew

	return s.repo.Create(ctx, item)
}

func (s *LeadService) UpdateStatus(ctx context.Context, id int64, status model.LeadStatus) (model.Lead, error) {
	if id == 0 {
		return model.Lead{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if _, ok := validLeadStatuses[status]; !ok {
		return model.Lead{}, errors.Join(ErrValidation, errors.New("invalid status"))
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *LeadService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
