package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/G20-00/task-management-service-go/internal/usecase/task"
	"github.com/G20-00/task-management-service-go/internal/usecase/tasklist"
)

// TaskListHandler handles HTTP requests for task list operations.
type TaskListHandler struct {
	service     *tasklist.Service
	taskService *task.Service
}

// NewTaskListHandler creates a new TaskListHandler instance.
func NewTaskListHandler(service *tasklist.Service, taskService *task.Service) *TaskListHandler {
	return &TaskListHandler{
		service:     service,
		taskService: taskService,
	}
}

// CreateTaskList handles the creation of a new task list.
func (h *TaskListHandler) CreateTaskList(c *fiber.Ctx) error {
	var req CreateTaskListRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required",
		})
	}

	list, err := h.service.Create(req.Name, req.Description)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := TaskListResponse{
		ID:                   list.ID,
		Name:                 list.Name,
		Description:          list.Description,
		CompletionPercentage: 0.0,
		CreatedAt:            list.CreatedAt,
		UpdatedAt:            list.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetTaskLists retrieves all task lists with their completion percentages.
func (h *TaskListHandler) GetTaskLists(c *fiber.Ctx) error {
	lists, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	responses := make([]TaskListResponse, len(lists))
	for i, list := range lists {
		percentage := h.calculateCompletionPercentage(list.ID)
		responses[i] = TaskListResponse{
			ID:                   list.ID,
			Name:                 list.Name,
			Description:          list.Description,
			CompletionPercentage: percentage,
			CreatedAt:            list.CreatedAt,
			UpdatedAt:            list.UpdatedAt,
		}
	}

	return c.JSON(responses)
}

// GetTaskList retrieves a single task list by ID with completion percentage.
func (h *TaskListHandler) GetTaskList(c *fiber.Ctx) error {
	id := c.Params("id")

	list, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := TaskListResponse{
		ID:                   list.ID,
		Name:                 list.Name,
		Description:          list.Description,
		CompletionPercentage: h.calculateCompletionPercentage(list.ID),
		CreatedAt:            list.CreatedAt,
		UpdatedAt:            list.UpdatedAt,
	}

	return c.JSON(response)
}

// UpdateTaskList updates an existing task list.
func (h *TaskListHandler) UpdateTaskList(c *fiber.Ctx) error {
	id := c.Params("id")
	var req UpdateTaskListRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	list, err := h.service.Update(id, req.Name, req.Description)
	if err != nil {
		if err.Error() == "task list not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := TaskListResponse{
		ID:          list.ID,
		Name:        list.Name,
		Description: list.Description,
		CreatedAt:   list.CreatedAt,
		UpdatedAt:   list.UpdatedAt,
	}

	return c.JSON(response)
}

// DeleteTaskList deletes a task list by ID.
func (h *TaskListHandler) DeleteTaskList(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.Delete(id); err != nil {
		if err.Error() == "task list not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TaskListHandler) calculateCompletionPercentage(listID string) float64 {
	tasks, err := h.taskService.GetByFilters("", "")
	if err != nil {
		return 0.0
	}

	var totalTasks, completedTasks int
	for _, t := range tasks {
		if t.ListID == listID {
			totalTasks++
			if t.Status == "completed" {
				completedTasks++
			}
		}
	}

	if totalTasks == 0 {
		return 0.0
	}

	return (float64(completedTasks) / float64(totalTasks)) * 100
}
