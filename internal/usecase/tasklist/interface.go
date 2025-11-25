package tasklist

import "github.com/G20-00/task-management-service-go/internal/domain"

type Repository interface {
	Create(list *domain.TaskList) error
	GetAll() ([]*domain.TaskList, error)
	GetByID(id string) (*domain.TaskList, error)
	Update(list *domain.TaskList) error
	Delete(id string) error
}
