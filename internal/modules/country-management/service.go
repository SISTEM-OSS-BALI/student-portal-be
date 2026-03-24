package country

import "github.com/username/gin-gorm-api/internal/schema"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(input CreateDTO) (schema.CountryManagement, error) {
	country := schema.CountryManagement{NameCountry: input.Name}
	if err := s.repo.Create(&country); err != nil {
		return schema.CountryManagement{}, err
	}
	return country, nil
}

func (s *Service) List() ([]schema.CountryManagement, error) {
	return s.repo.List()
}

func (s *Service) ListWithDocumentAndStepCounts() ([]schema.CountryManagement, map[string]int64, map[string]int64, error) {
	countries, err := s.repo.List()
	if err != nil {
		return nil, nil, nil, err
	}
	ids := make([]string, 0, len(countries))
	for _, c := range countries {
		ids = append(ids, c.ID)
	}
	documentCounts, err := s.repo.CountDocumentsByCountryIDs(ids)
	if err != nil {
		return nil, nil, nil, err
	}
	stepCounts, err := s.repo.CountStepsByCountryIDs(ids)
	if err != nil {
		return nil, nil, nil, err
	}
	return countries, documentCounts, stepCounts, nil
}

func (s *Service) ListWithDocumentCounts() ([]schema.CountryManagement, map[string]int64, error) {
	countries, err := s.repo.List()
	if err != nil {
		return nil, nil, err
	}
	ids := make([]string, 0, len(countries))
	for _, c := range countries {
		ids = append(ids, c.ID)
	}
	counts, err := s.repo.CountDocumentsByCountryIDs(ids)
	if err != nil {
		return nil, nil, err
	}
	return countries, counts, nil
}

func (s *Service) GetByID(id string) (schema.CountryManagement, error) {
	return s.repo.GetByID(id)
}

func (s *Service) GetByIDWithDocumentAndStepCount(id string) (schema.CountryManagement, int64, int64, error) {
	country, err := s.repo.GetByID(id)
	if err != nil {
		return schema.CountryManagement{}, 0, 0, err
	}
	documentCount, err := s.repo.CountDocumentsByCountryID(id)
	if err != nil {
		return schema.CountryManagement{}, 0, 0, err
	}
	stepCount, err := s.repo.CountStepsByCountryID(id)
	if err != nil {
		return schema.CountryManagement{}, 0, 0, err
	}
	return country, documentCount, stepCount, nil
}

func (s *Service) GetByIDWithDocumentCount(id string) (schema.CountryManagement, int64, error) {
	country, err := s.repo.GetByID(id)
	if err != nil {
		return schema.CountryManagement{}, 0, err
	}
	count, err := s.repo.CountDocumentsByCountryID(id)
	if err != nil {
		return schema.CountryManagement{}, 0, err
	}
	return country, count, nil
}

func (s *Service) DocumentCount(id string) (int64, error) {
	return s.repo.CountDocumentsByCountryID(id)
}

func (s *Service) StepCount(id string) (int64, error) {
	return s.repo.CountStepsByCountryID(id)
}

func (s *Service) Update(id string, input UpdateDTO) (schema.CountryManagement, error) {
	country, err := s.repo.GetByID(id)
	if err != nil {
		return schema.CountryManagement{}, err
	}
	if input.Name != nil {
		country.NameCountry = *input.Name
	}
	if err := s.repo.Update(&country); err != nil {
		return schema.CountryManagement{}, err
	}
	return country, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
