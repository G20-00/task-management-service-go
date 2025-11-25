package task

import "github.com/G20-00/task-management-service-go/internal/domain"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(task *domain.Task) error {
	return s.repo.Create(task)
}

func (s *Service) GetAll() ([]*domain.Task, error) {
	return s.repo.GetAll()
}

func (s *Service) GetByID(id string) (*domain.Task, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(task *domain.Task) error {
	return s.repo.Update(task)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
