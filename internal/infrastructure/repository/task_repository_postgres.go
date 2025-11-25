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
	query := `INSERT INTO tasks (id, title, description, status, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query, task.ID, task.Title, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)
	return err
}

func (r *PostgresTaskRepository) GetAll() ([]*domain.Task, error) {
	query := `SELECT id, title, description, status, created_at, updated_at 
	          FROM tasks ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*domain.Task{}
	for rows.Next() {
		task := &domain.Task{}
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (r *PostgresTaskRepository) GetByID(id string) (*domain.Task, error) {
	query := `SELECT id, title, description, status, created_at, updated_at 
	          FROM tasks WHERE id = $1`

	task := &domain.Task{}
	err := r.db.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("task not found")
	}
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *PostgresTaskRepository) Update(task *domain.Task) error {
	query := `UPDATE tasks SET title = $2, description = $3, status = $4, updated_at = $5 
	          WHERE id = $1`

	result, err := r.db.Exec(query, task.ID, task.Title, task.Description, task.Status, task.UpdatedAt)
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
