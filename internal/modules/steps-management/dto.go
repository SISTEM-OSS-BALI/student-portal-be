package steps

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type CreateDTO struct {
	Label    string   `json:"label" binding:"required"`
	ChildIDs []string `json:"child_ids,omitempty"`
}

type UpdateDTO struct {
	Label    *string  `json:"label"`
	ChildIDs *[]string `json:"child_ids,omitempty"`
}

type ResponseDTO struct {
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

func NewResponseDTO(step schema.StepsManagement) ResponseDTO {
	childIDs := make([]string, 0, len(step.Children))
	children := make([]ChildDTO, 0, len(step.Children))
	for _, child := range step.Children {
		childIDs = append(childIDs, child.ID)
		children = append(children, ChildDTO{ID: child.ID, Label: child.Label})
	}
	return ResponseDTO{
		ID:        step.ID,
		Label:     step.Label,
		ChildIDs:  childIDs,
		Children:  children,
		CreatedAt: step.CreatedAt,
		UpdatedAt: step.UpdatedAt,
	}
}

func NewResponseListDTO(steps []schema.StepsManagement) []ResponseDTO {
	out := make([]ResponseDTO, 0, len(steps))
	for _, s := range steps {
		out = append(out, NewResponseDTO(s))
	}
	return out
}

func newChildDTO(child *schema.ChildStepsManagement) *ChildDTO {
	if child == nil {
		return nil
	}
	return &ChildDTO{
		ID:    child.ID,
		Label: child.Label,
	}
}
