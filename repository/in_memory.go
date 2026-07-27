package repository

import (
	"context"
	"errors"
	"sync"

	"go-task-api/domain"
)

type InMemoryTaskRepo struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

func NewInMemoryTaskRepo() *InMemoryTaskRepo {
	return &InMemoryTaskRepo{tasks: make(map[string]*domain.Task)}
}

func (r *InMemoryTaskRepo) Save(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
	return nil
}

func (r *InMemoryTaskRepo) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, exists := r.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}