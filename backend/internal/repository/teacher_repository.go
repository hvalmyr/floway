package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

type TeacherRepository struct {
	db *pgxpool.Pool
}

func NewTeacherRepository(db *pgxpool.Pool) *TeacherRepository {
	return &TeacherRepository{db: db}
}

func (r *TeacherRepository) List(ctx context.Context) ([]model.Teacher, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, photo, description, sort_order, created_at, updated_at
		FROM teachers
		ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Teacher
	for rows.Next() {
		var item model.Teacher
		if err := rows.Scan(&item.ID, &item.Name, &item.Photo, &item.Description, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *TeacherRepository) Create(ctx context.Context, item model.Teacher) (model.Teacher, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO teachers (name, photo, description, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, item.Name, item.Photo, item.Description, item.SortOrder).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *TeacherRepository) Update(ctx context.Context, item model.Teacher) (model.Teacher, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE teachers
		SET name = $1, photo = $2, description = $3, sort_order = $4, updated_at = now()
		WHERE id = $5
		RETURNING updated_at
	`, item.Name, item.Photo, item.Description, item.SortOrder, item.ID).Scan(&item.UpdatedAt)
	return item, err
}

func (r *TeacherRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM teachers WHERE id = $1`, id)
	return err
}
