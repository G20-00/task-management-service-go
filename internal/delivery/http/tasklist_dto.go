package http

import "time"

type CreateTaskListRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateTaskListRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TaskListResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
