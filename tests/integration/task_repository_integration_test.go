package integration_test

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/G20-00/task-management-service-go/internal/domain"
	"github.com/G20-00/task-management-service-go/internal/infrastructure/repository"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load(); err != nil {
		println("Warning: .env file not loaded:", err.Error())
	}
	os.Exit(m.Run())
}

func cleanupTasks(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DELETE FROM tasks")
	if err != nil {
		t.Fatalf("Failed to cleanup tasks: %v", err)
	}
}

func TestPostgresTaskRepository_Create_Integration(t *testing.T) {
	db := getTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Failed to close db: %v", err)
		}
	}()

	cleanupTasks(t, db)

	listID := "test-list-id"
	_, err := db.Exec(`INSERT INTO task_lists (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING`, listID, "Test List", time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to create test list: %v", err)
	}

	repo := repository.NewPostgresTaskRepository(db)

	now := time.Now()
	task := &domain.Task{
		ID:          uuid.New().String(),
		ListID:      listID,
		Title:       "Integration Test Task",
		Description: "This is a test task",
		Status:      "pending",
		Priority:    "high",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = repo.Create(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM tasks WHERE id = $1", task.ID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query task: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 task, got %d", count)
	}

	// Cleanup
	cleanupTasks(t, db)
}

func TestPostgresTaskRepository_GetByID_Integration(t *testing.T) {
	db := getTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Failed to close db: %v", err)
		}
	}()

	cleanupTasks(t, db)

	repo := repository.NewPostgresTaskRepository(db)

	now := time.Now()
	taskID := uuid.New().String()
	_, err := db.Exec(`
		INSERT INTO tasks (id, list_id, title, description, status, priority, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, taskID, "test-list-id", "Test Task", "Test Description", "pending", "medium", now, now)

	if err != nil {
		t.Fatalf("Failed to insert test task: %v", err)
	}

	task, err := repo.GetByID(taskID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	if task.ID != taskID {
		t.Errorf("Expected task ID %s, got %s", taskID, task.ID)
	}

	if task.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%s'", task.Title)
	}

	if task.Status != "pending" {
		t.Errorf("Expected status 'pending', got '%s'", task.Status)
	}

	if task.Priority != "medium" {
		t.Errorf("Expected priority 'medium', got '%s'", task.Priority)
	}

	// Cleanup
	cleanupTasks(t, db)
}
