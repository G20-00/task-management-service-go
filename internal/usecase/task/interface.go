// Package task provides task-related business logic and repository interfaces.
package task

import "github.com/G20-00/task-management-service-go/internal/domain"

// Repository defines the interface for task data persistence operations.
type Repository interface {
	Create(task *domain.Task) error
	GetAll() ([]*domain.Task, error)
	GetByID(id string) (*domain.Task, error)
	Update(task *domain.Task) error
	Delete(id string) error
	GetByFilters(status, priority string) ([]*domain.Task, error)
	CountByListIDAndStatus(listID, status string) (int, error)
}
