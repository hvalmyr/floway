package service

import (
	"context"
	"errors"
	"log/slog"
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

// LeadNotifier tells someone a new lead came in — email, Telegram, or both
// (see internal/notify.MultiNotifier). Optional: a nil notifier just means
// no channel is configured, not an error. programName is the resolved
// course/masterclass title (see resolveProgramName below), passed in
// separately since notify has no repository dependency of its own.
type LeadNotifier interface {
	NotifyNewLead(ctx context.Context, lead model.Lead, programName string) error
}

// CourseLookup and MasterclassLookup resolve a lead's RelatedSlug to a
// human title for the notification (name/title, not the raw slug). Both
// optional — a nil lookup just means the notification falls back to the
// generic request-type label ("Курс"/"Мастер-класс") instead of a title.
type CourseLookup interface {
	FindBySlug(ctx context.Context, slug string) (model.Course, error)
}

type MasterclassLookup interface {
	FindBySlug(ctx context.Context, slug string) (model.Masterclass, error)
}

type LeadService struct {
	repo          LeadRepository
	notifier      LeadNotifier
	courses       CourseLookup
	masterclasses MasterclassLookup
}

func NewLeadService(repo LeadRepository, notifier LeadNotifier, courses CourseLookup, masterclasses MasterclassLookup) *LeadService {
	return &LeadService{repo: repo, notifier: notifier, courses: courses, masterclasses: masterclasses}
}

// resolveProgramName looks up the human title for item.RelatedSlug — empty
// for trial-lesson leads (no related entity), and best-effort otherwise: a
// missing lookup, an empty slug, or a "not found" (deleted/renamed course)
// all just fall back to an empty string rather than failing the lead.
func (s *LeadService) resolveProgramName(ctx context.Context, item model.Lead) string {
	if item.RelatedSlug == "" {
		return ""
	}
	switch item.RequestType {
	case model.LeadRequestTypeCourse:
		if s.courses == nil {
			return ""
		}
		course, err := s.courses.FindBySlug(ctx, item.RelatedSlug)
		if err != nil {
			return ""
		}
		return course.Name
	case model.LeadRequestTypeMasterclass:
		if s.masterclasses == nil {
			return ""
		}
		masterclass, err := s.masterclasses.FindBySlug(ctx, item.RelatedSlug)
		if err != nil {
			return ""
		}
		return masterclass.Title
	default:
		return ""
	}
}

func (s *LeadService) List(ctx context.Context) ([]model.Lead, error) {
	return s.repo.List(ctx)
}

func (s *LeadService) Create(ctx context.Context, item model.Lead) (model.Lead, error) {
	item.Name = strings.TrimSpace(item.Name)
	item.Phone = strings.TrimSpace(item.Phone)
	item.Email = strings.TrimSpace(item.Email)
	item.RelatedSlug = strings.TrimSpace(item.RelatedSlug)

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

	created, err := s.repo.Create(ctx, item)
	if err != nil {
		return model.Lead{}, err
	}

	if s.notifier != nil {
		// Best-effort: the lead is already saved and will show up in the
		// admin panel regardless — a notification hiccup (SMTP relay down,
		// Telegram API blip) must not turn a successful form submission
		// into a 500 for the visitor. Logged so it's not silently invisible.
		programName := s.resolveProgramName(ctx, created)
		if err := s.notifier.NotifyNewLead(ctx, created, programName); err != nil {
			slog.Default().Error("lead notification failed", "leadId", created.ID, "error", err)
		}
	}

	return created, nil
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
