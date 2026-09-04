package service

import (
	"context"
	"errors"
	"log/slog"
	"regexp"
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
	model.LeadStatusNew:           {},
	model.LeadStatusInProgress:    {},
	model.LeadStatusWaitingClient: {},
	model.LeadStatusBooked:        {},
	model.LeadStatusPostponed:     {},
	model.LeadStatusClosedWon:     {},
	model.LeadStatusClosedLost:    {},
}

type LeadRepository interface {
	ListWithClient(ctx context.Context) ([]model.LeadListItem, error)
	Create(ctx context.Context, item model.Lead) (model.Lead, error)
	UpdateStatus(ctx context.Context, id int64, status model.LeadStatus) (model.Lead, error)
	DismissReview(ctx context.Context, id int64) (model.Lead, error)
	CountByStatus(ctx context.Context, statuses ...model.LeadStatus) (map[model.LeadStatus]int, error)
	Delete(ctx context.Context, id int64) error
}

// ClientRepository is the subset of client persistence LeadService needs to
// dedup an incoming submission against an existing customer profile.
type ClientRepository interface {
	FindByPhoneOrEmail(ctx context.Context, phoneNormalized, email string) (model.Client, error)
	Create(ctx context.Context, item model.Client) (model.Client, error)
	RefreshContactInfo(ctx context.Context, id int64, item model.Client) (model.Client, error)
}

// normalizePhone keeps digits only, then folds the RU/CIS "8-first" dialing
// convention onto "+7" so "89991234567" and "+7 999 123-45-67" match as the
// same number — a common real-world input variance for this audience.
var nonDigit = regexp.MustCompile(`\D`)

func normalizePhone(phone string) string {
	digits := nonDigit.ReplaceAllString(phone, "")
	if len(digits) == 11 && digits[0] == '8' {
		digits = "7" + digits[1:]
	}
	return digits
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
	clients       ClientRepository
	notifier      LeadNotifier
	courses       CourseLookup
	masterclasses MasterclassLookup
}

func NewLeadService(repo LeadRepository, clients ClientRepository, notifier LeadNotifier, courses CourseLookup, masterclasses MasterclassLookup) *LeadService {
	return &LeadService{repo: repo, clients: clients, notifier: notifier, courses: courses, masterclasses: masterclasses}
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

func (s *LeadService) List(ctx context.Context) ([]model.LeadListItem, error) {
	return s.repo.ListWithClient(ctx)
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

	// Dedup against an existing customer by phone or email rather than
	// creating a disconnected client record for every submission. A match
	// gets its profile refreshed to the newest contact info; the lead row
	// itself keeps its own untouched snapshot of what was submitted here.
	normalizedPhone := normalizePhone(item.Phone)
	client, err := s.clients.FindByPhoneOrEmail(ctx, normalizedPhone, item.Email)
	switch {
	case err == nil:
		client, err = s.clients.RefreshContactInfo(ctx, client.ID, model.Client{
			Name: item.Name, Phone: item.Phone, PhoneNormalized: normalizedPhone, Email: item.Email,
		})
		if err != nil {
			return model.Lead{}, err
		}
	case errors.Is(err, ErrNotFound):
		client, err = s.clients.Create(ctx, model.Client{
			Name: item.Name, Phone: item.Phone, PhoneNormalized: normalizedPhone, Email: item.Email,
		})
		if err != nil {
			return model.Lead{}, err
		}
	default:
		return model.Lead{}, err
	}
	item.ClientID = client.ID

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

// DismissReview clears NeedsStatusReview without changing status — for a
// lead auto-migrated to closed_won (see migration 00031) where that
// default was already correct.
func (s *LeadService) DismissReview(ctx context.Context, id int64) (model.Lead, error) {
	if id == 0 {
		return model.Lead{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.DismissReview(ctx, id)
}

// ConversionRate reports closed_won vs. closed_lost counts and their
// ratio. Rate is nil rather than dividing by zero when neither status has
// any leads yet.
func (s *LeadService) ConversionRate(ctx context.Context) (won, lost int, rate *float64, err error) {
	counts, err := s.repo.CountByStatus(ctx, model.LeadStatusClosedWon, model.LeadStatusClosedLost)
	if err != nil {
		return 0, 0, nil, err
	}
	won = counts[model.LeadStatusClosedWon]
	lost = counts[model.LeadStatusClosedLost]
	if won+lost == 0 {
		return won, lost, nil, nil
	}
	computed := float64(won) / float64(won+lost)
	return won, lost, &computed, nil
}

func (s *LeadService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
