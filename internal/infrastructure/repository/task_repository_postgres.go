// Package repository provides database repository implementations.
package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/G20-00/task-management-service-go/internal/domain"
)

// PostgresTaskRepository is a PostgreSQL implementation of task repository.
type PostgresTaskRepository struct {
	db *sql.DB
}

// NewPostgresTaskRepository creates a new PostgresTaskRepository instance.
func NewPostgresTaskRepository(db *sql.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{
		db: db,
	}
}

// Create inserts a new task into the database.
func (r *PostgresTaskRepository) Create(task *domain.Task) error {
	query := `INSERT INTO tasks (id, list_id, title, description, status, priority, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(query, task.ID, task.ListID, task.Title, task.Description, task.Status, task.Priority, task.CreatedAt, task.UpdatedAt)
	return err
}

// GetAll retrieves all tasks from the database.
func (r *PostgresTaskRepository) GetAll() ([]*domain.Task, error) {
	query := `SELECT id, list_id, title, description, status, priority, created_at, updated_at 
	          FROM tasks ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close() //nolint:errcheck,gocritic
	}()

	tasks := []*domain.Task{}
	for rows.Next() {
		task := &domain.Task{}
		if err := rows.Scan(&task.ID, &task.ListID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

// GetByID retrieves a single task by ID.
func (r *PostgresTaskRepository) GetByID(id string) (*domain.Task, error) {
	query := `SELECT id, list_id, title, description, status, priority, created_at, updated_at 
	          FROM tasks WHERE id = $1`

	task := &domain.Task{}
	err := r.db.QueryRow(query, id).Scan(&task.ID, &task.ListID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("task not found")
	}
	if err != nil {
		return nil, err
	}

	return task, nil
}

// Update modifies an existing task in the database.
func (r *PostgresTaskRepository) Update(task *domain.Task) error {
	query := `UPDATE tasks SET list_id = $2, title = $3, description = $4, status = $5, priority = $6, updated_at = $7
	          WHERE id = $1`

	result, err := r.db.Exec(query, task.ID, task.ListID, task.Title, task.Description, task.Status, task.Priority, task.UpdatedAt)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("task not found")
	}

	return nil
}

// Delete removes a task from the database.
func (r *PostgresTaskRepository) Delete(id string) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("task not found")
	}

	return nil
}

// GetByFilters retrieves tasks filtered by status and/or priority.
func (r *PostgresTaskRepository) GetByFilters(status, priority string) ([]*domain.Task, error) {
	query := `SELECT id, list_id, title, description, status, priority, created_at, updated_at 
	          FROM tasks WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	if status != "" {
		query += ` AND status = $` + fmt.Sprintf("%d", argCount)
		args = append(args, status)
		argCount++
	}

	if priority != "" {
		query += ` AND priority = $` + fmt.Sprintf("%d", argCount)
		args = append(args, priority)
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close() //nolint:errcheck,gocritic
	}()

	tasks := []*domain.Task{}
	for rows.Next() {
		task := &domain.Task{}
		if err := rows.Scan(&task.ID, &task.ListID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

// CountByListIDAndStatus counts tasks by list ID and status.
func (r *PostgresTaskRepository) CountByListIDAndStatus(listID, status string) (int, error) {
	query := `SELECT COUNT(*) FROM tasks WHERE list_id = $1 AND status = $2`

	var count int
	err := r.db.QueryRow(query, listID, status).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
