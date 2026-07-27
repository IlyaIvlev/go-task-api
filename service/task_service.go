package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-task-api/domain"
	"go-task-api/repository"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, title string) (*domain.Task, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	task := &domain.Task{
		ID:        fmt.Sprintf("task_%d", time.Now().UnixNano()),
		Title:     title,
		IsDone:    false,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Save(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	return s.repo.GetByID(ctx, id)
}
