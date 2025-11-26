package task

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/G20-00/task-management-service-go/internal/domain"
	"github.com/G20-00/task-management-service-go/pkg/utils"
)

var validStatuses = map[string]bool{
	"pending":     true,
	"in-progress": true,
	"completed":   true,
}

var validPriorities = map[string]bool{
	"low":    true,
	"medium": true,
	"high":   true,
}

// Service implements the task business logic operations.
type Service struct {
	repo Repository
}

// NewService creates and returns a new task Service instance.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Create creates a new task with the provided details and returns the created task.
func (s *Service) Create(listID, title, description, priority string) (task *domain.Task, err error) {
	defer utils.RecoverPanic("service", "Create", &err)

	if strings.TrimSpace(title) == "" {
		return nil, errors.New("title cannot be empty")
	}

	if priority == "" {
		priority = "medium"
	}

	if !validPriorities[priority] {
		return nil, errors.New("invalid priority: must be low, medium, or high")
	}

	now := time.Now()
	newTask := &domain.Task{
		ID:          uuid.New().String(),
		ListID:      listID,
		Title:       title,
		Description: description,
		Status:      "pending",
		Priority:    priority,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(newTask); err != nil {
		return nil, err
	}

	return newTask, nil
}

// GetAll retrieves all tasks from the repository.
func (s *Service) GetAll() (tasks []*domain.Task, err error) {
	defer utils.RecoverPanic("service", "GetAll", &err)

	return s.repo.GetAll()
}

// GetByFilters retrieves tasks filtered by status and/or priority.
func (s *Service) GetByFilters(status, priority string) (tasks []*domain.Task, err error) {
	defer utils.RecoverPanic("service", "GetByFilters", &err)

	if status != "" && !validStatuses[status] {
		return nil, errors.New("invalid status")
	}

	if priority != "" && !validPriorities[priority] {
		return nil, errors.New("invalid priority")
	}

	return s.repo.GetByFilters(status, priority)
}

// GetByID retrieves a task by its ID.
func (s *Service) GetByID(id string) (task *domain.Task, err error) {
	defer utils.RecoverPanic("service", "GetByID", &err)

	return s.repo.GetByID(id)
}

// Update updates an existing task with the provided details and returns the updated task.
func (s *Service) Update(id, listID, title, description, status, priority string) (task *domain.Task, err error) {
	defer utils.RecoverPanic("service", "Update", &err)

	if strings.TrimSpace(title) == "" {
		return nil, errors.New("title cannot be empty")
	}

	if !validStatuses[status] {
		return nil, errors.New("invalid status: must be pending, in-progress, or completed")
	}

	if !validPriorities[priority] {
		return nil, errors.New("invalid priority: must be low, medium, or high")
	}

	existingTask, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	existingTask.ListID = listID
	existingTask.Title = title
	existingTask.Description = description
	existingTask.Status = status
	existingTask.Priority = priority
	existingTask.UpdatedAt = time.Now()

	if err := s.repo.Update(existingTask); err != nil {
		return nil, err
	}

	return existingTask, nil
}

// Delete removes a task by its ID from the repository.
func (s *Service) Delete(id string) (err error) {
	defer utils.RecoverPanic("service", "Delete", &err)

	return s.repo.Delete(id)
}
