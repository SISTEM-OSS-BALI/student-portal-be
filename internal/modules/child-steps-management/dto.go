package childsteps

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type CreateDTO struct {
	Label string `json:"label" binding:"required"`
}

type UpdateDTO struct {
	Label *string `json:"label"`
}

type ResponseDTO struct {
	ID        string    `json:"id"`
	Label     string    `json:"label"`
	Steps     []StepDTO `json:"steps,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StepDTO struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

func NewResponseDTO(child schema.ChildStepsManagement) ResponseDTO {
	return ResponseDTO{
		ID:        child.ID,
		Label:     child.Label,
		Steps:     newStepListDTO(child.Steps),
		CreatedAt: child.CreatedAt,
		UpdatedAt: child.UpdatedAt,
	}
}

func NewResponseListDTO(children []schema.ChildStepsManagement) []ResponseDTO {
	out := make([]ResponseDTO, 0, len(children))
	for _, c := range children {
		out = append(out, NewResponseDTO(c))
	}
	return out
}

func newStepListDTO(steps []schema.StepsManagement) []StepDTO {
	if len(steps) == 0 {
		return nil
	}
	out := make([]StepDTO, 0, len(steps))
	for _, s := range steps {
		out = append(out, StepDTO{
			ID:    s.ID,
			Label: s.Label,
		})
	}
	return out
}
