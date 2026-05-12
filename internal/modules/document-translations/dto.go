package documenttranslations

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type CreateDTO struct {
	StudentID             string     `json:"student_id" binding:"required"`
	DocumentID            string     `json:"document_id" binding:"required"`
	UploaderID            string     `json:"uploader_id" binding:"required"`
	AnswerDocumentID      *string    `json:"answer_document_id"`
	FileURL               string     `json:"file_url" binding:"required"`
	FilePath              *string    `json:"file_path"`
	FileName              *string    `json:"file_name"`
	FileType              *string    `json:"file_type"`
	PageCount             int        `json:"page_count"`
	IsExistingTranslation bool       `json:"is_existing_translation"`
	Status                *string    `json:"status"`
	CreatedAt             *time.Time `json:"created_at"`
}

type UpdateDTO struct {
	StudentID             *string `json:"student_id"`
	DocumentID            *string `json:"document_id"`
	UploaderID            *string `json:"uploader_id"`
	AnswerDocumentID      *string `json:"answer_document_id"`
	FileURL               *string `json:"file_url"`
	FilePath              *string `json:"file_path"`
	FileName              *string `json:"file_name"`
	FileType              *string `json:"file_type"`
	PageCount             *int    `json:"page_count"`
	IsExistingTranslation *bool   `json:"is_existing_translation"`
	Status                *string `json:"status"`
}

type ResponseDTO struct {
	ID                    string    `json:"id"`
	StudentID             string    `json:"student_id"`
	DocumentID            string    `json:"document_id"`
	UploaderID            string    `json:"uploader_id"`
	AnswerDocumentID      *string   `json:"answer_document_id,omitempty"`
	FileURL               string    `json:"file_url"`
	FilePath              *string   `json:"file_path,omitempty"`
	FileName              *string   `json:"file_name,omitempty"`
	FileType              *string   `json:"file_type,omitempty"`
	PageCount             int       `json:"page_count"`
	IsExistingTranslation bool      `json:"is_existing_translation"`
	Status                *string   `json:"status,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func NewResponseDTO(item schema.DocumentTranslation) ResponseDTO {
	return ResponseDTO{
		ID:                    item.ID,
		StudentID:             item.StudentID,
		DocumentID:            item.DocumentID,
		UploaderID:            item.UploaderID,
		AnswerDocumentID:      item.AnswerDocumentID,
		FileURL:               item.FileURL,
		FilePath:              item.FilePath,
		FileName:              item.FileName,
		FileType:              item.FileType,
		PageCount:             item.PageCount,
		IsExistingTranslation: item.IsExistingTranslation,
		Status:                item.Status,
		CreatedAt:             item.CreatedAt,
		UpdatedAt:             item.UpdatedAt,
	}
}

func NewResponseListDTO(items []schema.DocumentTranslation) []ResponseDTO {
	out := make([]ResponseDTO, 0, len(items))
	for _, item := range items {
		out = append(out, NewResponseDTO(item))
	}
	return out
}
