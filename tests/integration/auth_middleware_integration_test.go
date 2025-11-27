package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpdelivery "github.com/G20-00/task-management-service-go/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

func TestProtectedEndpoint_WithoutJWT_ShouldFail(t *testing.T) {
	app := fiber.New()
	app.Get("/api/tasks", httpdelivery.JWTMiddleware, func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/tasks", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error al hacer request: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Esperado 401, obtuve %d", resp.StatusCode)
	}
}

func TestProtectedEndpoint_WithJWT_ShouldPass(t *testing.T) {
	app := fiber.New()
	app.Get("/api/tasks", httpdelivery.JWTMiddleware, func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	token, err := httpdelivery.GenerateJWT("test-user")
	if err != nil {
		t.Fatalf("No se pudo generar token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/tasks", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error al hacer request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Esperado 200, obtuve %d", resp.StatusCode)
	}
}
