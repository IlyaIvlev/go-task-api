package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// --- 1. DOMAIN (Сущность) ---

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	IsDone    bool      `json:"is_done"`
	CreatedAt time.Time `json:"created_at"`
}

// --- 2. REPOSITORY (Интерфейс и реализация) ---
// Принцип I (Interface Segregation) и D (Dependency Inversion)

type TaskRepository interface {
	Save(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id string) (*Task, error)
}

type InMemoryTaskRepo struct {
	mu    sync.RWMutex
	tasks map[string]*Task
}

func NewInMemoryTaskRepo() *InMemoryTaskRepo {
	return &InMemoryTaskRepo{tasks: make(map[string]*Task)}
}

func (r *InMemoryTaskRepo) Save(ctx context.Context, task *Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
	return nil
}

func (r *InMemoryTaskRepo) GetByID(ctx context.Context, id string) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, exists := r.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

// --- 3. SERVICE (Бизнес-логика) ---
// Принцип S (Single Responsibility) и D (Dependency Inversion)

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, title string) (*Task, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	task := &Task{
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

func (s *TaskService) GetTask(ctx context.Context, id string) (*Task, error) {
	return s.repo.GetByID(ctx, id)
}

// --- 4. HANDLER (HTTP слой) ---
// Принцип S (Single Responsibility)

type TaskHandler struct {
	service *TaskService
}

func NewTaskHandler(service *TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	task, err := h.service.CreateTask(r.Context(), req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") // Go 1.22+ routing
	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// --- 5. MAIN (Точка входа и Dependency Injection) ---

func main() {
	// Инициализация зависимостей (Dependency Injection)
	repo := NewInMemoryTaskRepo()
	service := NewTaskService(repo)
	handler := NewTaskHandler(service)

	// Настройка роутера (Go 1.22+ синтаксис)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", handler.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", handler.GetTask)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}