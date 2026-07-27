# Go Task API

REST API for task management built with Go, following **Clean Architecture** and **SOLID principles**.

This project demonstrates a production-ready approach to building Go services without heavy frameworks — using only the standard library (`net/http`) with modern Go 1.22+ routing.

## Architecture

The project is structured according to **Clean Architecture** with clear separation of concerns:
├── Domain → Business entities (Task)
├── Repository → Data access abstraction (interface + in-memory impl)
├── Service → Business logic layer
├── Handler → HTTP transport layer
└── Main → Dependency injection & composition root


### Layers & Responsibilities

| Layer        | Responsibility                              |
|--------------|---------------------------------------------|
| **Domain**   | Pure business entities, no external deps    |
| **Repository** | Data persistence behind an interface     |
| **Service**  | Business rules, validation, orchestration   |
| **Handler**  | HTTP request/response handling only         |

## SOLID Principles Applied

- **S (Single Responsibility)** — Each layer has one reason to change. Handlers don't know about DB, Services don't know about HTTP.
- **O (Open/Closed)** — New repository implementations (e.g., PostgreSQL) can be added without modifying the service layer.
- **L (Liskov Substitution)** — Any `TaskRepository` implementation can be swapped transparently.
- **I (Interface Segregation)** — Small, focused interfaces (`Save`, `GetByID`).
- **D (Dependency Inversion)** — `TaskService` depends on the `TaskRepository` abstraction, not a concrete implementation.

## Features

- ✅ Pure Go — zero third-party dependencies
- ✅ Go 1.22+ native routing (`r.PathValue`)
- ✅ Thread-safe in-memory storage (`sync.RWMutex`)
- ✅ Context propagation for cancellation & timeouts
- ✅ Proper HTTP status codes & JSON responses
- ✅ Clean, testable architecture ready for PostgreSQL integration

## Tech Stack

- **Go 1.22+**
- **Standard library only** (`net/http`, `encoding/json`, `sync`, `context`)
- **Clean Architecture**
- **SOLID principles**

## Installation & Run

```bash
# Clone the repository
git clone https://github.com/IlyaIvlev/go-task-api.git
cd go-task-api

# Run the server
go run main.go

The server starts on http://localhost:8080.

API Endpoints

Create a Task
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Go concurrency"}'

Response (201 Created):
{
  "id": "task_1722098765432109876",
  "title": "Learn Go concurrency",
  "is_done": false,
  "created_at": "2026-07-27T18:50:00Z"
}

Get a Task by ID
curl http://localhost:8080/tasks/task_1722098765432109876

Response (200 OK):
{
  "id": "task_1722098765432109876",
  "title": "Learn Go concurrency",
  "is_done": false,
  "created_at": "2026-07-27T18:50:00Z"
}

Testing Strategy
The architecture is designed for easy testing:
Unit tests for the Service layer using mock repositories
Integration tests for Handlers using httptest
Repository tests with both in-memory and real DB implementations
Next Steps
Add PostgreSQL repository implementation (pgx)
Add task update & delete endpoints
Add structured logging (slog)
Add graceful shutdown
Add Dockerfile & docker-compose
Add integration tests with Testcontainers
Author: Ilya Ivlev
GitHub: github.com/IlyaIvlev