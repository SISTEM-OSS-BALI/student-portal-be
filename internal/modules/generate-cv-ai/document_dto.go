package generatecvai

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type GeneratedDocumentUpsertDTO struct {
	StudentID    string  `json:"student_id" binding:"required"`
	FileURL      string  `json:"file_url" binding:"required"`
	FilePath     *string `json:"file_path"`
	FileName     *string `json:"file_name"`
	FileType     *string `json:"file_type"`
	WordFileURL  *string `json:"word_file_url"`
	WordFilePath *string `json:"word_file_path"`
	WordFileName *string `json:"word_file_name"`
	WordFileType *string `json:"word_file_type"`
	Status       *string `json:"status"`
}

type GeneratedDocumentResponseDTO struct {
	ID           string                       `json:"id"`
	StudentID    string                       `json:"student_id"`
	Student      *GeneratedDocumentStudentDTO `json:"student,omitempty"`
	FileURL      string                       `json:"file_url"`
	FilePath     *string                      `json:"file_path,omitempty"`
	FileName     *string                      `json:"file_name,omitempty"`
	FileType     *string                      `json:"file_type,omitempty"`
	WordFileURL  *string                      `json:"word_file_url,omitempty"`
	WordFilePath *string                      `json:"word_file_path,omitempty"`
	WordFileName *string                      `json:"word_file_name,omitempty"`
	WordFileType *string                      `json:"word_file_type,omitempty"`
	Status       *string                      `json:"status,omitempty"`
	CreatedAt    time.Time                    `json:"created_at"`
	UpdatedAt    time.Time                    `json:"updated_at"`
}

type GeneratedDocumentStudentDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewGeneratedDocumentResponseDTO(doc schema.GeneratedCVAIDocument) GeneratedDocumentResponseDTO {
	response := GeneratedDocumentResponseDTO{
		ID:           doc.ID,
		StudentID:    doc.StudentID,
		FileURL:      doc.FileURL,
		FilePath:     doc.FilePath,
		FileName:     doc.FileName,
		FileType:     doc.FileType,
		WordFileURL:  doc.WordFileURL,
		WordFilePath: doc.WordFilePath,
		WordFileName: doc.WordFileName,
		WordFileType: doc.WordFileType,
		Status:       doc.Status,
		CreatedAt:    doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
	}

	if doc.Student != nil {
		response.Student = &GeneratedDocumentStudentDTO{
			ID:    doc.Student.ID,
			Name:  doc.Student.Name,
			Email: doc.Student.Email,
		}
	}

	return response
}

func NewGeneratedDocumentResponseListDTO(docs []schema.GeneratedCVAIDocument) []GeneratedDocumentResponseDTO {
	out := make([]GeneratedDocumentResponseDTO, 0, len(docs))
	for _, doc := range docs {
		out = append(out, NewGeneratedDocumentResponseDTO(doc))
	}
	return out
}
