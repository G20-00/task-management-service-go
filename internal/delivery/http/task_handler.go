package http

import "github.com/gofiber/fiber/v2"

type TaskHandler struct {
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	return nil
}

func (h *TaskHandler) GetTasks(c *fiber.Ctx) error {
	return nil
}

func (h *TaskHandler) GetTask(c *fiber.Ctx) error {
	return nil
}

func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	return nil
}

func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	return nil
}
