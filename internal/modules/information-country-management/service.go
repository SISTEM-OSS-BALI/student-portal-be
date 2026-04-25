package informationcountrymanagement

import (
	"strings"

	"github.com/username/gin-gorm-api/internal/schema"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func normalizeString(value string) string {
	return strings.TrimSpace(value)
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	return &trimmed
}

func normalizePriority(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "normal"
	}
	return strings.ToLower(trimmed)
}

func (s *Service) Create(
	slug string,
	title string,
	description *string,
	priority string,
	countryID string,
) (schema.InformationCountryManagement, error) {
	data := schema.InformationCountryManagement{
		Slug:        normalizeString(slug),
		Title:       normalizeString(title),
		Description: normalizeOptionalString(description),
		Priority:    normalizePriority(priority),
		CountryID:   normalizeString(countryID),
	}

	if err := s.repo.Create(&data); err != nil {
		return schema.InformationCountryManagement{}, err
	}

	return data, nil
}

func (s *Service) List() ([]schema.InformationCountryManagement, error) {
	return s.repo.List()
}

func (s *Service) GetByID(id string) (schema.InformationCountryManagement, error) {
	return s.repo.GetByID(id)
}

func (s *Service) GetBySlug(slug string) (schema.InformationCountryManagement, error) {
	return s.repo.GetBySlug(slug)
}

func (s *Service) ListByCountryID(countryID string) ([]schema.InformationCountryManagement, error) {
	return s.repo.ListByCountryID(countryID)
}

func (s *Service) Update(
	id string,
	slug, title, description, priority, countryID *string,
) (schema.InformationCountryManagement, error) {
	data, err := s.repo.GetByID(id)
	if err != nil {
		return schema.InformationCountryManagement{}, err
	}

	if slug != nil {
		data.Slug = normalizeString(*slug)
	}

	if title != nil {
		data.Title = normalizeString(*title)
	}

	if description != nil {
		data.Description = normalizeOptionalString(description)
	}

	if priority != nil {
		data.Priority = normalizePriority(*priority)
	}

	if countryID != nil {
		data.CountryID = normalizeString(*countryID)
	}

	if err := s.repo.Update(&data); err != nil {
		return schema.InformationCountryManagement{}, err
	}

	return data, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
