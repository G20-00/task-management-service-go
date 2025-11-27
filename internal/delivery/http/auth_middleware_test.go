package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestJWTMiddleware_MissingHeader(t *testing.T) {
	app := fiber.New()
	app.Use(JWTMiddleware)
	app.Get("/protected", func(c *fiber.Ctx) error { return c.SendStatus(200) })
	resp, err := app.Test(httptest.NewRequest("GET", "/protected", http.NoBody))
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestJWTMiddleware_InvalidHeader(t *testing.T) {
	app := fiber.New()
	app.Use(JWTMiddleware)
	app.Get("/protected", func(c *fiber.Ctx) error { return c.SendStatus(200) })
	req := httptest.NewRequest("GET", "/protected", http.NoBody)
	req.Header.Set("Authorization", "InvalidToken")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	app := fiber.New()
	app.Use(JWTMiddleware)
	app.Get("/protected", func(c *fiber.Ctx) error { return c.SendStatus(200) })
	req := httptest.NewRequest("GET", "/protected", http.NoBody)
	req.Header.Set("Authorization", "Bearer invalid.jwt.token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestJWTMiddleware_ValidToken(t *testing.T) {
	token, err := GenerateJWT("user1")
	if err != nil {
		t.Fatalf("error generando JWT: %v", err)
	}
	app := fiber.New()
	app.Use(JWTMiddleware)
	app.Get("/protected", func(c *fiber.Ctx) error { return c.SendStatus(200) })
	req := httptest.NewRequest("GET", "/protected", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}
