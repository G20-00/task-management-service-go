package unit

import (
	"errors"
	"testing"

	"net/http/httptest"
	"strings"

	"github.com/gofiber/fiber/v2"

	httpdelivery "github.com/G20-00/task-management-service-go/internal/delivery/http"
	"github.com/G20-00/task-management-service-go/internal/domain"
	taskusecase "github.com/G20-00/task-management-service-go/internal/usecase/task"
)

func TestTaskHandler_CreateTask_BadRequest(t *testing.T) {
	svc := taskusecase.NewService(&mockRepo{})
	h := httpdelivery.NewTaskHandler(svc)
	app := fiber.New()
	app.Post("/tasks", h.CreateTask)
	// Enviar body inválido
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("esperado 400, obtuve %d", resp.StatusCode)
	}
}

func TestTaskHandler_CreateTask_TitleRequired(t *testing.T) {
	svc := taskusecase.NewService(&mockRepo{})
	h := httpdelivery.NewTaskHandler(svc)
	app := fiber.New()
	app.Post("/tasks", h.CreateTask)
	body := `{"list_id": "1"}`
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("esperado 400 por título vacío, obtuve %d", resp.StatusCode)
	}
}

func TestTaskHandler_CreateTask_ServiceError(t *testing.T) {
	svc := taskusecase.NewService(&mockRepo{
		CreateFn: func(t *domain.Task) error { return errors.New("fail") },
	})
	h := httpdelivery.NewTaskHandler(svc)
	app := fiber.New()
	app.Post("/tasks", h.CreateTask)
	body := `{"list_id": "1", "title": "t"}`
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("esperado 500, obtuve %d", resp.StatusCode)
	}
}
