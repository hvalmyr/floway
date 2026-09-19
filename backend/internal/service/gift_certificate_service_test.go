package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/model"
	"floway-backend/internal/service"
)

type fakeGiftCertificateRepository struct {
	items  []model.GiftCertificate
	nextID int64
	err    error
}

func newFakeGiftCertificateRepository() *fakeGiftCertificateRepository {
	return &fakeGiftCertificateRepository{nextID: 1}
}

func (f *fakeGiftCertificateRepository) List(ctx context.Context) ([]model.GiftCertificate, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func (f *fakeGiftCertificateRepository) Get(ctx context.Context, id int64) (model.GiftCertificate, error) {
	for _, existing := range f.items {
		if existing.ID == id {
			return existing, nil
		}
	}
	return model.GiftCertificate{}, errors.New("not found")
}

func (f *fakeGiftCertificateRepository) Create(ctx context.Context, item model.GiftCertificate) (model.GiftCertificate, error) {
	if f.err != nil {
		return model.GiftCertificate{}, f.err
	}
	item.ID = f.nextID
	f.nextID++
	f.items = append(f.items, item)
	return item, nil
}

func (f *fakeGiftCertificateRepository) Delete(ctx context.Context, id int64) error {
	for i, existing := range f.items {
		if existing.ID == id {
			f.items = append(f.items[:i], f.items[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func (f *fakeGiftCertificateRepository) CountIssuedOn(ctx context.Context, datePrefix string) (int, error) {
	if f.err != nil {
		return 0, f.err
	}
	count := 0
	prefix := "#" + datePrefix
	for _, existing := range f.items {
		if len(existing.Number) >= len(prefix) && existing.Number[:len(prefix)] == prefix {
			count++
		}
	}
	return count, nil
}

func TestGiftCertificateService_Create(t *testing.T) {
	t.Run("first and second certificate of the day get sequential numbers", func(t *testing.T) {
		repo := newFakeGiftCertificateRepository()
		svc := service.NewGiftCertificateService(repo)
		wantPrefix := "#" + time.Now().Format("020106")

		first, err := svc.Create(context.Background(), model.GiftCertificate{
			Kind:      model.GiftCertificateKindAmount,
			Value:     "5000 рублей",
			Recipient: "Васильевой Василисе Васильевне",
		})
		require.NoError(t, err)
		assert.Equal(t, wantPrefix+"01", first.Number)

		second, err := svc.Create(context.Background(), model.GiftCertificate{
			Kind:      model.GiftCertificateKindAnyMasterclass,
			Recipient: "Ивановой Марии Петровне",
		})
		require.NoError(t, err)
		assert.Equal(t, wantPrefix+"02", second.Number)
	})

	t.Run("rejects an invalid kind", func(t *testing.T) {
		repo := newFakeGiftCertificateRepository()
		svc := service.NewGiftCertificateService(repo)

		_, err := svc.Create(context.Background(), model.GiftCertificate{
			Kind:      model.GiftCertificateKind("bogus"),
			Value:     "5000 рублей",
			Recipient: "Кто-то",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
		assert.Empty(t, repo.items)
	})

	t.Run("rejects a missing value for a kind that needs one", func(t *testing.T) {
		repo := newFakeGiftCertificateRepository()
		svc := service.NewGiftCertificateService(repo)

		_, err := svc.Create(context.Background(), model.GiftCertificate{
			Kind:      model.GiftCertificateKindCourse,
			Recipient: "Кто-то",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})

	t.Run("any_masterclass kind doesn't need a value", func(t *testing.T) {
		repo := newFakeGiftCertificateRepository()
		svc := service.NewGiftCertificateService(repo)

		_, err := svc.Create(context.Background(), model.GiftCertificate{
			Kind:      model.GiftCertificateKindAnyMasterclass,
			Recipient: "Кто-то",
		})

		require.NoError(t, err)
	})

	t.Run("rejects a missing recipient", func(t *testing.T) {
		repo := newFakeGiftCertificateRepository()
		svc := service.NewGiftCertificateService(repo)

		_, err := svc.Create(context.Background(), model.GiftCertificate{
			Kind:  model.GiftCertificateKindAnyMasterclass,
			Value: "",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrValidation)
	})
}

func TestGiftCertificateService_Delete(t *testing.T) {
	repo := newFakeGiftCertificateRepository()
	svc := service.NewGiftCertificateService(repo)
	created, err := svc.Create(context.Background(), model.GiftCertificate{
		Kind:      model.GiftCertificateKindAnyMasterclass,
		Recipient: "Кто-то",
	})
	require.NoError(t, err)

	t.Run("deletes an existing item", func(t *testing.T) {
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
