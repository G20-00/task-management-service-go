package main

import (
	"github.com/G20-00/task-management-service-go/internal/delivery/http"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	http.RegisterRoutes(app)

	app.Listen(":8080")
}
