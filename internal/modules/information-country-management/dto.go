package informationcountrymanagement

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type CreateDTO struct {
	Slug        string  `json:"slug"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Priority    string  `json:"priority"`
	CountryID   string  `json:"country_id"`
}

type UpdateDTO struct {
	Slug        *string `json:"slug"`
	Title       *string `json:"title"`
	Description *string `json:"description,omitempty"`
	Priority    *string `json:"priority"`
	CountryID   *string `json:"country_id"`
}

type ResponseDTO struct {
	ID          string      `json:"id"`
	Slug        string      `json:"slug"`
	Title       string      `json:"title"`
	Description *string     `json:"description,omitempty"`
	Priority    string      `json:"priority"`
	CountryID   string      `json:"country_id"`
	Country     *CountryDTO `json:"country,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type CountryDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewResponseDTO(data schema.InformationCountryManagement) ResponseDTO {
	return ResponseDTO{
		ID:          data.ID,
		Slug:        data.Slug,
		Title:       data.Title,
		Description: data.Description,
		Priority:    data.Priority,
		CountryID:   data.CountryID,
		Country:     newCountryDTO(data.Country),
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func NewResponseListDTO(items []schema.InformationCountryManagement) []ResponseDTO {
	out := make([]ResponseDTO, 0, len(items))
	for _, item := range items {
		out = append(out, NewResponseDTO(item))
	}
	return out
}

func newCountryDTO(country schema.CountryManagement) *CountryDTO {
	if country.ID == "" {
		return nil
	}

	return &CountryDTO{
		ID:   country.ID,
		Name: country.NameCountry,
	}
}
