package main

import (
	"log"
	"net/http"

	"go-task-api/handler"
	"go-task-api/repository"
	"go-task-api/service"
)

func main() {
	repo := repository.NewInMemoryTaskRepo()
	taskService := service.NewTaskService(repo)
	taskHandler := handler.NewTaskHandler(taskService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", taskHandler.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", taskHandler.GetTask)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
