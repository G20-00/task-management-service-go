// Package usecase defines the business logic interfaces for the application.
package usecase

import "github.com/G20-00/task-management-service-go/internal/domain"

// TaskUsecase defines the interface for task business logic operations.
type TaskUsecase interface {
	CreateTask(task *domain.Task) error
	GetTaskByID(id string) (*domain.Task, error)
	GetAllTasks() ([]*domain.Task, error)
	UpdateTask(task *domain.Task) error
	DeleteTask(id string) error
}
