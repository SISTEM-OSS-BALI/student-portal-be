package notesstudent

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type CreateDTO struct {
	UserID  string `json:"user_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateDTO struct {
	Content *string `json:"content"`
}

type ResponseDTO struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewResponseDTO(note schema.NoteStudent) ResponseDTO {
	return ResponseDTO{
		ID:        note.ID,
		UserID:    note.UserID,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

func NewResponseListDTO(notes []schema.NoteStudent) []ResponseDTO {
	out := make([]ResponseDTO, 0, len(notes))
	for _, note := range notes {
		out = append(out, NewResponseDTO(note))
	}
	return out
}
