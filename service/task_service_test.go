package service

import (
	"context"
	"errors"
	"testing"

	"go-task-api/domain"
)

type MockTaskRepo struct {
	tasks map[string]*domain.Task
}

func NewMockTaskRepo() *MockTaskRepo {
	return &MockTaskRepo{tasks: make(map[string]*domain.Task)}
}

func (m *MockTaskRepo) Save(ctx context.Context, task *domain.Task) error {
	m.tasks[task.ID] = task
	return nil
}

func (m *MockTaskRepo) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func TestTaskService_CreateTask_Success(t *testing.T) {
	repo := NewMockTaskRepo()
	service := NewTaskService(repo)

	task, err := service.CreateTask(context.Background(), "Test task")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task == nil {
		t.Fatal("expected task, got nil")
	}

	if task.Title != "Test task" {
		t.Errorf("expected title 'Test task', got %s", task.Title)
	}

	if task.IsDone != false {
		t.Errorf("expected IsDone to be false, got %v", task.IsDone)
	}

	if task.ID == "" {
		t.Error("expected ID to be set")
	}
}

func TestTaskService_CreateTask_EmptyTitle(t *testing.T) {
	repo := NewMockTaskRepo()
	service := NewTaskService(repo)

	task, err := service.CreateTask(context.Background(), "")

	if err == nil {
		t.Fatal("expected error for empty title, got nil")
	}

	if err.Error() != "title cannot be empty" {
		t.Errorf("expected 'title cannot be empty', got %v", err)
	}

	if task != nil {
		t.Error("expected nil task on error")
	}
}

func TestTaskService_GetTask_Success(t *testing.T) {
	repo := NewMockTaskRepo()
	service := NewTaskService(repo)

	createdTask, _ := service.CreateTask(context.Background(), "Get test")

	task, err := service.GetTask(context.Background(), createdTask.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.ID != createdTask.ID {
		t.Errorf("expected ID %s, got %s", createdTask.ID, task.ID)
	}
}

func TestTaskService_GetTask_NotFound(t *testing.T) {
	repo := NewMockTaskRepo()
	service := NewTaskService(repo)

	task, err := service.GetTask(context.Background(), "nonexistent")

	if err == nil {
		t.Fatal("expected error for nonexistent task, got nil")
	}

	if task != nil {
		t.Error("expected nil task on error")
	}
}
