package visatype

import (
	"errors"
	"strings"

	"github.com/lucsky/cuid"

	"github.com/username/gin-gorm-api/internal/schema"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(input CreateDTO) (schema.VisaTypeManagement, error) {
	name := strings.TrimSpace(input.Name)
	countryID := strings.TrimSpace(input.CountryID)
	if name == "" || countryID == "" {
		return schema.VisaTypeManagement{}, errors.New("name and country_id are required")
	}

	item := schema.VisaTypeManagement{
		ID:        cuid.New(),
		Name:      name,
		CountryID: countryID,
	}
	if err := s.repo.Create(&item); err != nil {
		return schema.VisaTypeManagement{}, err
	}
	return item, nil
}

func (s *Service) List(filter Filter) ([]schema.VisaTypeManagement, error) {
	return s.repo.List(filter)
}

func (s *Service) GetByID(id string) (schema.VisaTypeManagement, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(id string, input UpdateDTO) (schema.VisaTypeManagement, error) {
	updates := map[string]interface{}{}
	if input.Name != nil {
		updates["name"] = strings.TrimSpace(*input.Name)
	}
	if input.CountryID != nil {
		updates["country_id"] = strings.TrimSpace(*input.CountryID)
	}
	if len(updates) == 0 {
		return s.repo.GetByID(id)
	}
	return s.repo.Update(id, updates)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

