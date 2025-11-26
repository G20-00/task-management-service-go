package http

import "time"

// CreateTaskListRequest represents the request body for creating a task list.
type CreateTaskListRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateTaskListRequest represents the request body for updating a task list.
type UpdateTaskListRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TaskListResponse represents the response body for a task list.
type TaskListResponse struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	CompletionPercentage float64   `json:"completion_percentage"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
