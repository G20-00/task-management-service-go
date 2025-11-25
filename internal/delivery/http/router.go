package http

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, taskHandler *TaskHandler, taskListHandler *TaskListHandler) {
	api := app.Group("/api")

	tasks := api.Group("/tasks")
	tasks.Post("/", taskHandler.CreateTask)
	tasks.Get("/", taskHandler.GetTasks)
	tasks.Get("/:id", taskHandler.GetTask)
	tasks.Put("/:id", taskHandler.UpdateTask)
	tasks.Delete("/:id", taskHandler.DeleteTask)

	lists := api.Group("/lists")
	lists.Post("/", taskListHandler.CreateTaskList)
	lists.Get("/", taskListHandler.GetTaskLists)
	lists.Get("/:id", taskListHandler.GetTaskList)
	lists.Put("/:id", taskListHandler.UpdateTaskList)
	lists.Delete("/:id", taskListHandler.DeleteTaskList)
}
