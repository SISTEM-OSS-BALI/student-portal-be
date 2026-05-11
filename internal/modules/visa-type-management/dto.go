package visatype

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type CreateDTO struct {
	Name      string `json:"name" binding:"required"`
	CountryID string `json:"country_id" binding:"required"`
}

type UpdateDTO struct {
	Name      *string `json:"name"`
	CountryID *string `json:"country_id"`
}

type ResponseDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CountryID string    `json:"country_id"`
	CountryName *string  `json:"country_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewResponseDTO(item schema.VisaTypeManagement) ResponseDTO {
	var countryName *string
	if item.Country != nil {
		name := item.Country.NameCountry
		countryName = &name
	}
	return ResponseDTO{
		ID:        item.ID,
		Name:      item.Name,
		CountryID: item.CountryID,
		CountryName: countryName,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func NewResponseListDTO(items []schema.VisaTypeManagement) []ResponseDTO {
	out := make([]ResponseDTO, 0, len(items))
	for _, item := range items {
		out = append(out, NewResponseDTO(item))
	}
	return out
}

type Filter struct {
	CountryID *string
}
