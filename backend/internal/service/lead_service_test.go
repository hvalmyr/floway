package service_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type fakeLeadRepository struct {
	items  []model.Lead
	nextID int64
	err    error
}

func newFakeLeadRepository() *fakeLeadRepository {
	return &fakeLeadRepository{nextID: 1}
}

func (f *fakeLeadRepository) ListWithClient(ctx context.Context) ([]model.LeadListItem, error) {
	if f.err != nil {
		return nil, f.err
	}
	items := make([]model.LeadListItem, len(f.items))
	for i, item := range f.items {
		items[i] = model.LeadListItem{Lead: item}
	}
	return items, nil
}

func (f *fakeLeadRepository) Create(ctx context.Context, item model.Lead) (model.Lead, error) {
	if f.err != nil {
		return model.Lead{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeLeadRepository) UpdateStatus(ctx context.Context, id int64, status model.LeadStatus) (model.Lead, error) {
	if f.err != nil {
		return model.Lead{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == id {
			f.items[i].Status = status
			f.items[i].NeedsStatusReview = false
			return f.items[i], nil
		}
	}
	return model.Lead{}, errors.New("not found")
}

func (f *fakeLeadRepository) DismissReview(ctx context.Context, id int64) (model.Lead, error) {
	if f.err != nil {
		return model.Lead{}, f.err
	}
	for i, existing := range f.items {
		if existing.ID == id {
			f.items[i].NeedsStatusReview = false
			return f.items[i], nil
		}
	}
	return model.Lead{}, errors.New("not found")
}

func (f *fakeLeadRepository) CountByStatus(ctx context.Context, statuses ...model.LeadStatus) (map[model.LeadStatus]int, error) {
	if f.err != nil {
		return nil, f.err
	}
	wanted := make(map[model.LeadStatus]struct{}, len(statuses))
	for _, s := range statuses {
		wanted[s] = struct{}{}
	}
	counts := map[model.LeadStatus]int{}
	for _, item := range f.items {
		if _, ok := wanted[item.Status]; ok {
			counts[item.Status]++
		}
	}
	return counts, nil
}

func (f *fakeLeadRepository) Delete(ctx context.Context, id int64) error {
	if f.err != nil {
		return f.err
	}
	for i, existing := range f.items {
		if existing.ID == id {
			f.items = append(f.items[:i], f.items[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

// fakeClientRepository is an in-memory stand-in for ClientRepository,
// tracking call counts so dedup tests can assert whether Create or
// RefreshContactInfo actually ran.
type fakeClientRepository struct {
	byID         map[int64]model.Client
	nextID       int64
	createCalls  int
	refreshCalls int
	err          error
}

func newFakeClientRepository() *fakeClientRepository {
	return &fakeClientRepository{byID: map[int64]model.Client{}, nextID: 1}
}

func (f *fakeClientRepository) FindByPhoneOrEmail(ctx context.Context, phoneNormalized, email string) (model.Client, error) {
	if f.err != nil {
		return model.Client{}, f.err
	}
	for _, c := range f.byID {
		if phoneNormalized != "" && c.PhoneNormalized == phoneNormalized {
			return c, nil
		}
		if email != "" && c.Email != "" && strings.EqualFold(c.Email, email) {
			return c, nil
		}
	}
	return model.Client{}, service.ErrNotFound
}

func (f *fakeClientRepository) Create(ctx context.Context, item model.Client) (model.Client, error) {
	if f.err != nil {
		return model.Client{}, f.err
	}
	f.createCalls++
	item.ID = f.nextID
	f.nextID++
	f.byID[item.ID] = item
	return item, nil
}

func (f *fakeClientRepository) RefreshContactInfo(ctx context.Context, id int64, item model.Client) (model.Client, error) {
	if f.err != nil {
		return model.Client{}, f.err
	}
	f.refreshCalls++
	existing, ok := f.byID[id]
	if !ok {
		return model.Client{}, service.ErrNotFound
	}
	existing.Name = item.Name
	existing.Phone = item.Phone
	existing.PhoneNormalized = item.PhoneNormalized
	existing.Email = item.Email
	f.byID[id] = existing
	return existing, nil
}

type fakeLeadNotifier struct {
	notified     []model.Lead
	programNames []string
	err          error
}

func (f *fakeLeadNotifier) NotifyNewLead(ctx context.Context, lead model.Lead, programName string) error {
	f.notified = append(f.notified, lead)
	f.programNames = append(f.programNames, programName)
	return f.err
}

type fakeCourseLookup struct {
	bySlug map[string]model.Course
}

func (f *fakeCourseLookup) FindBySlug(ctx context.Context, slug string) (model.Course, error) {
	c, ok := f.bySlug[slug]
	if !ok {
		return model.Course{}, errors.New("not found")
	}
	return c, nil
}

type fakeMasterclassLookup struct {
	bySlug map[string]model.Masterclass
}

func (f *fakeMasterclassLookup) FindBySlug(ctx context.Context, slug string) (model.Masterclass, error) {
	mc, ok := f.bySlug[slug]
	if !ok {
		return model.Masterclass{}, errors.New("not found")
	}
	return mc, nil
}

func validLead() model.Lead {
	return model.Lead{
		Name:          "Иван Иванов",
		Phone:         "+79991234567",
		Email:         "ivan@example.com",
		ContactMethod: model.ContactMethodTelegram,
		Source:        model.LeadSourceAds,
		RequestType:   model.LeadRequestTypeCourse,
	}
}

func TestLeadService_Create(t *testing.T) {
	t.Run("creates a valid lead", func(t *testing.T) {
		repo := newFakeLeadRepository()
		svc := service.NewLeadService(repo, newFakeClientRepository(), nil, nil, nil)

		lead, err := svc.Create(context.Background(), validLead())

		require.NoError(t, err)
		assert.Equal(t, int64(1), lead.ID)
		assert.Equal(t, "Иван Иванов", lead.Name)
		assert.Equal(t, model.LeadStatusNew, lead.Status)
	})

	t.Run("trims name and phone", func(t *testing.T) {
		repo := newFakeLeadRepository()
		svc := service.NewLeadService(repo, newFakeClientRepository(), nil, nil, nil)

		item := validLead()
		item.Name = "  Иван Иванов  "
		item.Phone = "  +79991234567  "

		lead, err := svc.Create(context.Background(), item)

		require.NoError(t, err)
		assert.Equal(t, "Иван Иванов", lead.Name)
		assert.Equal(t, "+79991234567", lead.Phone)
	})

	t.Run("rejects an empty name", func(t *testing.T) {
		repo := newFakeLeadRepository()
		svc := service.NewLeadService(repo, newFakeClientRepository(), nil, nil, nil)

		item := validLead()
		item.Name = "   "

		_, err := svc.Create(context.Background(), item)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})

	t.Run("rejects an empty phone", func(t *testing.T) {
		repo := newFakeLeadRepository()
		svc := service.NewLeadService(repo, newFakeClientRepository(), nil, nil, nil)

		item := validLead()
		item.Phone = "   "

		_, err := svc.Create(context.Background(), item)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects an invalid contact method", func(t *testing.T) {
		repo := newFakeLeadRepository()
		svc := service.NewLeadService(repo, newFakeClientRepository(), nil, nil, nil)

		item := validLead()
		item.ContactMethod = "carrier_pigeon"

		_, err := svc.Create(context.Background(), item)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects an invalid source", func(t *testing.T) {
		repo := newFakeLeadRepository()
		svc := service.NewLeadService(repo, newFakeClientRepository(), nil, nil, nil)

		item := validLead()
		item.Source = "unknown"

		_, err := svc.Create(context.Background(), item)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects an invalid request type", func(t *testing.T) {
		repo := newFakeLeadRepository()
		svc := service.NewLeadService(repo, newFakeClientRepository(), nil, nil, nil)

		item := validLead()
		item.RequestType = "unknown"

		_, err := svc.Create(context.Background(), item)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("forces status to new regardless of input", func(t *testing.T) {
		repo := newFakeLeadRepository()
		svc := service.NewLeadService(repo, newFakeClientRepository(), nil, nil, nil)

		item := validLead()
		item.Status = model.LeadStatusClosedWon

		lead, err := svc.Create(context.Background(), item)

		require.NoError(t, err)
		assert.Equal(t, model.LeadStatusNew, lead.Status)
	})
}

func TestLeadService_Create_ClientDedup(t *testing.T) {
	t.Run("attaches to an existing client matched by phone, despite formatting differences", func(t *testing.T) {
		repo := newFakeLeadRepository()
		clients := newFakeClientRepository()
		existing, err := clients.Create(context.Background(), model.Client{Name: "Иван", PhoneNormalized: "79991234567"})
		require.NoError(t, err)
		clients.createCalls = 0 // reset past the seed insert above
		svc := service.NewLeadService(repo, clients, nil, nil, nil)

		item := validLead()
		item.Phone = "8 (999) 123-45-67"
		item.Email = "" // no email match possible, only phone

		lead, err := svc.Create(context.Background(), item)

		require.NoError(t, err)
		assert.Equal(t, existing.ID, lead.ClientID)
		assert.Zero(t, clients.createCalls, "must not create a second client for a phone that already matches")
	})

	t.Run("attaches to an existing client matched by email when the phone differs", func(t *testing.T) {
		repo := newFakeLeadRepository()
		clients := newFakeClientRepository()
		existing, err := clients.Create(context.Background(), model.Client{Name: "Иван", PhoneNormalized: "70000000000", Email: "ivan@example.com"})
		require.NoError(t, err)
		clients.createCalls = 0 // reset past the seed insert above
		svc := service.NewLeadService(repo, clients, nil, nil, nil)

		item := validLead()
		item.Phone = "+79991234567" // different phone
		item.Email = "IVAN@example.com"

		lead, err := svc.Create(context.Background(), item)

		require.NoError(t, err)
		assert.Equal(t, existing.ID, lead.ClientID)
		assert.Zero(t, clients.createCalls)
	})

	t.Run("creates a new client when no phone or email matches", func(t *testing.T) {
		repo := newFakeLeadRepository()
		clients := newFakeClientRepository()
		svc := service.NewLeadService(repo, clients, nil, nil, nil)

		lead, err := svc.Create(context.Background(), validLead())

		require.NoError(t, err)
		assert.Equal(t, 1, clients.createCalls)
		assert.NotZero(t, lead.ClientID)
	})

	t.Run("refreshes the matched client's contact info while the lead keeps its own submitted snapshot", func(t *testing.T) {
		repo := newFakeLeadRepository()
		clients := newFakeClientRepository()
		existing, err := clients.Create(context.Background(), model.Client{Name: "Old Name", PhoneNormalized: "79991234567", Email: "old@example.com"})
		require.NoError(t, err)
		svc := service.NewLeadService(repo, clients, nil, nil, nil)

		item := validLead()
		item.Name = "New Name"
		item.Phone = "+79991234567"
		item.Email = "new@example.com"

		lead, err := svc.Create(context.Background(), item)

		require.NoError(t, err)
		assert.Equal(t, 1, clients.refreshCalls)
		refreshed := clients.byID[existing.ID]
		assert.Equal(t, "New Name", refreshed.Name, "client profile should pick up the newest submitted name")
		assert.Equal(t, "new@example.com", refreshed.Email)
		assert.Equal(t, "New Name", lead.Name, "the lead's own snapshot is whatever was submitted this time")
	})
}

func TestLeadService_Create_Notifications(t *testing.T) {
	t.Run("notifies with the created lead on success", func(t *testing.T) {
		repo := newFakeLeadRepository()
		notifier := &fakeLeadNotifier{}
		svc := service.NewLeadService(repo, newFakeClientRepository(), notifier, nil, nil)

		lead, err := svc.Create(context.Background(), validLead())

		require.NoError(t, err)
		require.Len(t, notifier.notified, 1)
		assert.Equal(t, lead.ID, notifier.notified[0].ID, "must be notified with the ID assigned by the repository, not the pre-insert value")
	})

	t.Run("a notification failure does not fail lead creation", func(t *testing.T) {
		repo := newFakeLeadRepository()
		notifier := &fakeLeadNotifier{err: errors.New("smtp relay unreachable")}
		svc := service.NewLeadService(repo, newFakeClientRepository(), notifier, nil, nil)

		lead, err := svc.Create(context.Background(), validLead())

		require.NoError(t, err, "the lead is already saved — a notification hiccup must not turn a successful submission into an error")
		assert.NotZero(t, lead.ID)
	})

	t.Run("does not notify when validation rejects the lead", func(t *testing.T) {
		repo := newFakeLeadRepository()
		notifier := &fakeLeadNotifier{}
		svc := service.NewLeadService(repo, newFakeClientRepository(), notifier, nil, nil)

		item := validLead()
		item.Name = "   "
		_, err := svc.Create(context.Background(), item)

		require.Error(t, err)
		assert.Empty(t, notifier.notified, "nothing was actually created — there's nothing to notify about")
	})

	t.Run("does not notify when the repository fails", func(t *testing.T) {
		repo := newFakeLeadRepository()
		repo.err = errors.New("boom")
		notifier := &fakeLeadNotifier{}
		svc := service.NewLeadService(repo, newFakeClientRepository(), notifier, nil, nil)

		_, err := svc.Create(context.Background(), validLead())

		require.Error(t, err)
		assert.Empty(t, notifier.notified)
	})
}

func TestLeadService_Create_ResolvesProgramNameForNotification(t *testing.T) {
	t.Run("resolves a course title by slug", func(t *testing.T) {
		repo := newFakeLeadRepository()
		notifier := &fakeLeadNotifier{}
		courses := &fakeCourseLookup{bySlug: map[string]model.Course{
			"aktualnaya-floristika": {Name: "Актуальная флористика"},
		}}
		svc := service.NewLeadService(repo, newFakeClientRepository(), notifier, courses, nil)

		item := validLead()
		item.RequestType = model.LeadRequestTypeCourse
		item.RelatedSlug = "aktualnaya-floristika"
		_, err := svc.Create(context.Background(), item)

		require.NoError(t, err)
		require.Len(t, notifier.programNames, 1)
		assert.Equal(t, "Актуальная флористика", notifier.programNames[0])
	})

	t.Run("resolves a masterclass title by slug", func(t *testing.T) {
		repo := newFakeLeadRepository()
		notifier := &fakeLeadNotifier{}
		masterclasses := &fakeMasterclassLookup{bySlug: map[string]model.Masterclass{
			"novogodnyaya-floristika": {Title: "Новогодняя флористика"},
		}}
		svc := service.NewLeadService(repo, newFakeClientRepository(), notifier, nil, masterclasses)

		item := validLead()
		item.RequestType = model.LeadRequestTypeMasterclass
		item.RelatedSlug = "novogodnyaya-floristika"
		_, err := svc.Create(context.Background(), item)

		require.NoError(t, err)
		require.Len(t, notifier.programNames, 1)
		assert.Equal(t, "Новогодняя флористика", notifier.programNames[0])
	})

	t.Run("falls back to an empty program name when the slug isn't found", func(t *testing.T) {
		repo := newFakeLeadRepository()
		notifier := &fakeLeadNotifier{}
		courses := &fakeCourseLookup{bySlug: map[string]model.Course{}}
		svc := service.NewLeadService(repo, newFakeClientRepository(), notifier, courses, nil)

		item := validLead()
		item.RequestType = model.LeadRequestTypeCourse
		item.RelatedSlug = "deleted-course"
		_, err := svc.Create(context.Background(), item)

		require.NoError(t, err, "a lookup miss must not fail the lead")
		require.Len(t, notifier.programNames, 1)
		assert.Empty(t, notifier.programNames[0])
	})

	t.Run("does not look anything up for a trial lesson", func(t *testing.T) {
		repo := newFakeLeadRepository()
		notifier := &fakeLeadNotifier{}
		svc := service.NewLeadService(repo, newFakeClientRepository(), notifier, nil, nil)

		item := validLead()
		item.RequestType = model.LeadRequestTypeTrialLesson
		item.RelatedSlug = ""
		_, err := svc.Create(context.Background(), item)

		require.NoError(t, err)
		require.Len(t, notifier.programNames, 1)
		assert.Empty(t, notifier.programNames[0])
	})
}

func TestLeadService_UpdateStatus(t *testing.T) {
	repo := newFakeLeadRepository()
	svc := service.NewLeadService(repo, newFakeClientRepository(), nil, nil, nil)
	created, err := svc.Create(context.Background(), validLead())
	require.NoError(t, err)

	allStatuses := []model.LeadStatus{
		model.LeadStatusNew,
		model.LeadStatusInProgress,
		model.LeadStatusWaitingClient,
		model.LeadStatusBooked,
		model.LeadStatusPostponed,
		model.LeadStatusClosedWon,
		model.LeadStatusClosedLost,
	}
	for _, status := range allStatuses {
		t.Run("accepts "+string(status), func(t *testing.T) {
			lead, err := svc.UpdateStatus(context.Background(), created.ID, status)

			require.NoError(t, err)
			assert.Equal(t, status, lead.Status)
		})
	}

	t.Run("rejects the old pre-migration closed value", func(t *testing.T) {
		_, err := svc.UpdateStatus(context.Background(), created.ID, "closed")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects an invalid status", func(t *testing.T) {
		_, err := svc.UpdateStatus(context.Background(), created.ID, "unknown")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects a zero id", func(t *testing.T) {
		_, err := svc.UpdateStatus(context.Background(), 0, model.LeadStatusClosedWon)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("clears needs_status_review on an explicit status write", func(t *testing.T) {
		repo := newFakeLeadRepository()
		svc := service.NewLeadService(repo, newFakeClientRepository(), nil, nil, nil)
		created, err := svc.Create(context.Background(), validLead())
		require.NoError(t, err)
		repo.items[0].NeedsStatusReview = true

		lead, err := svc.UpdateStatus(context.Background(), created.ID, model.LeadStatusClosedLost)

		require.NoError(t, err)
		assert.False(t, lead.NeedsStatusReview)
	})
}

func TestLeadService_DismissReview(t *testing.T) {
	t.Run("clears the review flag without changing status", func(t *testing.T) {
		repo := newFakeLeadRepository()
		svc := service.NewLeadService(repo, newFakeClientRepository(), nil, nil, nil)
		created, err := svc.Create(context.Background(), validLead())
		require.NoError(t, err)
		repo.items[0].Status = model.LeadStatusClosedWon
		repo.items[0].NeedsStatusReview = true

		lead, err := svc.DismissReview(context.Background(), created.ID)

		require.NoError(t, err)
		assert.False(t, lead.NeedsStatusReview)
		assert.Equal(t, model.LeadStatusClosedWon, lead.Status)
	})

	t.Run("rejects a zero id", func(t *testing.T) {
		svc := service.NewLeadService(newFakeLeadRepository(), newFakeClientRepository(), nil, nil, nil)

		_, err := svc.DismissReview(context.Background(), 0)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestLeadService_ConversionRate(t *testing.T) {
	t.Run("computes the ratio of won to won+lost", func(t *testing.T) {
		repo := newFakeLeadRepository()
		svc := service.NewLeadService(repo, newFakeClientRepository(), nil, nil, nil)
		for _, status := range []model.LeadStatus{
			model.LeadStatusClosedWon, model.LeadStatusClosedWon, model.LeadStatusClosedWon,
			model.LeadStatusClosedLost,
			model.LeadStatusNew, // must not count toward either bucket
		} {
			created, err := svc.Create(context.Background(), validLead())
			require.NoError(t, err)
			_, err = svc.UpdateStatus(context.Background(), created.ID, status)
			require.NoError(t, err)
		}

		won, lost, rate, err := svc.ConversionRate(context.Background())

		require.NoError(t, err)
		assert.Equal(t, 3, won)
		assert.Equal(t, 1, lost)
		require.NotNil(t, rate)
		assert.InDelta(t, 0.75, *rate, 0.0001)
	})

	t.Run("returns a nil rate instead of dividing by zero", func(t *testing.T) {
		svc := service.NewLeadService(newFakeLeadRepository(), newFakeClientRepository(), nil, nil, nil)

		won, lost, rate, err := svc.ConversionRate(context.Background())

		require.NoError(t, err)
		assert.Zero(t, won)
		assert.Zero(t, lost)
		assert.Nil(t, rate)
	})
}

func TestLeadService_Delete(t *testing.T) {
	repo := newFakeLeadRepository()
	svc := service.NewLeadService(repo, newFakeClientRepository(), nil, nil, nil)
	created, err := svc.Create(context.Background(), validLead())
	require.NoError(t, err)

	t.Run("deletes an existing lead", func(t *testing.T) {
		require.NoError(t, svc.Delete(context.Background(), created.ID))
		items, err := svc.List(context.Background())
		require.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("rejects a zero id", func(t *testing.T) {
		err := svc.Delete(context.Background(), 0)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestLeadService_List_PropagatesRepositoryError(t *testing.T) {
	repo := newFakeLeadRepository()
	repo.err = errors.New("boom")
	svc := service.NewLeadService(repo, newFakeClientRepository(), nil, nil, nil)

	_, err := svc.List(context.Background())

	require.Error(t, err)
}
