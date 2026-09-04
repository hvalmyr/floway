package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type ReminderRepository struct {
	db *pgxpool.Pool
}

func NewReminderRepository(db *pgxpool.Pool) *ReminderRepository {
	return &ReminderRepository{db: db}
}

const reminderColumns = "id, client_id, remind_at, note, completed_at, created_at"

func (r *ReminderRepository) ListForClient(ctx context.Context, clientID int64) ([]model.Reminder, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+reminderColumns+`
		FROM reminders
		WHERE client_id = $1
		ORDER BY remind_at ASC, id ASC
	`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Reminder{}
	for rows.Next() {
		var item model.Reminder
		if err := rows.Scan(&item.ID, &item.ClientID, &item.RemindAt, &item.Note, &item.CompletedAt, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *ReminderRepository) Create(ctx context.Context, item model.Reminder) (model.Reminder, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO reminders (client_id, remind_at, note)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, item.ClientID, item.RemindAt, item.Note).Scan(&item.ID, &item.CreatedAt)
	return item, err
}

func (r *ReminderRepository) Complete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `UPDATE reminders SET completed_at = now() WHERE id = $1`, id)
	return checkDeleted(tag, err)
}
