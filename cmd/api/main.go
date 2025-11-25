package main

import (
	"github.com/G20-00/task-management-service-go/internal/delivery/http"
	"github.com/G20-00/task-management-service-go/internal/repository"
	"github.com/G20-00/task-management-service-go/internal/usecase/task"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	taskRepo := repository.NewInMemoryTaskRepository()
	taskService := task.NewService(taskRepo)
	taskHandler := http.NewTaskHandler(taskService)

	http.RegisterRoutes(app, taskHandler)

	app.Listen(":8080")
}
