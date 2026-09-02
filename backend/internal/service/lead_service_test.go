package service_test

import (
	"context"
	"errors"
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

func (f *fakeLeadRepository) List(ctx context.Context) ([]model.Lead, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
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
			return f.items[i], nil
		}
	}
	return model.Lead{}, errors.New("not found")
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
		svc := service.NewLeadService(repo, nil, nil, nil)

		lead, err := svc.Create(context.Background(), validLead())

		require.NoError(t, err)
		assert.Equal(t, int64(1), lead.ID)
		assert.Equal(t, "Иван Иванов", lead.Name)
		assert.Equal(t, model.LeadStatusNew, lead.Status)
	})

	t.Run("trims name and phone", func(t *testing.T) {
		repo := newFakeLeadRepository()
		svc := service.NewLeadService(repo, nil, nil, nil)

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
		svc := service.NewLeadService(repo, nil, nil, nil)

		item := validLead()
		item.Name = "   "

		_, err := svc.Create(context.Background(), item)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})

	t.Run("rejects an empty phone", func(t *testing.T) {
		repo := newFakeLeadRepository()
		svc := service.NewLeadService(repo, nil, nil, nil)

		item := validLead()
		item.Phone = "   "

		_, err := svc.Create(context.Background(), item)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects an invalid contact method", func(t *testing.T) {
		repo := newFakeLeadRepository()
		svc := service.NewLeadService(repo, nil, nil, nil)

		item := validLead()
		item.ContactMethod = "carrier_pigeon"

		_, err := svc.Create(context.Background(), item)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects an invalid source", func(t *testing.T) {
		repo := newFakeLeadRepository()
		svc := service.NewLeadService(repo, nil, nil, nil)

		item := validLead()
		item.Source = "unknown"

		_, err := svc.Create(context.Background(), item)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects an invalid request type", func(t *testing.T) {
		repo := newFakeLeadRepository()
		svc := service.NewLeadService(repo, nil, nil, nil)

		item := validLead()
		item.RequestType = "unknown"

		_, err := svc.Create(context.Background(), item)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("forces status to new regardless of input", func(t *testing.T) {
		repo := newFakeLeadRepository()
		svc := service.NewLeadService(repo, nil, nil, nil)

		item := validLead()
		item.Status = model.LeadStatusClosed

		lead, err := svc.Create(context.Background(), item)

		require.NoError(t, err)
		assert.Equal(t, model.LeadStatusNew, lead.Status)
	})
}

func TestLeadService_Create_Notifications(t *testing.T) {
	t.Run("notifies with the created lead on success", func(t *testing.T) {
		repo := newFakeLeadRepository()
		notifier := &fakeLeadNotifier{}
		svc := service.NewLeadService(repo, notifier, nil, nil)

		lead, err := svc.Create(context.Background(), validLead())

		require.NoError(t, err)
		require.Len(t, notifier.notified, 1)
		assert.Equal(t, lead.ID, notifier.notified[0].ID, "must be notified with the ID assigned by the repository, not the pre-insert value")
	})

	t.Run("a notification failure does not fail lead creation", func(t *testing.T) {
		repo := newFakeLeadRepository()
		notifier := &fakeLeadNotifier{err: errors.New("smtp relay unreachable")}
		svc := service.NewLeadService(repo, notifier, nil, nil)

		lead, err := svc.Create(context.Background(), validLead())

		require.NoError(t, err, "the lead is already saved — a notification hiccup must not turn a successful submission into an error")
		assert.NotZero(t, lead.ID)
	})

	t.Run("does not notify when validation rejects the lead", func(t *testing.T) {
		repo := newFakeLeadRepository()
		notifier := &fakeLeadNotifier{}
		svc := service.NewLeadService(repo, notifier, nil, nil)

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
		svc := service.NewLeadService(repo, notifier, nil, nil)

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
		svc := service.NewLeadService(repo, notifier, courses, nil)

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
		svc := service.NewLeadService(repo, notifier, nil, masterclasses)

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
		svc := service.NewLeadService(repo, notifier, courses, nil)

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
		svc := service.NewLeadService(repo, notifier, nil, nil)

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
	svc := service.NewLeadService(repo, nil, nil, nil)
	created, err := svc.Create(context.Background(), validLead())
	require.NoError(t, err)

	t.Run("updates status of an existing lead", func(t *testing.T) {
		lead, err := svc.UpdateStatus(context.Background(), created.ID, model.LeadStatusInProgress)

		require.NoError(t, err)
		assert.Equal(t, model.LeadStatusInProgress, lead.Status)
	})

	t.Run("rejects an invalid status", func(t *testing.T) {
		_, err := svc.UpdateStatus(context.Background(), created.ID, "unknown")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("rejects a zero id", func(t *testing.T) {
		_, err := svc.UpdateStatus(context.Background(), 0, model.LeadStatusClosed)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestLeadService_Delete(t *testing.T) {
	repo := newFakeLeadRepository()
	svc := service.NewLeadService(repo, nil, nil, nil)
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
	svc := service.NewLeadService(repo, nil, nil, nil)

	_, err := svc.List(context.Background())

	require.Error(t, err)
}
