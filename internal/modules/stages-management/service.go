package stages

import "github.com/username/gin-gorm-api/internal/schema"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(input CreateDTO) (schema.StageManagement, error) {
	stage := schema.StageManagement{
		CountryID:  input.CountryID,
		DocumentID: input.DocumentID,
	}
	if err := s.repo.Create(&stage); err != nil {
		return schema.StageManagement{}, err
	}
	return stage, nil
}

func (s *Service) List() ([]schema.StageManagement, error) {
	return s.repo.List()
}

func (s *Service) ListWithCountryAndDocument() ([]schema.StageManagement, error) {
	return s.repo.ListWithCountryAndDocument()
}

func (s *Service) GetByID(id string) (schema.StageManagement, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(id string, input UpdateDTO) (schema.StageManagement, error) {
	stage, err := s.repo.GetByID(id)
	if err != nil {
		return schema.StageManagement{}, err
	}

	if input.CountryID != nil {
		stage.CountryID = *input.CountryID
	}
	if input.DocumentID != nil {
		stage.DocumentID = *input.DocumentID
	}

	if err := s.repo.Update(&stage); err != nil {
		return schema.StageManagement{}, err
	}
	return stage, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
