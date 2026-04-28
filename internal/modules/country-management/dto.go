package country

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type CreateDTO struct {
	Name string `json:"name" binding:"required"`
}

type UpdateDTO struct {
	Name *string `json:"name"`
}

type ResponseDTO struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	DocumentTotal int64     `json:"document_total"`
	StepTotal     int64     `json:"step_total"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewResponseDTO(country schema.CountryManagement, documentTotal int64, stepTotal int64) ResponseDTO {
	return ResponseDTO{
		ID:            country.ID,
		Name:          country.NameCountry,
		DocumentTotal: documentTotal,
		StepTotal:     stepTotal,
		CreatedAt:     country.CreatedAt,
		UpdatedAt:     country.UpdatedAt,
	}
}

func NewResponseListDTO(countries []schema.CountryManagement, totals map[string]int64, stepTotals map[string]int64) []ResponseDTO {
	out := make([]ResponseDTO, 0, len(countries))
	for _, c := range countries {
		out = append(out, NewResponseDTO(c, totals[c.ID], stepTotals[c.ID]))
	}
	return out
}