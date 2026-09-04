package service_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type fakeClientLookupRepository struct {
	byID map[int64]model.Client
}

func newFakeClientLookupRepository(clients ...model.Client) *fakeClientLookupRepository {
	byID := map[int64]model.Client{}
	for _, c := range clients {
		byID[c.ID] = c
	}
	return &fakeClientLookupRepository{byID: byID}
}

func (f *fakeClientLookupRepository) FindByID(ctx context.Context, id int64) (model.Client, error) {
	c, ok := f.byID[id]
	if !ok {
		return model.Client{}, service.ErrNotFound
	}
	return c, nil
}

func (f *fakeClientLookupRepository) List(ctx context.Context) ([]model.ClientListItem, error) {
	items := make([]model.ClientListItem, 0, len(f.byID))
	for _, c := range f.byID {
		items = append(items, model.ClientListItem{Client: c})
	}
	return items, nil
}

type fakeClientLeadRepository struct {
	byClientID map[int64][]model.Lead
}

func (f *fakeClientLeadRepository) ListByClientID(ctx context.Context, clientID int64) ([]model.Lead, error) {
	return f.byClientID[clientID], nil
}

// fakeTagTypeRepository is an in-memory stand-in for one of the two
// independent tag tables (product or client-type) — one instance per type,
// exactly like the real TagRepository is constructed twice in production.
type fakeTagTypeRepository struct {
	byID             map[int64]model.Tag
	nextID           int64
	assignedByClient map[int64][]int64
}

func newFakeTagTypeRepository() *fakeTagTypeRepository {
	return &fakeTagTypeRepository{byID: map[int64]model.Tag{}, nextID: 1, assignedByClient: map[int64][]int64{}}
}

func (f *fakeTagTypeRepository) FindOrCreateByName(ctx context.Context, name string) (model.Tag, error) {
	for _, tag := range f.byID {
		if strings.EqualFold(tag.Name, name) {
			return tag, nil
		}
	}
	tag := model.Tag{ID: f.nextID, Name: name}
	f.byID[tag.ID] = tag
	f.nextID++
	return tag, nil
}

func (f *fakeTagTypeRepository) SetForClient(ctx context.Context, clientID int64, tagIDs []int64) error {
	f.assignedByClient[clientID] = tagIDs
	return nil
}

func (f *fakeTagTypeRepository) ListForClient(ctx context.Context, clientID int64) ([]model.Tag, error) {
	ids := f.assignedByClient[clientID]
	tags := make([]model.Tag, 0, len(ids))
	for _, id := range ids {
		tags = append(tags, f.byID[id])
	}
	return tags, nil
}

type fakeClientCommentRepository struct {
	byClientID map[int64][]model.ClientComment
	nextID     int64
}

func newFakeClientCommentRepository() *fakeClientCommentRepository {
	return &fakeClientCommentRepository{byClientID: map[int64][]model.ClientComment{}, nextID: 1}
}

func (f *fakeClientCommentRepository) ListForClient(ctx context.Context, clientID int64) ([]model.ClientComment, error) {
	return f.byClientID[clientID], nil
}

func (f *fakeClientCommentRepository) Create(ctx context.Context, item model.ClientComment) (model.ClientComment, error) {
	item.ID = f.nextID
	f.nextID++
	f.byClientID[item.ClientID] = append(f.byClientID[item.ClientID], item)
	return item, nil
}

type fakeReminderRepository struct {
	byClientID map[int64][]model.Reminder
	byID       map[int64]*model.Reminder
	nextID     int64
}

func newFakeReminderRepository() *fakeReminderRepository {
	return &fakeReminderRepository{byClientID: map[int64][]model.Reminder{}, byID: map[int64]*model.Reminder{}, nextID: 1}
}

func (f *fakeReminderRepository) ListForClient(ctx context.Context, clientID int64) ([]model.Reminder, error) {
	return f.byClientID[clientID], nil
}

func (f *fakeReminderRepository) Create(ctx context.Context, item model.Reminder) (model.Reminder, error) {
	item.ID = f.nextID
	f.nextID++
	f.byClientID[item.ClientID] = append(f.byClientID[item.ClientID], item)
	f.byID[item.ID] = &f.byClientID[item.ClientID][len(f.byClientID[item.ClientID])-1]
	return item, nil
}

func (f *fakeReminderRepository) Complete(ctx context.Context, id int64) error {
	r, ok := f.byID[id]
	if !ok {
		return service.ErrNotFound
	}
	now := time.Now()
	r.CompletedAt = &now
	return nil
}

func newTestClientService() (*service.ClientService, *fakeClientLookupRepository, *fakeClientLeadRepository, *fakeTagTypeRepository, *fakeTagTypeRepository, *fakeClientCommentRepository, *fakeReminderRepository) {
	clients := newFakeClientLookupRepository(model.Client{ID: 1, Name: "Иван"})
	leads := &fakeClientLeadRepository{byClientID: map[int64][]model.Lead{}}
	productTags := newFakeTagTypeRepository()
	typeTags := newFakeTagTypeRepository()
	comments := newFakeClientCommentRepository()
	reminders := newFakeReminderRepository()
	svc := service.NewClientService(clients, leads, productTags, typeTags, comments, reminders)
	return svc, clients, leads, productTags, typeTags, comments, reminders
}

func TestClientService_List(t *testing.T) {
	svc, _, _, _, _, _, _ := newTestClientService()

	items, err := svc.List(context.Background())

	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "Иван", items[0].Name)
}

func TestClientService_SetProductTags(t *testing.T) {
	t.Run("creates a tag on the fly when the name doesn't exist yet", func(t *testing.T) {
		svc, _, _, productTags, _, _, _ := newTestClientService()

		tags, err := svc.SetProductTags(context.Background(), 1, []string{"Свадьба"})

		require.NoError(t, err)
		require.Len(t, tags, 1)
		assert.Equal(t, "Свадьба", tags[0].Name)
		assert.Len(t, productTags.byID, 1)
	})

	t.Run("reuses an existing tag case-insensitively instead of duplicating it", func(t *testing.T) {
		svc, _, _, productTags, _, _, _ := newTestClientService()

		_, err := svc.SetProductTags(context.Background(), 1, []string{"Свадьба"})
		require.NoError(t, err)
		tags, err := svc.SetProductTags(context.Background(), 1, []string{"свадьба"})
		require.NoError(t, err)

		require.Len(t, tags, 1)
		assert.Len(t, productTags.byID, 1, "must not create a second row for the same name in a different case")
	})

	t.Run("does not touch the client-type tag table", func(t *testing.T) {
		svc, _, _, _, typeTags, _, _ := newTestClientService()

		_, err := svc.SetProductTags(context.Background(), 1, []string{"Свадьба"})

		require.NoError(t, err)
		assert.Empty(t, typeTags.byID, "product tags must never land in the client-type tag table")
	})

	t.Run("rejects a zero client id", func(t *testing.T) {
		svc, _, _, _, _, _, _ := newTestClientService()

		_, err := svc.SetProductTags(context.Background(), 0, []string{"Свадьба"})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestClientService_SetClientTypeTags(t *testing.T) {
	t.Run("creates a tag on the fly, independent of product tags", func(t *testing.T) {
		svc, _, _, _, typeTags, _, _ := newTestClientService()

		tags, err := svc.SetClientTypeTags(context.Background(), 1, []string{"Постоянный"})

		require.NoError(t, err)
		require.Len(t, tags, 1)
		assert.Equal(t, "Постоянный", tags[0].Name)
		assert.Len(t, typeTags.byID, 1)
	})
}

func TestClientService_AddComment(t *testing.T) {
	t.Run("adds a comment and stamps it against the client", func(t *testing.T) {
		svc, _, _, _, _, comments, _ := newTestClientService()

		comment, err := svc.AddComment(context.Background(), 1, "  Перезвонить завтра  ")

		require.NoError(t, err)
		assert.Equal(t, "Перезвонить завтра", comment.Text, "text should be trimmed")
		assert.Len(t, comments.byClientID[1], 1)
	})

	t.Run("rejects blank text", func(t *testing.T) {
		svc, _, _, _, _, _, _ := newTestClientService()

		_, err := svc.AddComment(context.Background(), 1, "   ")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestClientService_AddReminder(t *testing.T) {
	t.Run("creates a reminder N days out", func(t *testing.T) {
		svc, _, _, _, _, _, reminders := newTestClientService()

		reminder, err := svc.AddReminder(context.Background(), 1, 3, "перезвонить")

		require.NoError(t, err)
		assert.Len(t, reminders.byClientID[1], 1)
		assert.WithinDuration(t, time.Now().AddDate(0, 0, 3), reminder.RemindAt, 5*time.Second)
	})

	t.Run("rejects a non-positive day count", func(t *testing.T) {
		svc, _, _, _, _, _, _ := newTestClientService()

		_, err := svc.AddReminder(context.Background(), 1, 0, "")

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestClientService_CompleteReminder(t *testing.T) {
	t.Run("marks a reminder complete so it drops out of the due set", func(t *testing.T) {
		svc, _, _, _, _, _, reminders := newTestClientService()
		reminder, err := svc.AddReminder(context.Background(), 1, 1, "")
		require.NoError(t, err)

		err = svc.CompleteReminder(context.Background(), reminder.ID)

		require.NoError(t, err)
		assert.NotNil(t, reminders.byID[reminder.ID].CompletedAt)
	})

	t.Run("rejects a zero id", func(t *testing.T) {
		svc, _, _, _, _, _, _ := newTestClientService()

		err := svc.CompleteReminder(context.Background(), 0)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestClientService_GetDetail(t *testing.T) {
	t.Run("aggregates requests, comments, tags, and reminders for the client", func(t *testing.T) {
		svc, _, leads, _, _, _, _ := newTestClientService()
		leads.byClientID[1] = []model.Lead{{ID: 10, ClientID: 1, Name: "Иван"}}
		_, err := svc.SetProductTags(context.Background(), 1, []string{"Свадьба"})
		require.NoError(t, err)
		_, err = svc.AddComment(context.Background(), 1, "звонили")
		require.NoError(t, err)
		_, err = svc.AddReminder(context.Background(), 1, 3, "")
		require.NoError(t, err)

		detail, err := svc.GetDetail(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, int64(1), detail.ID)
		require.Len(t, detail.Requests, 1)
		require.Len(t, detail.Comments, 1)
		require.Len(t, detail.ProductTags, 1)
		require.Len(t, detail.Reminders, 1)
	})

	t.Run("propagates a not-found client", func(t *testing.T) {
		svc, _, _, _, _, _, _ := newTestClientService()

		_, err := svc.GetDetail(context.Background(), 999)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrNotFound)
	})

	t.Run("rejects a zero id", func(t *testing.T) {
		svc, _, _, _, _, _, _ := newTestClientService()

		_, err := svc.GetDetail(context.Background(), 0)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}
