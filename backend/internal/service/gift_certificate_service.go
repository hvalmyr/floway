package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"floway-backend/internal/model"
)

type GiftCertificateRepository interface {
	List(ctx context.Context) ([]model.GiftCertificate, error)
	Get(ctx context.Context, id int64) (model.GiftCertificate, error)
	Create(ctx context.Context, item model.GiftCertificate) (model.GiftCertificate, error)
	Delete(ctx context.Context, id int64) error
	CountIssuedOn(ctx context.Context, datePrefix string) (int, error)
}

var validGiftCertificateKinds = map[model.GiftCertificateKind]bool{
	model.GiftCertificateKindAmount:         true,
	model.GiftCertificateKindCourse:         true,
	model.GiftCertificateKindMasterclass:    true,
	model.GiftCertificateKindAnyMasterclass: true,
}

type GiftCertificateService struct {
	repo GiftCertificateRepository
}

func NewGiftCertificateService(repo GiftCertificateRepository) *GiftCertificateService {
	return &GiftCertificateService{repo: repo}
}

func (s *GiftCertificateService) List(ctx context.Context) ([]model.GiftCertificate, error) {
	return s.repo.List(ctx)
}

func (s *GiftCertificateService) Get(ctx context.Context, id int64) (model.GiftCertificate, error) {
	return s.repo.Get(ctx, id)
}

// Create validates the admin's input, then generates the certificate number
// itself — Number/IssuedAt on the incoming item are ignored, never
// client-supplied. Numbering is #DDMMYYNN: day/month/2-digit-year of the
// issue date, then a 1-based 2-digit same-day sequence from
// CountIssuedOn+1. This is a single-admin system (see ClientComment's doc
// comment in model.go) — read-count-then-insert has a theoretical race only
// if two admins issue a certificate in the same instant, which doesn't
// happen in practice here, so no advisory lock/transaction is used. If it
// ever did collide, the UNIQUE index on number turns it into a request
// error instead of silent data corruption.
func (s *GiftCertificateService) Create(ctx context.Context, item model.GiftCertificate) (model.GiftCertificate, error) {
	if !validGiftCertificateKinds[item.Kind] {
		return model.GiftCertificate{}, errors.Join(ErrValidation, errors.New("invalid kind"))
	}
	if item.Kind != model.GiftCertificateKindAnyMasterclass && item.Value == "" {
		return model.GiftCertificate{}, errors.Join(ErrValidation, errors.New("value is required for this kind"))
	}
	if item.Recipient == "" {
		return model.GiftCertificate{}, errors.Join(ErrValidation, errors.New("recipient is required"))
	}

	issuedAt := time.Now()
	item.IssuedAt = issuedAt

	datePrefix := issuedAt.Format("020106") // DDMMYY
	count, err := s.repo.CountIssuedOn(ctx, datePrefix)
	if err != nil {
		return model.GiftCertificate{}, err
	}
	item.Number = fmt.Sprintf("#%s%02d", datePrefix, count+1)

	return s.repo.Create(ctx, item)
}

func (s *GiftCertificateService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
