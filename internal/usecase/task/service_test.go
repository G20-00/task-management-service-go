package task

import (
	"testing"

	"github.com/G20-00/task-management-service-go/internal/domain"
)

// MockRepository is a mock implementation of the task repository
type MockRepository struct {
	tasks []*domain.Task
}

func (m *MockRepository) Create(task *domain.Task) error {
	m.tasks = append(m.tasks, task)
	return nil
}

func (m *MockRepository) GetAll() ([]*domain.Task, error) {
	return m.tasks, nil
}

func (m *MockRepository) GetByID(id string) (*domain.Task, error) {
	for _, t := range m.tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, nil
}

func (m *MockRepository) Update(task *domain.Task) error {
	return nil
}

func (m *MockRepository) Delete(id string) error {
	return nil
}

func (m *MockRepository) GetByFilters(status, priority string) ([]*domain.Task, error) {
	return m.tasks, nil
}

func (m *MockRepository) CountByListIDAndStatus(listID, status string) (int, error) {
	return 0, nil
}

func TestCreateTask_Success(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	task, err := service.Create("list-123", "Test Task", "Description", "high")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if task.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%s'", task.Title)
	}

	if task.Status != "pending" {
		t.Errorf("Expected status 'pending', got '%s'", task.Status)
	}
}

func TestCreateTask_EmptyTitle(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	_, err := service.Create("list-123", "", "Description", "high")

	if err == nil {
		t.Error("Expected error for empty title, got nil")
	}
}

func TestCreateTask_InvalidPriority(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)

	_, err := service.Create("list-123", "Test Task", "Description", "invalid")

	// This test will FAIL intentionally to test CI
	if err != nil {
		t.Error("Expected no error, but got one (intentional fail for CI test)")
	}
}
