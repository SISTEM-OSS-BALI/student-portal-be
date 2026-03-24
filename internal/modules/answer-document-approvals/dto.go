package answerdocumentapprovals

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type CreateDTO struct {
	AnswerDocumentID string     `json:"answer_document_id" binding:"required"`
	StudentID        string     `json:"student_id" binding:"required"`
	ReviewerID       string     `json:"reviewer_id" binding:"required"`
	Status           *string    `json:"status"`
	Note             *string    `json:"note"`
	ReviewedAt       *time.Time `json:"reviewed_at"`
}

type UpdateDTO struct {
	AnswerDocumentID *string    `json:"answer_document_id"`
	StudentID        *string    `json:"student_id"`
	ReviewerID       *string    `json:"reviewer_id"`
	Status           *string    `json:"status"`
	Note             *string    `json:"note"`
	ReviewedAt       *time.Time `json:"reviewed_at"`
}

type ResponseDTO struct {
	ID               string     `json:"id"`
	AnswerDocumentID string     `json:"answer_document_id"`
	StudentID        string     `json:"student_id"`
	ReviewerID       string     `json:"reviewer_id"`
	Status           string     `json:"status"`
	Note             *string    `json:"note,omitempty"`
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func NewResponseDTO(approval schema.AnswerDocumentApproval) ResponseDTO {
	return ResponseDTO{
		ID:               approval.ID,
		AnswerDocumentID: approval.AnswerDocumentID,
		StudentID:        approval.StudentID,
		ReviewerID:       approval.ReviewerID,
		Status:           approval.Status,
		Note:             approval.Note,
		ReviewedAt:       approval.ReviewedAt,
		CreatedAt:        approval.CreatedAt,
		UpdatedAt:        approval.UpdatedAt,
	}
}

func NewResponseListDTO(items []schema.AnswerDocumentApproval) []ResponseDTO {
	out := make([]ResponseDTO, 0, len(items))
	for _, item := range items {
		out = append(out, NewResponseDTO(item))
	}
	return out
}
