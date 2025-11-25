package repository

import "github.com/G20-00/task-management-service-go/internal/domain"

type TaskRepository interface {
	Create(task *domain.Task) error
	GetByID(id string) (*domain.Task, error)
	GetAll() ([]*domain.Task, error)
	Update(task *domain.Task) error
	Delete(id string) error
}
