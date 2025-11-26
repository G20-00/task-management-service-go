// Package http provides HTTP handlers and routing for the task management API.
package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/G20-00/task-management-service-go/internal/domain"
	"github.com/G20-00/task-management-service-go/internal/usecase/task"
	"github.com/G20-00/task-management-service-go/pkg/logger"
)

// TaskHandler handles HTTP requests for task operations.
type TaskHandler struct {
	service *task.Service
}

// NewTaskHandler creates a new TaskHandler instance.
func NewTaskHandler(service *task.Service) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

// CreateTask handles the creation of a new task.
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var req CreateTaskRequest

	if err := c.BodyParser(&req); err != nil {
		logger.GetLogger().WithFields(map[string]interface{}{
			"layer":  "handler",
			"method": "CreateTask",
			"error":  err.Error(),
		}).Error("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	if req.Status == "" {
		req.Status = "pending"
	}

	if req.Priority == "" {
		req.Priority = "medium"
	}

	createdTask, err := h.service.Create(req.ListID, req.Title, req.Description, req.Priority)
	if err != nil {
		logger.GetLogger().WithFields(map[string]interface{}{
			"layer":  "handler",
			"method": "CreateTask",
			"error":  err.Error(),
		}).Error("Failed to create task")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create task",
		})
	}

	response := TaskResponse{
		ID:          createdTask.ID,
		ListID:      createdTask.ListID,
		Title:       createdTask.Title,
		Description: createdTask.Description,
		Status:      createdTask.Status,
		Priority:    createdTask.Priority,
		CreatedAt:   createdTask.CreatedAt,
		UpdatedAt:   createdTask.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetTasks retrieves all tasks or filters them by status and priority.
func (h *TaskHandler) GetTasks(c *fiber.Ctx) error {
	status := c.Query("status")
	priority := c.Query("priority")

	var tasks []*domain.Task
	var err error

	if status != "" || priority != "" {
		tasks, err = h.service.GetByFilters(status, priority)
	} else {
		tasks, err = h.service.GetAll()
	}

	if err != nil {
		logger.GetLogger().WithFields(map[string]interface{}{
			"layer":  "handler",
			"method": "GetTasks",
			"error":  err.Error(),
		}).Error("Failed to get tasks")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get tasks",
		})
	}

	responses := make([]TaskResponse, len(tasks))
	for i, t := range tasks {
		responses[i] = TaskResponse{
			ID:          t.ID,
			ListID:      t.ListID,
			Title:       t.Title,
			Description: t.Description,
			Status:      t.Status,
			Priority:    t.Priority,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		}
	}

	return c.Status(fiber.StatusOK).JSON(responses)
}

// GetTask retrieves a single task by ID.
func (h *TaskHandler) GetTask(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task ID is required",
		})
	}

	t, err := h.service.GetByID(id)
	if err != nil {
		if err.Error() == "task not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Task not found",
			})
		}

		logger.GetLogger().WithFields(map[string]interface{}{
			"layer":  "handler",
			"method": "GetTask",
			"taskID": id,
			"error":  err.Error(),
		}).Error("Failed to get task")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get task",
		})
	}

	response := TaskResponse{
		ID:          t.ID,
		ListID:      t.ListID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Priority:    t.Priority,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// UpdateTask updates an existing task.
func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task ID is required",
		})
	}

	var req UpdateTaskRequest

	if err := c.BodyParser(&req); err != nil {
		logger.GetLogger().WithFields(map[string]interface{}{
			"layer":  "handler",
			"method": "UpdateTask",
			"taskID": id,
			"error":  err.Error(),
		}).Error("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	if req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status is required",
		})
	}

	if req.Priority == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Priority is required",
		})
	}

	updatedTask, err := h.service.Update(id, req.ListID, req.Title, req.Description, req.Status, req.Priority)
	if err != nil {
		if err.Error() == "task not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Task not found",
			})
		}

		logger.GetLogger().WithFields(map[string]interface{}{
			"layer":  "handler",
			"method": "UpdateTask",
			"taskID": id,
			"error":  err.Error(),
		}).Error("Failed to update task")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update task",
		})
	}

	response := TaskResponse{
		ID:          updatedTask.ID,
		ListID:      updatedTask.ListID,
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		Status:      updatedTask.Status,
		Priority:    updatedTask.Priority,
		CreatedAt:   updatedTask.CreatedAt,
		UpdatedAt:   updatedTask.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteTask deletes a task by ID.
func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task ID is required",
		})
	}

	err := h.service.Delete(id)
	if err != nil {
		if err.Error() == "task not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Task not found",
			})
		}

		logger.GetLogger().WithFields(map[string]interface{}{
			"layer":  "handler",
			"method": "DeleteTask",
			"taskID": id,
			"error":  err.Error(),
		}).Error("Failed to delete task")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete task",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
