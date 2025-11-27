package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/G20-00/task-management-service-go/internal/domain"
)

func TestPostgresTaskRepository_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	r := NewPostgresTaskRepository(db)
	mock.ExpectExec("INSERT INTO tasks").WillReturnError(errors.New("fail"))
	task := &domain.Task{ID: "1", ListID: "1", Title: "t", Status: "pending", Priority: "medium"}
	err = r.Create(task)
	if err == nil {
		t.Error("esperado error en Create")
	}
}

func TestPostgresTaskRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	r := NewPostgresTaskRepository(db)
	mock.ExpectQuery("SELECT id, list_id, title, description, status, priority, created_at, updated_at FROM tasks WHERE id = \\$1").
		WithArgs("no-task").WillReturnError(sql.ErrNoRows)
	_, err = r.GetByID("no-task")
	if err == nil {
		t.Error("esperado error por no encontrado")
	}
}
