package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Тесты для TaskHandler
func TestTaskHandler_CreateTask_Success(t *testing.T) {
	repo := NewMockTaskRepo()
	service := NewTaskService(repo)
	handler := NewTaskHandler(service)

	// Создаём HTTP тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.CreateTask(w, r)
	}))
	defer server.Close()

	// Готовим запрос
	body := bytes.NewBufferString(`{"title": "Handler test"}`)
	resp, err := http.Post(server.URL+"/tasks", "application/json", body)

	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	// Проверяем ответ
	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if task.Title != "Handler test" {
		t.Errorf("expected title 'Handler test', got %s", task.Title)
	}
}

func TestTaskHandler_CreateTask_InvalidJSON(t *testing.T) {
	repo := NewMockTaskRepo()
	service := NewTaskService(repo)
	handler := NewTaskHandler(service)

	req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestTaskHandler_CreateTask_EmptyTitle(t *testing.T) {
	repo := NewMockTaskRepo()
	service := NewTaskService(repo)
	handler := NewTaskHandler(service)

	body := bytes.NewBufferString(`{"title": ""}`)
	req := httptest.NewRequest("POST", "/tasks", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestTaskHandler_GetTask_Success(t *testing.T) {
	repo := NewMockTaskRepo()
	service := NewTaskService(repo)
	handler := NewTaskHandler(service)

	// Создаём задачу
	createdTask, _ := service.CreateTask(context.Background(), "Get handler test")

	// Запрашиваем её
	req := httptest.NewRequest("GET", "/tasks/"+createdTask.ID, nil)
	w := httptest.NewRecorder()

	// Используем роутер Go 1.22+
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}", handler.GetTask)
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var task Task
	if err := json.NewDecoder(w.Body).Decode(&task); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if task.ID != createdTask.ID {
		t.Errorf("expected ID %s, got %s", createdTask.ID, task.ID)
	}
}

func TestTaskHandler_GetTask_NotFound(t *testing.T) {
	repo := NewMockTaskRepo()
	service := NewTaskService(repo)
	handler := NewTaskHandler(service)

	req := httptest.NewRequest("GET", "/tasks/nonexistent", nil)
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}", handler.GetTask)
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}
