package documents

import "github.com/username/gin-gorm-api/internal/schema"

type CreateDTO struct {
	Label             string `json:"label" binding:"required"`
	InternalCode      string `json:"internal_code" binding:"required"`
	FileType          string `json:"file_type" binding:"required"`
	Category          string `json:"category" binding:"required"`
	TranslationNeeded string `json:"translation_needed" binding:"required"`
	Required          *bool  `json:"required" binding:"required"`
	AutoRenamePattern schema.AutoRenamePattern `json:"auto_rename_pattern"`
	Notes             string `json:"notes"`
}

type UpdateDTO struct {
	Label             *string `json:"label"`
	InternalCode      *string `json:"internal_code"`
	FileType          *string `json:"file_type"`
	Category          *string `json:"category"`
	TranslationNeeded *string `json:"translation_needed"`
	Required          *bool   `json:"required"`
	AutoRenamePattern *schema.AutoRenamePattern `json:"auto_rename_pattern"`
	Notes             *string `json:"notes"`
}

type ResponseDTO struct {
	ID                string                   `json:"id"`
	Label             string                   `json:"label"`
	InternalCode      string                   `json:"internal_code"`
	FileType          string                   `json:"file_type"`
	Category          string                   `json:"category"`
	TranslationNeeded schema.TranslationNeeded `json:"translation_needed"`
	Required          bool                     `json:"required"`
	AutoRenamePattern schema.AutoRenamePattern `json:"auto_rename_pattern"`
	Notes             string                   `json:"notes"`
}

func NewResponseDTO(doc schema.DocumentsManagement) ResponseDTO {
	return ResponseDTO{
		ID:                doc.ID,
		Label:             doc.Label,
		InternalCode:      doc.InternalCode,
		FileType:          doc.FileType,
		Category:          doc.Category,
		TranslationNeeded: doc.TranslationNeeded,
		Required:          doc.Required,
		AutoRenamePattern: doc.AutoRenamePattern,
		Notes:             doc.Notes,
	}
}

func NewResponseListDTO(docs []schema.DocumentsManagement) []ResponseDTO {
	out := make([]ResponseDTO, 0, len(docs))
	for _, d := range docs {
		out = append(out, NewResponseDTO(d))
	}
	return out
}

type CountPDFPagesRequestDTO struct {
	URL string `json:"url" binding:"required,url"`
}

type CountPDFPagesResponseDTO struct {
	URL       string `json:"url"`
	PageCount int    `json:"page_count"`
}

