package repository

import (
	"context"
	"go-task-api/domain"
)

type TaskRepository interface {
	Save(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id string) (*domain.Task, error)
}