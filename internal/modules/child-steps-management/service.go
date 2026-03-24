package childsteps

import "github.com/username/gin-gorm-api/internal/schema"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(input CreateDTO) (schema.ChildStepsManagement, error) {
	child := schema.ChildStepsManagement{Label: input.Label}
	if err := s.repo.Create(&child); err != nil {
		return schema.ChildStepsManagement{}, err
	}
	return child, nil
}

func (s *Service) List() ([]schema.ChildStepsManagement, error) {
	return s.repo.List()
}

func (s *Service) ListWithSteps() ([]schema.ChildStepsManagement, error) {
	return s.repo.ListWithSteps()
}

func (s *Service) GetByID(id string) (schema.ChildStepsManagement, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(id string, input UpdateDTO) (schema.ChildStepsManagement, error) {
	child, err := s.repo.GetByID(id)
	if err != nil {
		return schema.ChildStepsManagement{}, err
	}

	if input.Label != nil {
		child.Label = *input.Label
	}

	if err := s.repo.Update(&child); err != nil {
		return schema.ChildStepsManagement{}, err
	}
	return child, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
