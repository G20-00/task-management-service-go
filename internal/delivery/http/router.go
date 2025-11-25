package http

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App) {
	handler := NewTaskHandler()

	api := app.Group("/api")
	tasks := api.Group("/tasks")

	tasks.Post("/", handler.CreateTask)
	tasks.Get("/", handler.GetTasks)
	tasks.Get("/:id", handler.GetTask)
	tasks.Put("/:id", handler.UpdateTask)
	tasks.Delete("/:id", handler.DeleteTask)
}
