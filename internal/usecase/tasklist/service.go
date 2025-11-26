package tasklist

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/G20-00/task-management-service-go/internal/domain"
)

// Service implements the task list business logic operations.
type Service struct {
	repo Repository
}

// NewService creates and returns a new task list Service instance.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Create creates a new task list with the provided name and description.
func (s *Service) Create(name, description string) (*domain.TaskList, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name cannot be empty")
	}

	now := time.Now()
	list := &domain.TaskList{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(list); err != nil {
		return nil, err
	}

	return list, nil
}

// GetAll retrieves all task lists from the repository.
func (s *Service) GetAll() ([]*domain.TaskList, error) {
	return s.repo.GetAll()
}

// GetByID retrieves a task list by its ID.
func (s *Service) GetByID(id string) (*domain.TaskList, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	return s.repo.GetByID(id)
}

// Update updates an existing task list with the provided details.
func (s *Service) Update(id, name, description string) (*domain.TaskList, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		existing.Name = name
	}
	if description != "" {
		existing.Description = description
	}
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// Delete removes a task list by its ID from the repository.
func (s *Service) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	return s.repo.Delete(id)
}
