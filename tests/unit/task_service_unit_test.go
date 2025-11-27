package unit

import (
	"errors"
	"testing"

	"github.com/G20-00/task-management-service-go/internal/domain"
	taskusecase "github.com/G20-00/task-management-service-go/internal/usecase/task"
)

type mockRepo struct {
	CountByListIDAndStatusFn func(string, string) (int, error)
	CreateFn                 func(*domain.Task) error
	GetByIDFn                func(string) (*domain.Task, error)
	UpdateFn                 func(*domain.Task) error
	DeleteFn                 func(string) error
}

func (m *mockRepo) Create(t *domain.Task) error                         { return m.CreateFn(t) }
func (m *mockRepo) GetByID(id string) (*domain.Task, error)             { return m.GetByIDFn(id) }
func (m *mockRepo) Update(t *domain.Task) error                         { return m.UpdateFn(t) }
func (m *mockRepo) Delete(id string) error                              { return m.DeleteFn(id) }
func (m *mockRepo) GetAll() ([]*domain.Task, error)                     { return nil, nil }
func (m *mockRepo) GetByFilters(string, string) ([]*domain.Task, error) { return nil, nil }
func (m *mockRepo) CountByListIDAndStatus(listID, status string) (int, error) {
	if m.CountByListIDAndStatusFn != nil {
		return m.CountByListIDAndStatusFn(listID, status)
	}
	return 0, nil
}

func TestService_Create_Success(t *testing.T) {
	repo := &mockRepo{
		CreateFn: func(tk *domain.Task) error { return nil },
	}
	svc := taskusecase.NewService(repo)
	task, err := svc.Create("list-1", "titulo", "desc", "medium")
	if err != nil || task == nil {
		t.Fatalf("esperado crear tarea sin error, obtuve %v", err)
	}
}

func TestService_Create_InvalidPriority(t *testing.T) {
	repo := &mockRepo{CreateFn: func(tk *domain.Task) error { return nil }}
	svc := taskusecase.NewService(repo)
	_, err := svc.Create("list-1", "titulo", "desc", "super")
	if err == nil {
		t.Error("esperado error por prioridad inválida")
	}
}

func TestService_Update_NotFound(t *testing.T) {
	repo := &mockRepo{
		GetByIDFn: func(id string) (*domain.Task, error) { return nil, errors.New("task not found") },
		UpdateFn:  func(tk *domain.Task) error { return nil },
	}
	svc := taskusecase.NewService(repo)
	_, err := svc.Update("id", "list-1", "titulo", "desc", "pending", "medium")
	if err == nil {
		t.Error("esperado error por tarea no encontrada")
	}
}

func TestService_Update_InvalidStatus(t *testing.T) {
	repo := &mockRepo{
		GetByIDFn: func(id string) (*domain.Task, error) { return &domain.Task{ID: id}, nil },
		UpdateFn:  func(tk *domain.Task) error { return nil },
	}
	svc := taskusecase.NewService(repo)
	_, err := svc.Update("id", "list-1", "titulo", "desc", "hecho", "medium")
	if err == nil {
		t.Error("esperado error por status inválido")
	}
}
