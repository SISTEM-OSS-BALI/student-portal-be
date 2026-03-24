package generatecvai

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/schema"
)

type GeneratedDocumentService struct {
	repo GeneratedDocumentRepository
}

func NewGeneratedDocumentService(repo GeneratedDocumentRepository) *GeneratedDocumentService {
	return &GeneratedDocumentService{repo: repo}
}

func (s *GeneratedDocumentService) Upsert(input GeneratedDocumentUpsertDTO) (schema.GeneratedCVAIDocument, error) {
	studentID := strings.TrimSpace(input.StudentID)
	fileURL := strings.TrimSpace(input.FileURL)
	if studentID == "" {
		return schema.GeneratedCVAIDocument{}, errors.New("student_id is required")
	}
	if fileURL == "" {
		return schema.GeneratedCVAIDocument{}, errors.New("file_url is required")
	}

	doc, err := s.repo.GetByStudentID(studentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newDoc := schema.GeneratedCVAIDocument{
				StudentID:    studentID,
				FileURL:      fileURL,
				FilePath:     input.FilePath,
				FileName:     input.FileName,
				FileType:     input.FileType,
				WordFileURL:  input.WordFileURL,
				WordFilePath: input.WordFilePath,
				WordFileName: input.WordFileName,
				WordFileType: input.WordFileType,
				Status:       input.Status,
			}
			if err := s.repo.Create(&newDoc); err != nil {
				return schema.GeneratedCVAIDocument{}, err
			}
			return s.repo.GetByStudentID(studentID)
		}
		return schema.GeneratedCVAIDocument{}, err
	}

	doc.FileURL = fileURL
	doc.FilePath = input.FilePath
	doc.FileName = input.FileName
	doc.FileType = input.FileType
	doc.WordFileURL = input.WordFileURL
	doc.WordFilePath = input.WordFilePath
	doc.WordFileName = input.WordFileName
	doc.WordFileType = input.WordFileType
	doc.Status = input.Status
	if err := s.repo.Update(&doc); err != nil {
		return schema.GeneratedCVAIDocument{}, err
	}
	return s.repo.GetByStudentID(studentID)
}

func (s *GeneratedDocumentService) ListByStudentID(studentID string) ([]schema.GeneratedCVAIDocument, error) {
	return s.repo.ListByStudentID(strings.TrimSpace(studentID))
}

func (s *GeneratedDocumentService) GetByStudentID(studentID string) (schema.GeneratedCVAIDocument, error) {
	return s.repo.GetByStudentID(strings.TrimSpace(studentID))
}
