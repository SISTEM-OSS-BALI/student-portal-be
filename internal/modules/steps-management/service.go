package steps

import "github.com/username/gin-gorm-api/internal/schema"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(input CreateDTO) (schema.StepsManagement, error) {
	step := schema.StepsManagement{Label: input.Label}
	if err := s.repo.Create(&step); err != nil {
		return schema.StepsManagement{}, err
	}
	if len(input.ChildIDs) > 0 {
		if err := s.repo.ReplaceChildren(step.ID, input.ChildIDs); err != nil {
			return schema.StepsManagement{}, err
		}
	}
	return s.repo.GetByID(step.ID)
}

func (s *Service) List() ([]schema.StepsManagement, error) {
	return s.repo.List()
}

func (s *Service) GetByID(id string) (schema.StepsManagement, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(id string, input UpdateDTO) (schema.StepsManagement, error) {
	step, err := s.repo.GetByID(id)
	if err != nil {
		return schema.StepsManagement{}, err
	}

	if input.Label != nil {
		step.Label = *input.Label
	}

	if err := s.repo.Update(&step); err != nil {
		return schema.StepsManagement{}, err
	}
	if input.ChildIDs != nil {
		if err := s.repo.ReplaceChildren(step.ID, *input.ChildIDs); err != nil {
			return schema.StepsManagement{}, err
		}
	}
	return s.repo.GetByID(step.ID)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
