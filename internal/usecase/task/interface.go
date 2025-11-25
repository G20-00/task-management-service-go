package task

import "github.com/G20-00/task-management-service-go/internal/domain"

type Repository interface {
	Create(task *domain.Task) error
	GetAll() ([]*domain.Task, error)
	GetByID(id string) (*domain.Task, error)
	Update(task *domain.Task) error
	Delete(id string) error
}
