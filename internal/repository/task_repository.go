package repository

import (
	"errors"
	"sync"

	"github.com/G20-00/task-management-service-go/internal/domain"
	"github.com/G20-00/task-management-service-go/pkg/utils"
)

type InMemoryTaskRepository struct {
	tasks map[string]domain.Task
	mu    sync.RWMutex
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[string]domain.Task),
	}
}

func (r *InMemoryTaskRepository) Create(task *domain.Task) (err error) {
	defer utils.RecoverPanic("repository", "Create", &err)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; exists {
		return errors.New("task already exists")
	}

	r.tasks[task.ID] = *task
	return nil
}

func (r *InMemoryTaskRepository) GetAll() (tasks []*domain.Task, err error) {
	defer utils.RecoverPanic("repository", "GetAll", &err)

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		taskCopy := task
		result = append(result, &taskCopy)
	}

	return result, nil
}

func (r *InMemoryTaskRepository) GetByID(id string) (task *domain.Task, err error) {
	defer utils.RecoverPanic("repository", "GetByID", &err)

	r.mu.RLock()
	defer r.mu.RUnlock()

	t, exists := r.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}

	return &t, nil
}

func (r *InMemoryTaskRepository) Update(task *domain.Task) (err error) {
	defer utils.RecoverPanic("repository", "Update", &err)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return errors.New("task not found")
	}

	r.tasks[task.ID] = *task
	return nil
}

func (r *InMemoryTaskRepository) Delete(id string) (err error) {
	defer utils.RecoverPanic("repository", "Delete", &err)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return errors.New("task not found")
	}

	delete(r.tasks, id)
	return nil
}
