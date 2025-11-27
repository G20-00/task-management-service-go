package repository

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/G20-00/task-management-service-go/internal/domain"
)

var db *sql.DB
var mock sqlmock.Sqlmock

func TestPostgresTaskRepository_Create_Success(t *testing.T) {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	r := NewPostgresTaskRepository(db)
	mock.ExpectExec("INSERT INTO tasks").WillReturnResult(sqlmock.NewResult(1, 1))
	task := &domain.Task{ID: "1", ListID: "1", Title: "t", Status: "pending", Priority: "medium", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err = r.Create(task)
	if err != nil {
		t.Errorf("no se esperaba error en Create: %v", err)
	}
}

func TestPostgresTaskRepository_Update_NotFound(t *testing.T) {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	r := NewPostgresTaskRepository(db)
	mock.ExpectExec("UPDATE tasks SET").WillReturnResult(sqlmock.NewResult(0, 0))
	task := &domain.Task{ID: "1", ListID: "1", Title: "t", Status: "pending", Priority: "medium", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err = r.Update(task)
	if err == nil {
		t.Error("esperado error por tarea no encontrada en Update")
	}
}

func TestPostgresTaskRepository_Delete_NotFound(t *testing.T) {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	r := NewPostgresTaskRepository(db)
	mock.ExpectExec("DELETE FROM tasks").WillReturnResult(sqlmock.NewResult(0, 0))
	err = r.Delete("no-task")
	if err == nil {
		t.Error("esperado error por tarea no encontrada en Delete")
	}
}

func TestPostgresTaskRepository_GetAll_Success(t *testing.T) {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	r := NewPostgresTaskRepository(db)

	rows := sqlmock.NewRows([]string{"id", "list_id", "title", "description", "status", "priority", "created_at", "updated_at"}).
		AddRow("1", "1", "t", "desc", "pending", "medium", time.Now(), time.Now())
	mock.ExpectQuery("SELECT id, list_id, title, description, status, priority, created_at, updated_at FROM tasks ORDER BY created_at DESC").WillReturnRows(rows)

	tasks, err := r.GetAll()
	if err != nil {
		t.Errorf("no se esperaba error en GetAll: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("esperado 1 tarea, obtuve %d", len(tasks))
	}
}

func TestPostgresTaskRepository_GetByID_Success(t *testing.T) {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("error creando sqlmock: %v", err)
	}
	r := NewPostgresTaskRepository(db)

	row := sqlmock.NewRows([]string{"id", "list_id", "title", "description", "status", "priority", "created_at", "updated_at"}).
		AddRow("1", "1", "t", "desc", "pending", "medium", time.Now(), time.Now())
	mock.ExpectQuery("SELECT id, list_id, title, description, status, priority, created_at, updated_at FROM tasks WHERE id = \\$1").WithArgs("1").WillReturnRows(row)

	task, err := r.GetByID("1")
	if err != nil {
		t.Errorf("no se esperaba error en GetByID: %v", err)
	}
	if task == nil || task.ID != "1" {
		t.Errorf("esperado tarea con ID '1', obtuve %+v", task)
	}
}

func TestPostgresTaskRepository_Update_Success(t *testing.T) {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("error creando sqlmock: %v", err)
	}
	r := NewPostgresTaskRepository(db)

	mock.ExpectExec("UPDATE tasks SET").WillReturnResult(sqlmock.NewResult(0, 1))
	task := &domain.Task{ID: "1", ListID: "1", Title: "t", Status: "pending", Priority: "medium", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err = r.Update(task)
	if err != nil {
		t.Errorf("no se esperaba error en Update: %v", err)
	}
}

func TestPostgresTaskRepository_Update_ErrorRowsAffected(t *testing.T) {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("error creando sqlmock: %v", err)
	}
	r := NewPostgresTaskRepository(db)

	mock.ExpectExec("UPDATE tasks SET").WillReturnError(errors.New("fail"))
	task := &domain.Task{ID: "1", ListID: "1", Title: "t", Status: "pending", Priority: "medium", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err = r.Update(task)
	if err == nil {
		t.Error("esperado error en Update por error de base de datos")
	}
}

func TestPostgresTaskRepository_Delete_Success(t *testing.T) {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("error creando sqlmock: %v", err)
	}
	r := NewPostgresTaskRepository(db)

	mock.ExpectExec("DELETE FROM tasks").WillReturnResult(sqlmock.NewResult(0, 1))
	err = r.Delete("1")
	if err != nil {
		t.Errorf("no se esperaba error en Delete: %v", err)
	}
}

func TestPostgresTaskRepository_GetByFilters_StatusAndPriority(t *testing.T) {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("error creando sqlmock: %v", err)
	}
	r := NewPostgresTaskRepository(db)

	rows := sqlmock.NewRows([]string{"id", "list_id", "title", "description", "status", "priority", "created_at", "updated_at"}).
		AddRow("1", "1", "t", "desc", "pending", "medium", time.Now(), time.Now())
	mock.ExpectQuery(`(?s)SELECT id, list_id, title, description, status, priority, created_at, updated_at FROM tasks WHERE 1=1 AND status = \$1 AND priority = \$2 ORDER BY created_at DESC`).WillReturnRows(rows)

	tasks, err := r.GetByFilters("pending", "medium")
	if err != nil {
		t.Errorf("no se esperaba error en GetByFilters: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("esperado 1 tarea, obtuve %d", len(tasks))
	}
}

func TestPostgresTaskRepository_CountByListIDAndStatus(t *testing.T) {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("error creando sqlmock: %v", err)
	}
	r := NewPostgresTaskRepository(db)

	expectedSQL := `SELECT id, list_id, title, description, status, priority, created_at, updated_at 
		 FROM tasks WHERE 1=1 AND status = \$1 AND priority = \$2 ORDER BY created_at DESC`

	rows := sqlmock.NewRows([]string{
		"id", "list_id", "title", "description", "status", "priority", "created_at", "updated_at",
	}).AddRow(
		"1", "1", "Task A", "Description A", "pending", "high", time.Now(), time.Now(),
	)

	mock.ExpectQuery(expectedSQL).WithArgs("pending", "high").WillReturnRows(rows)

	tasks, err := r.GetByFilters("pending", "high")
	if err != nil {
		t.Errorf("no se esperaba error en GetByFilters: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("esperado 1 tarea, obtuve %d", len(tasks))
	}
}
