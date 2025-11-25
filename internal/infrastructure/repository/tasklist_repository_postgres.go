package repository

import (
	"database/sql"
	"errors"

	"github.com/G20-00/task-management-service-go/internal/domain"
)

type PostgresTaskListRepository struct {
	db *sql.DB
}

func NewPostgresTaskListRepository(db *sql.DB) *PostgresTaskListRepository {
	return &PostgresTaskListRepository{
		db: db,
	}
}

func (r *PostgresTaskListRepository) Create(list *domain.TaskList) error {
	query := `INSERT INTO task_lists (id, name, description, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query, list.ID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt)
	return err
}

func (r *PostgresTaskListRepository) GetAll() ([]*domain.TaskList, error) {
	query := `SELECT id, name, description, created_at, updated_at 
	          FROM task_lists ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lists := []*domain.TaskList{}
	for rows.Next() {
		list := &domain.TaskList{}
		if err := rows.Scan(&list.ID, &list.Name, &list.Description, &list.CreatedAt, &list.UpdatedAt); err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}

	return lists, rows.Err()
}

func (r *PostgresTaskListRepository) GetByID(id string) (*domain.TaskList, error) {
	query := `SELECT id, name, description, created_at, updated_at 
	          FROM task_lists WHERE id = $1`

	list := &domain.TaskList{}
	err := r.db.QueryRow(query, id).Scan(&list.ID, &list.Name, &list.Description, &list.CreatedAt, &list.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("task list not found")
	}
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (r *PostgresTaskListRepository) Update(list *domain.TaskList) error {
	query := `UPDATE task_lists SET name = $2, description = $3, updated_at = $4 
	          WHERE id = $1`

	result, err := r.db.Exec(query, list.ID, list.Name, list.Description, list.UpdatedAt)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("task list not found")
	}

	return nil
}

func (r *PostgresTaskListRepository) Delete(id string) error {
	query := `DELETE FROM task_lists WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("task list not found")
	}

	return nil
}
