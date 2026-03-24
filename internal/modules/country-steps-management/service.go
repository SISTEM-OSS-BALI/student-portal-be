package countrysteps

import "github.com/username/gin-gorm-api/internal/schema"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(input CreateDTO) (schema.CountryStepsManagement, error) {
	item := schema.CountryStepsManagement{
		CountryID: input.CountryID,
		StepID:    input.StepID,
	}
	if err := s.repo.Create(&item); err != nil {
		return schema.CountryStepsManagement{}, err
	}
	return item, nil
}

func (s *Service) List() ([]schema.CountryStepsManagement, error) {
	return s.repo.List()
}

func (s *Service) ListWithCountryAndStep() ([]schema.CountryStepsManagement, error) {
	return s.repo.ListWithCountryAndStep()
}

func (s *Service) GetByID(id string) (schema.CountryStepsManagement, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(id string, input UpdateDTO) (schema.CountryStepsManagement, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return schema.CountryStepsManagement{}, err
	}

	if input.CountryID != nil {
		item.CountryID = *input.CountryID
	}
	if input.StepID != nil {
		item.StepID = *input.StepID
	}

	if err := s.repo.Update(&item); err != nil {
		return schema.CountryStepsManagement{}, err
	}
	return item, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
