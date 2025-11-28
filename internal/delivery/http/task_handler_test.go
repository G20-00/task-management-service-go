package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/G20-00/task-management-service-go/internal/domain"
)

func TestGetTasks_EmptyList(t *testing.T) {
	app := fiber.New()
	mockService := &mockTaskService{
		GetAllFn: func() ([]*domain.Task, error) { return []*domain.Task{}, nil },
	}
	h := NewTaskHandler(mockService)
	app.Get("/tasks", h.GetTasks)
	req := httptest.NewRequest("GET", "/tasks", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetTasks_Filtered_Empty(t *testing.T) {
	app := fiber.New()
	mockService := &mockTaskService{
		GetByFiltersFn: func(status, priority string) ([]*domain.Task, error) { return []*domain.Task{}, nil },
	}
	h := NewTaskHandler(mockService)
	app.Get("/tasks", h.GetTasks)
	req := httptest.NewRequest("GET", "/tasks?status=done&priority=low", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetTasks_Filtered_OnlyStatus(t *testing.T) {
	app := fiber.New()
	mockService := &mockTaskService{
		GetByFiltersFn: func(status, priority string) ([]*domain.Task, error) {
			if status == "pending" && priority == "" {
				return []*domain.Task{{ID: "1", Status: status, CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil
			}
			return nil, errors.New("unexpected params")
		},
	}
	h := NewTaskHandler(mockService)
	app.Get("/tasks", h.GetTasks)
	req := httptest.NewRequest("GET", "/tasks?status=pending", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetTasks_Filtered_OnlyPriority(t *testing.T) {
	app := fiber.New()
	mockService := &mockTaskService{
		GetByFiltersFn: func(status, priority string) ([]*domain.Task, error) {
			if status == "" && priority == "high" {
				return []*domain.Task{{ID: "1", Priority: priority, CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil
			}
			return nil, errors.New("unexpected params")
		},
	}
	h := NewTaskHandler(mockService)
	app.Get("/tasks", h.GetTasks)
	req := httptest.NewRequest("GET", "/tasks?priority=high", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetTask_ServiceErrorOther(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{
		GetByIDFn: func(id string) (*domain.Task, error) { return nil, errors.New("db error") },
	})
	app.Get("/tasks/:id", h.GetTask)
	req := httptest.NewRequest("GET", "/tasks/1", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
}

func TestUpdateTask_ServiceErrorOther(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{
		UpdateFn: func(id, listID, title, description, status, priority string) (*domain.Task, error) {
			return nil, errors.New("db error")
		},
	})
	app.Put("/tasks/:taskId", h.UpdateTask)
	req := httptest.NewRequest("PUT", "/tasks/1", strings.NewReader(`{"title":"T","status":"pending","priority":"high"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
}

func TestDeleteTask_ServiceErrorOther(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{
		DeleteFn: func(id string) error { return errors.New("db error") },
	})
	app.Delete("/tasks/:id", h.DeleteTask)
	req := httptest.NewRequest("DELETE", "/tasks/1", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
}

func TestCreateTask_InvalidBody(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{})
	app.Post("/tasks", h.CreateTask)
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader("{malformed_json"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestUpdateTask_InvalidBody(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{})
	app.Put("/tasks/:id", h.UpdateTask)
	req := httptest.NewRequest("PUT", "/tasks/1", strings.NewReader("{malformed_json"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
}

func TestGetTask_EmptyID(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{})
	app.Get("/tasks/:id", h.GetTask)
	req := httptest.NewRequest("GET", "/tasks/", http.NoBody)
	req.RequestURI = "/tasks/"
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest && resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("expected status 400 or 404, got %d", resp.StatusCode)
	}
}

func TestDeleteTask_EmptyID(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{})
	app.Delete("/tasks/:id", h.DeleteTask)
	req := httptest.NewRequest("DELETE", "/tasks/", http.NoBody)
	req.RequestURI = "/tasks/"
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest && resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("expected status 400 or 404, got %d", resp.StatusCode)
	}
}

// mockTaskService implements TaskService for testing
type mockTaskService struct {
	CreateFn       func(listID, title, description, priority string) (*domain.Task, error)
	GetAllFn       func() ([]*domain.Task, error)
	GetByFiltersFn func(status, priority string) ([]*domain.Task, error)
	GetByIDFn      func(id string) (*domain.Task, error)
	UpdateFn       func(id, listID, title, description, status, priority string) (*domain.Task, error)
	DeleteFn       func(id string) error
}

func (m *mockTaskService) Create(listID, title, description, priority string) (*domain.Task, error) {
	if m.CreateFn != nil {
		return m.CreateFn(listID, title, description, priority)
	}
	return nil, nil
}
func (m *mockTaskService) GetAll() ([]*domain.Task, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn()
	}
	return nil, nil
}
func (m *mockTaskService) GetByFilters(status, priority string) ([]*domain.Task, error) {
	if m.GetByFiltersFn != nil {
		return m.GetByFiltersFn(status, priority)
	}
	return nil, nil
}
func (m *mockTaskService) GetByID(id string) (*domain.Task, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return nil, nil
}
func (m *mockTaskService) Update(id, listID, title, description, status, priority string) (*domain.Task, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(id, listID, title, description, status, priority)
	}
	return nil, nil
}
func (m *mockTaskService) Delete(id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

func TestGetTasks_Filtered_Success(t *testing.T) {
	app := fiber.New()
	mockService := &mockTaskService{
		GetByFiltersFn: func(status, priority string) ([]*domain.Task, error) {
			return []*domain.Task{{ID: "2", Status: status, Priority: priority, CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil
		},
	}
	h := NewTaskHandler(mockService)
	app.Get("/tasks", h.GetTasks)
	req := httptest.NewRequest("GET", "/tasks?status=pending&priority=high", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetTasks_Filtered_Error(t *testing.T) {
	app := fiber.New()
	mockService := &mockTaskService{
		GetByFiltersFn: func(status, priority string) ([]*domain.Task, error) {
			return nil, errors.New("fail")
		},
	}
	h := NewTaskHandler(mockService)
	app.Get("/tasks", h.GetTasks)
	req := httptest.NewRequest("GET", "/tasks?status=pending&priority=high", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
}

func TestUpdateTask_TitleRequired(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{})
	app.Put("/tasks/:id", h.UpdateTask)
	req := httptest.NewRequest("PUT", "/tasks/1", strings.NewReader(`{"status":"pending","priority":"high"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestUpdateTask_StatusRequired(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{})
	app.Put("/tasks/:taskId", h.UpdateTask)
	req := httptest.NewRequest("PUT", "/tasks/1", strings.NewReader(`{"title":"T","priority":"high"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestUpdateTask_PriorityRequired(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{})
	app.Put("/tasks/:taskId", h.UpdateTask)
	req := httptest.NewRequest("PUT", "/tasks/1", strings.NewReader(`{"title":"T","status":"pending"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestUpdateTask_NotFound(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{
		UpdateFn: func(id, listID, title, description, status, priority string) (*domain.Task, error) {
			return nil, errors.New("task not found")
		},
	})
	app.Put("/tasks/:taskId", h.UpdateTask)
	req := httptest.NewRequest("PUT", "/tasks/1", strings.NewReader(`{"title":"T","status":"pending","priority":"high"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestDeleteTask_NotFound_Error(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{
		DeleteFn: func(id string) error { return errors.New("task not found") },
	})
	app.Delete("/tasks/:id", h.DeleteTask)
	req := httptest.NewRequest("DELETE", "/tasks/1", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestCreateTask_Success(t *testing.T) {
	app := fiber.New()
	mockService := &mockTaskService{
		CreateFn: func(listID, title, description, priority string) (*domain.Task, error) {
			return &domain.Task{ID: "1", ListID: listID, Title: title, Description: description, Priority: priority, Status: "pending", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
		},
	}
	h := NewTaskHandler(mockService)
	app.Post("/tasks", h.CreateTask)
	body := `{"list_id":"l1","title":"T","description":"D","priority":"high"}`
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
}

func TestCreateTask_BadRequest(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{})
	app.Post("/tasks", h.CreateTask)
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestGetTasks_Success(t *testing.T) {
	app := fiber.New()
	mockService := &mockTaskService{
		GetAllFn: func() ([]*domain.Task, error) {
			return []*domain.Task{{ID: "1", Title: "T", CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil
		},
		GetByFiltersFn: func(status, priority string) ([]*domain.Task, error) {
			return []*domain.Task{{ID: "2", Status: status, Priority: priority, CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil
		},
	}
	h := NewTaskHandler(mockService)
	app.Get("/tasks", h.GetTasks)
	req := httptest.NewRequest("GET", "/tasks", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	app := fiber.New()
	mockService := &mockTaskService{
		GetByIDFn: func(id string) (*domain.Task, error) { return nil, errors.New("task not found") },
	}
	h := NewTaskHandler(mockService)
	app.Get("/tasks/:id", h.GetTask)
	req := httptest.NewRequest("GET", "/tasks/123", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestUpdateTask_BadRequest(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{})
	app.Put("/tasks/:id", h.UpdateTask)
	req := httptest.NewRequest("PUT", "/tasks/1", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestDeleteTask_Success(t *testing.T) {
	app := fiber.New()
	mockService := &mockTaskService{
		DeleteFn: func(id string) error { return nil },
	}
	h := NewTaskHandler(mockService)
	app.Delete("/tasks/:id", h.DeleteTask)
	req := httptest.NewRequest("DELETE", "/tasks/1", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestCreateTask_TitleRequired(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{})
	app.Post("/tasks", h.CreateTask)
	body := `{"list_id":"l1","description":"D","priority":"high"}`
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateTask_ServiceError(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{
		CreateFn: func(listID, title, description, priority string) (*domain.Task, error) {
			return nil, errors.New("fail")
		},
	})
	app.Post("/tasks", h.CreateTask)
	body := `{"list_id":"l1","title":"T","description":"D","priority":"high"}`
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
}

func TestGetTasks_ServiceError(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{
		GetAllFn: func() ([]*domain.Task, error) { return nil, errors.New("fail") },
	})
	app.Get("/tasks", h.GetTasks)
	req := httptest.NewRequest("GET", "/tasks", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
}

func TestGetTask_BadRequest(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{})
	app.Get("/tasks/:id", h.GetTask)
	req := httptest.NewRequest("GET", "/tasks/", http.NoBody)
	req.RequestURI = "/tasks/"
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest && resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("expected status 400 or 404, got %d", resp.StatusCode)
	}
}

func TestGetTask_ServiceError(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{
		GetByIDFn: func(id string) (*domain.Task, error) { return nil, errors.New("fail") },
	})
	app.Get("/tasks/:id", h.GetTask)
	req := httptest.NewRequest("GET", "/tasks/1", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
}

func TestGetTask_Success(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{
		GetByIDFn: func(id string) (*domain.Task, error) {
			return &domain.Task{ID: id, Title: "T", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
		},
	})
	app.Get("/tasks/:id", h.GetTask)
	req := httptest.NewRequest("GET", "/tasks/1", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestUpdateTask_MissingID(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{})
	app.Put("/tasks/:taskId", h.UpdateTask)
	req := httptest.NewRequest("PUT", "/tasks/", strings.NewReader("{}"))
	req.RequestURI = "/tasks/"
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest && resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("expected status 400 or 404, got %d", resp.StatusCode)
	}
}

func TestUpdateTask_MissingFields(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{})
	app.Put("/tasks/:taskId", h.UpdateTask)
	req := httptest.NewRequest("PUT", "/tasks/1", strings.NewReader(`{"title":"","status":"","priority":""}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestUpdateTask_ServiceError(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{
		UpdateFn: func(id, listID, title, description, status, priority string) (*domain.Task, error) {
			return nil, errors.New("fail")
		},
	})
	app.Put("/tasks/:taskId", h.UpdateTask)
	req := httptest.NewRequest("PUT", "/tasks/1", strings.NewReader(`{"title":"T","status":"pending","priority":"high"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
}

func TestUpdateTask_Success(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{
		UpdateFn: func(id, listID, title, description, status, priority string) (*domain.Task, error) {
			return &domain.Task{ID: id, Title: title, Status: status, Priority: priority, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
		},
	})
	app.Put("/tasks/:taskId", h.UpdateTask)
	req := httptest.NewRequest("PUT", "/tasks/1", strings.NewReader(`{"title":"T","status":"pending","priority":"high"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestDeleteTask_MissingID(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{})
	app.Delete("/tasks/:id", h.DeleteTask)
	req := httptest.NewRequest("DELETE", "/tasks/", http.NoBody)
	req.RequestURI = "/tasks/"
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest && resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("expected status 400 or 404, got %d", resp.StatusCode)
	}
}

func TestDeleteTask_ServiceError(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{
		DeleteFn: func(id string) error { return errors.New("fail") },
	})
	app.Delete("/tasks/:taskId", h.DeleteTask)
	req := httptest.NewRequest("DELETE", "/tasks/1", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	app := fiber.New()
	h := NewTaskHandler(&mockTaskService{
		DeleteFn: func(id string) error { return errors.New("task not found") },
	})
	app.Delete("/tasks/:taskId", h.DeleteTask)
	req := httptest.NewRequest("DELETE", "/tasks/1", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}
