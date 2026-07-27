package domain

import "time"

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	IsDone    bool      `json:"is_done"`
	CreatedAt time.Time `json:"created_at"`
}