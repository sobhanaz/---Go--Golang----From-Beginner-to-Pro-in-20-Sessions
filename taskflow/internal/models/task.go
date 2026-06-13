// Package models holds the domain types for TaskFlow.
package models

import "time"

// Task is the core domain entity.
// JSON tags define the API's wire format (Session 15).
type Task struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateTaskInput is the payload accepted when creating a task.
// Keeping it separate from Task means clients can't set ID or CreatedAt.
type CreateTaskInput struct {
	Title string `json:"title"`
}

// UpdateTaskInput is the payload accepted when updating a task.
type UpdateTaskInput struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}
