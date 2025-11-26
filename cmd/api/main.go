// Package main is the entry point for the task management service application.
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/G20-00/task-management-service-go/internal/delivery/http"
	"github.com/G20-00/task-management-service-go/internal/infrastructure/db"
	"github.com/G20-00/task-management-service-go/internal/infrastructure/repository"
	"github.com/G20-00/task-management-service-go/internal/usecase/task"
	"github.com/G20-00/task-management-service-go/internal/usecase/tasklist"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	database, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if closeErr := database.Close(); closeErr != nil {
			log.Printf("Failed to close database connection: %v", closeErr)
		}
	}()

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	taskRepo := repository.NewPostgresTaskRepository(database)
	taskService := task.NewService(taskRepo)
	taskHandler := http.NewTaskHandler(taskService)

	taskListRepo := repository.NewPostgresTaskListRepository(database)
	taskListService := tasklist.NewService(taskListRepo)
	taskListHandler := http.NewTaskListHandler(taskListService, taskService)

	http.RegisterRoutes(app, taskHandler, taskListHandler)

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
