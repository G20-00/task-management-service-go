package http

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRegisterRoutes(t *testing.T) {
	app := fiber.New()
	RegisterRoutes(app, nil, nil)

}
