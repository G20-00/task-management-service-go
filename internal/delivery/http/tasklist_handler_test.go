package http

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/G20-00/task-management-service-go/internal/domain"
)

type mockTaskListService struct {
	CreateFn  func(name, description string) (*domain.TaskList, error)
	GetAllFn  func() ([]*domain.TaskList, error)
	GetByIDFn func(id string) (*domain.TaskList, error)
	UpdateFn  func(id, name, description string) (*domain.TaskList, error)
	DeleteFn  func(id string) error
}

func (m *mockTaskListService) Create(name, description string) (*domain.TaskList, error) {
	if m.CreateFn != nil {
		return m.CreateFn(name, description)
	}
	return nil, nil
}
func (m *mockTaskListService) GetAll() ([]*domain.TaskList, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn()
	}
	return nil, nil
}
func (m *mockTaskListService) GetByID(id string) (*domain.TaskList, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return nil, nil
}
func (m *mockTaskListService) Update(id, name, description string) (*domain.TaskList, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(id, name, description)
	}
	return nil, nil
}
func (m *mockTaskListService) Delete(id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

func TestCreateTaskList_Success(t *testing.T) {
	app := fiber.New()
	h := &TaskListHandler{service: &mockTaskListService{
		CreateFn: func(name, description string) (*domain.TaskList, error) {
			return &domain.TaskList{ID: "1", Name: name, Description: description, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
		},
	}}
	app.Post("/lists", h.CreateTaskList)
	body := `{"name":"Lista","description":"desc"}`
	req := httptest.NewRequest("POST", "/lists", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}
}

func TestCreateTaskList_BadRequest(t *testing.T) {
	app := fiber.New()
	h := &TaskListHandler{service: &mockTaskListService{}}
	app.Post("/lists", h.CreateTaskList)
	req := httptest.NewRequest("POST", "/lists", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCreateTaskList_NameRequired(t *testing.T) {
	app := fiber.New()
	h := &TaskListHandler{service: &mockTaskListService{}}
	app.Post("/lists", h.CreateTaskList)
	body := `{"description":"desc"}`
	req := httptest.NewRequest("POST", "/lists", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCreateTaskList_ServiceError(t *testing.T) {
	app := fiber.New()
	h := &TaskListHandler{service: &mockTaskListService{
		CreateFn: func(name, description string) (*domain.TaskList, error) {
			return nil, errors.New("fail")
		},
	}}
	app.Post("/lists", h.CreateTaskList)
	body := `{"name":"Lista","description":"desc"}`
	req := httptest.NewRequest("POST", "/lists", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}
