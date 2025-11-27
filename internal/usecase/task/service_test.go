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
	_, err := service.Create("list-123", "Task", "Description", "urgent")
	if err == nil {
		t.Error("Expected error for invalid priority, got nil")
	}
}

func TestGetAllTasks(t *testing.T) {
	repo := &MockRepository{tasks: []*domain.Task{{ID: "1", Title: "A"}}}
	service := NewService(repo)
	tasks, err := service.GetAll()
	if err != nil || len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %v, err: %v", tasks, err)
	}
}

func TestGetByFilters_Valid(t *testing.T) {
	repo := &MockRepository{tasks: []*domain.Task{{ID: "1", Status: "pending", Priority: "high"}}}
	service := NewService(repo)
	tasks, err := service.GetByFilters("pending", "high")
	if err != nil || len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %v, err: %v", tasks, err)
	}
}

func TestGetByFilters_InvalidStatus(t *testing.T) {
	service := NewService(&MockRepository{})
	_, err := service.GetByFilters("invalid", "high")
	if err == nil {
		t.Error("Expected error for invalid status")
	}
}

func TestGetByFilters_InvalidPriority(t *testing.T) {
	service := NewService(&MockRepository{})
	_, err := service.GetByFilters("pending", "urgent")
	if err == nil {
		t.Error("Expected error for invalid priority")
	}
}

func TestGetByID(t *testing.T) {
	repo := &MockRepository{tasks: []*domain.Task{{ID: "1", Title: "A"}}}
	service := NewService(repo)
	task, err := service.GetByID("1")
	if err != nil || task == nil || task.ID != "1" {
		t.Errorf("Expected id '1', got %v, err: %v", task, err)
	}
}

func TestUpdateTask_Success(t *testing.T) {
	repo := &MockRepository{tasks: []*domain.Task{{ID: "1", Title: "Old", Status: "pending", Priority: "medium"}}}
	service := NewService(repo)
	task, err := service.Update("1", "list-123", "New", "desc", "completed", "high")
	if err != nil || task.Title != "New" || task.Status != "completed" || task.Priority != "high" {
		t.Errorf("Unexpected result: %+v, err: %v", task, err)
	}
}

func TestUpdateTask_EmptyTitle(t *testing.T) {
	service := NewService(&MockRepository{})
	_, err := service.Update("1", "list-123", "", "desc", "pending", "high")
	if err == nil {
		t.Error("Expected error for empty title")
	}
}

func TestUpdateTask_InvalidStatus(t *testing.T) {
	service := NewService(&MockRepository{})
	_, err := service.Update("1", "list-123", "title", "desc", "invalid", "high")
	if err == nil {
		t.Error("Expected error for invalid status")
	}
}

func TestUpdateTask_InvalidPriority(t *testing.T) {
	service := NewService(&MockRepository{})
	_, err := service.Update("1", "list-123", "title", "desc", "pending", "urgent")
	if err == nil {
		t.Error("Expected error for invalid priority")
	}
}

func TestDeleteTask(t *testing.T) {
	repo := &MockRepository{}
	service := NewService(repo)
	err := service.Delete("1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
