package tasklist

import (
	"errors"
	"testing"
	"time"

	"github.com/G20-00/task-management-service-go/internal/domain"
)

type mockRepo struct {
	CreateFn  func(list *domain.TaskList) error
	GetAllFn  func() ([]*domain.TaskList, error)
	GetByIDFn func(id string) (*domain.TaskList, error)
	UpdateFn  func(list *domain.TaskList) error
	DeleteFn  func(id string) error
}

func (m *mockRepo) Create(list *domain.TaskList) error          { return m.CreateFn(list) }
func (m *mockRepo) GetAll() ([]*domain.TaskList, error)         { return m.GetAllFn() }
func (m *mockRepo) GetByID(id string) (*domain.TaskList, error) { return m.GetByIDFn(id) }
func (m *mockRepo) Update(list *domain.TaskList) error          { return m.UpdateFn(list) }
func (m *mockRepo) Delete(id string) error                      { return m.DeleteFn(id) }

func TestService_Create(t *testing.T) {
	repo := &mockRepo{
		CreateFn: func(list *domain.TaskList) error { return nil },
	}
	s := NewService(repo)
	list, err := s.Create("Test List", "desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if list.Name != "Test List" {
		t.Errorf("expected name 'Test List', got %s", list.Name)
	}
}

func TestService_Create_EmptyName(t *testing.T) {
	s := NewService(&mockRepo{})
	_, err := s.Create("", "desc")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestService_GetAll(t *testing.T) {
	repo := &mockRepo{
		GetAllFn: func() ([]*domain.TaskList, error) {
			return []*domain.TaskList{{ID: "1", Name: "L", Description: "D", CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil
		},
	}
	s := NewService(repo)
	lists, err := s.GetAll()
	if err != nil || len(lists) != 1 {
		t.Errorf("expected 1 list, got %v, err: %v", lists, err)
	}
}

func TestService_GetByID(t *testing.T) {
	repo := &mockRepo{
		GetByIDFn: func(id string) (*domain.TaskList, error) {
			return &domain.TaskList{ID: id, Name: "L", Description: "D", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
		},
	}
	s := NewService(repo)
	list, err := s.GetByID("1")
	if err != nil || list.ID != "1" {
		t.Errorf("expected id '1', got %v, err: %v", list, err)
	}
}

func TestService_GetByID_EmptyID(t *testing.T) {
	s := NewService(&mockRepo{})
	_, err := s.GetByID("")
	if err == nil {
		t.Error("expected error for empty id")
	}
}

func TestService_Update(t *testing.T) {
	now := time.Now()
	repo := &mockRepo{
		GetByIDFn: func(id string) (*domain.TaskList, error) {
			return &domain.TaskList{ID: id, Name: "Old", Description: "Old", CreatedAt: now, UpdatedAt: now}, nil
		},
		UpdateFn: func(list *domain.TaskList) error {
			if list.Name != "New" || list.Description != "NewDesc" {
				t.Errorf("unexpected update: %+v", list)
			}
			return nil
		},
	}
	s := NewService(repo)
	list, err := s.Update("1", "New", "NewDesc")
	if err != nil || list.Name != "New" || list.Description != "NewDesc" {
		t.Errorf("unexpected result: %+v, err: %v", list, err)
	}
}

func TestService_Update_EmptyID(t *testing.T) {
	s := NewService(&mockRepo{})
	_, err := s.Update("", "n", "d")
	if err == nil {
		t.Error("expected error for empty id")
	}
}

func TestService_Update_RepoError(t *testing.T) {
	repo := &mockRepo{
		GetByIDFn: func(id string) (*domain.TaskList, error) { return nil, errors.New("not found") },
	}
	s := NewService(repo)
	_, err := s.Update("1", "n", "d")
	if err == nil {
		t.Error("expected error from repo")
	}
}

func TestService_Delete(t *testing.T) {
	repo := &mockRepo{
		DeleteFn: func(id string) error { return nil },
	}
	s := NewService(repo)
	if err := s.Delete("1"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestService_Delete_EmptyID(t *testing.T) {
	s := NewService(&mockRepo{})
	if err := s.Delete(""); err == nil {
		t.Error("expected error for empty id")
	}
}
