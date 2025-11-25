package repository

import (
	"database/sql"
	"errors"

	"github.com/G20-00/task-management-service-go/internal/domain"
)

type PostgresTaskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{
		db: db,
	}
}

func (r *PostgresTaskRepository) Create(task *domain.Task) error {
	query := `INSERT INTO tasks (id, list_id, title, description, status, priority, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(query, task.ID, task.ListID, task.Title, task.Description, task.Status, task.Priority, task.CreatedAt, task.UpdatedAt)
	return err
}

func (r *PostgresTaskRepository) GetAll() ([]*domain.Task, error) {
	query := `SELECT id, list_id, title, description, status, priority, created_at, updated_at 
	          FROM tasks ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (r *PostgresTaskRepository) Update(task *domain.Task) error {
	query := `UPDATE tasks SET list_id = $2, title = $3, description = $4, status = $5, priority = $6, updated_at = $7 
	          WHERE id = $1`

	result, err := r.db.Exec(query, task.ID, task.ListID, task.Title, task.Description, task.Status, task.Priority, task.UpdatedAt)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (r *PostgresTaskRepository) Delete(id string) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("task not found")
	}

	return nil
}
