package countrysteps

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type CreateDTO struct {
	CountryID string `json:"country_id" binding:"required"`
	StepID    string `json:"step_id" binding:"required"`
}

type UpdateDTO struct {
	CountryID *string `json:"country_id"`
	StepID    *string `json:"step_id"`
}

type ResponseDTO struct {
	ID        string      `json:"id"`
	CountryID string      `json:"country_id"`
	StepID    string      `json:"step_id"`
	Country   *CountryDTO `json:"country,omitempty"`
	Step      *StepDTO    `json:"step,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type CountryDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type StepDTO struct {
	ID        string     `json:"id"`
	Label     string     `json:"label"`
	ChildIDs  []string   `json:"child_ids,omitempty"`
	Children  []ChildDTO `json:"children,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type ChildDTO struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

func NewResponseDTO(item schema.CountryStepsManagement) ResponseDTO {
	return ResponseDTO{
		ID:        item.ID,
		CountryID: item.CountryID,
		StepID:    item.StepID,
		Country:   newCountryDTO(item.Country),
		Step:      newStepDTO(item.Step),
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func NewResponseListDTO(items []schema.CountryStepsManagement) []ResponseDTO {
	out := make([]ResponseDTO, 0, len(items))
	for _, item := range items {
		out = append(out, NewResponseDTO(item))
	}
	return out
}

func newCountryDTO(country *schema.CountryManagement) *CountryDTO {
	if country == nil {
		return nil
	}
	return &CountryDTO{
		ID:   country.ID,
		Name: country.NameCountry,
	}
}

func newStepDTO(step *schema.StepsManagement) *StepDTO {
	if step == nil {
		return nil
	}
	childIDs := make([]string, 0, len(step.Children))
	children := make([]ChildDTO, 0, len(step.Children))
	for _, child := range step.Children {
		childIDs = append(childIDs, child.ID)
		children = append(children, ChildDTO{ID: child.ID, Label: child.Label})
	}
	return &StepDTO{
		ID:        step.ID,
		Label:     step.Label,
		ChildIDs:  childIDs,
		Children:  children,
		CreatedAt: step.CreatedAt,
		UpdatedAt: step.UpdatedAt,
	}
}
