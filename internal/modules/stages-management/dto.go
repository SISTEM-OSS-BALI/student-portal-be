package stages

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type CreateDTO struct {
	CountryID  string `json:"country_id" binding:"required"`
	DocumentID string `json:"document_id" binding:"required"`
}

type UpdateDTO struct {
	CountryID  *string `json:"country_id"`
	DocumentID *string `json:"document_id"`
}

type ResponseDTO struct {
	ID          string    `json:"id"`
	CountryID   string    `json:"country_id"`
	DocumentID  string    `json:"document_id"`
	Country     *CountryDTO  `json:"country,omitempty"`
	Document    *DocumentDTO `json:"document,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CountryDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DocumentDTO struct {
	ID                string                    `json:"id"`
	Label             string                    `json:"label"`
	InternalCode      string                    `json:"internal_code"`
	FileType          string                    `json:"file_type"`
	Category          string                    `json:"category"`
	TranslationNeeded schema.TranslationNeeded  `json:"translation_needed"`
	Required          bool                      `json:"required"`
	AutoRenamePattern schema.AutoRenamePattern  `json:"auto_rename_pattern"`
	Notes             string                    `json:"notes"`
	CreatedAt         time.Time                 `json:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at"`
}

func NewResponseDTO(stage schema.StageManagement) ResponseDTO {
	return ResponseDTO{
		ID:         stage.ID,
		CountryID:  stage.CountryID,
		DocumentID: stage.DocumentID,
		Country:    newCountryDTO(stage.Country),
		Document:   newDocumentDTO(stage.Document),
		CreatedAt:  stage.CreatedAt,
		UpdatedAt:  stage.UpdatedAt,
	}
}

func NewResponseListDTO(stages []schema.StageManagement) []ResponseDTO {
	out := make([]ResponseDTO, 0, len(stages))
	for _, s := range stages {
		out = append(out, NewResponseDTO(s))
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

func newDocumentDTO(doc *schema.DocumentsManagement) *DocumentDTO {
	if doc == nil {
		return nil
	}
	return &DocumentDTO{
		ID:                doc.ID,
		Label:             doc.Label,
		InternalCode:      doc.InternalCode,
		FileType:          doc.FileType,
		Category:          doc.Category,
		TranslationNeeded: doc.TranslationNeeded,
		Required:          doc.Required,
		AutoRenamePattern: doc.AutoRenamePattern,
		Notes:             doc.Notes,
		CreatedAt:         doc.CreatedAt,
		UpdatedAt:         doc.UpdatedAt,
	}
}
