package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type GiftCertificateRepository struct {
	db *pgxpool.Pool
}

func NewGiftCertificateRepository(db *pgxpool.Pool) *GiftCertificateRepository {
	return &GiftCertificateRepository{db: db}
}

func (r *GiftCertificateRepository) List(ctx context.Context) ([]model.GiftCertificate, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, number, issued_at, kind, value, recipient, created_at
		FROM gift_certificates
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.GiftCertificate{}
	for rows.Next() {
		var item model.GiftCertificate
		if err := rows.Scan(&item.ID, &item.Number, &item.IssuedAt, &item.Kind, &item.Value, &item.Recipient, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *GiftCertificateRepository) Get(ctx context.Context, id int64) (model.GiftCertificate, error) {
	var item model.GiftCertificate
	err := r.db.QueryRow(ctx, `
		SELECT id, number, issued_at, kind, value, recipient, created_at
		FROM gift_certificates WHERE id = $1
	`, id).Scan(&item.ID, &item.Number, &item.IssuedAt, &item.Kind, &item.Value, &item.Recipient, &item.CreatedAt)
	return item, translateNotFound(err)
}

func (r *GiftCertificateRepository) Create(ctx context.Context, item model.GiftCertificate) (model.GiftCertificate, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO gift_certificates (number, issued_at, kind, value, recipient)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`, item.Number, item.IssuedAt, item.Kind, item.Value, item.Recipient).Scan(&item.ID, &item.CreatedAt)
	return item, err
}

func (r *GiftCertificateRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM gift_certificates WHERE id = $1`, id)
	return checkDeleted(tag, err)
}

// CountIssuedOn returns how many certificates already carry today's DDMMYY
// prefix in their number, so the service can compute the next 2-digit daily
// sequence. Counting by number prefix (not a separate sequence table) keeps
// the number format and the counting logic in one place.
func (r *GiftCertificateRepository) CountIssuedOn(ctx context.Context, datePrefix string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		SELECT count(*) FROM gift_certificates WHERE number LIKE $1 || '%'
	`, "#"+datePrefix).Scan(&count)
	return count, err
}
