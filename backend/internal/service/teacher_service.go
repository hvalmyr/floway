package service

import (
	"context"
	"errors"
	"strings"

	"floway-backend/internal/model"
)

type TeacherRepository interface {
	List(ctx context.Context) ([]model.Teacher, error)
	Create(ctx context.Context, item model.Teacher) (model.Teacher, error)
	Update(ctx context.Context, item model.Teacher) (model.Teacher, error)
	Delete(ctx context.Context, id int64) error
}

type TeacherService struct {
	repo TeacherRepository
}

func NewTeacherService(repo TeacherRepository) *TeacherService {
	return &TeacherService{repo: repo}
}

func (s *TeacherService) List(ctx context.Context) ([]model.Teacher, error) {
	return s.repo.List(ctx)
}

func (s *TeacherService) Create(ctx context.Context, item model.Teacher) (model.Teacher, error) {
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return model.Teacher{}, errors.Join(ErrValidation, errors.New("name is required"))
	}

	return s.repo.Create(ctx, item)
}

func (s *TeacherService) Update(ctx context.Context, item model.Teacher) (model.Teacher, error) {
	item.Name = strings.TrimSpace(item.Name)
	if item.ID == 0 {
		return model.Teacher{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if item.Name == "" {
		return model.Teacher{}, errors.Join(ErrValidation, errors.New("name is required"))
	}

	return s.repo.Update(ctx, item)
}

func (s *TeacherService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	return s.repo.Delete(ctx, id)
}
