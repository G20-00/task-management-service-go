// Package http provides HTTP handlers and routing for the task management API.
package http

import "github.com/gofiber/fiber/v2"

// RegisterRoutes configures all API routes for tasks and task lists.
func RegisterRoutes(app *fiber.App, taskHandler *TaskHandler, taskListHandler *TaskListHandler) {
	api := app.Group("/api")

	api.Post("/login", AuthHandler)

	tasks := api.Group("/tasks", JWTMiddleware)
	tasks.Post("/", taskHandler.CreateTask)
	tasks.Get("/", taskHandler.GetTasks)
	tasks.Get(":id", taskHandler.GetTask)
	tasks.Put(":id", taskHandler.UpdateTask)
	tasks.Patch(":id", taskHandler.UpdateTask) // Permitir PATCH directo
	tasks.Delete(":id", taskHandler.DeleteTask)

	// Rutas anidadas para compatibilidad con integración
	lists := api.Group("/lists", JWTMiddleware)
	lists.Post("/", taskListHandler.CreateTaskList)
	lists.Get("/", taskListHandler.GetTaskLists)
	lists.Get(":id", taskListHandler.GetTaskList)
	lists.Put(":id", taskListHandler.UpdateTaskList)
	lists.Delete(":id", taskListHandler.DeleteTaskList)

	// Tareas bajo listas (para integración)
	lists.Post(":id/tasks", taskHandler.CreateTask)
	lists.Get(":id/tasks/:taskId", taskHandler.GetTask)
	lists.Patch(":id/tasks/:taskId/state", taskHandler.UpdateTask)
	lists.Delete(":id/tasks/:taskId", taskHandler.DeleteTask)

	// (Ya declarado arriba)
}
