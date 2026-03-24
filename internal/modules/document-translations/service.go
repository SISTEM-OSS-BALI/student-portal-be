package documenttranslations

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(input CreateDTO) (schema.DocumentTranslation, error) {
	status := input.Status
	if status == nil {
		value := "pending"
		status = &value
	}

	item := schema.DocumentTranslation{
		StudentID:        input.StudentID,
		DocumentID:       input.DocumentID,
		UploaderID:       input.UploaderID,
		AnswerDocumentID: input.AnswerDocumentID,
		FileURL:          input.FileURL,
		FilePath:         input.FilePath,
		FileName:         input.FileName,
		FileType:         input.FileType,
		PageCount:        input.PageCount,
		Status:           status,
	}

	if input.CreatedAt != nil {
		item.CreatedAt = *input.CreatedAt
		item.UpdatedAt = *input.CreatedAt
	}

	if err := s.repo.Create(&item); err != nil {
		return schema.DocumentTranslation{}, err
	}
	return item, nil
}

func (s *Service) List(filter Filter) ([]schema.DocumentTranslation, error) {
	return s.repo.List(filter)
}

func (s *Service) GetByID(id string) (schema.DocumentTranslation, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(id string, input UpdateDTO) (schema.DocumentTranslation, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return schema.DocumentTranslation{}, err
	}

	if input.StudentID != nil {
		item.StudentID = *input.StudentID
	}
	if input.DocumentID != nil {
		item.DocumentID = *input.DocumentID
	}
	if input.UploaderID != nil {
		item.UploaderID = *input.UploaderID
	}
	if input.AnswerDocumentID != nil {
		item.AnswerDocumentID = input.AnswerDocumentID
	}
	if input.FileURL != nil {
		item.FileURL = *input.FileURL
	}
	if input.FilePath != nil {
		item.FilePath = input.FilePath
	}
	if input.FileName != nil {
		item.FileName = input.FileName
	}
	if input.FileType != nil {
		item.FileType = input.FileType
	}
	if input.PageCount != nil {
		item.PageCount = *input.PageCount
	}
	if input.Status != nil {
		item.Status = input.Status
	}

	item.UpdatedAt = time.Now()

	if err := s.repo.Update(&item); err != nil {
		return schema.DocumentTranslation{}, err
	}
	return item, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *Service) UpdateUserTranslationQuota(studentID string, pageCount int) error {
	return s.repo.UpdateUserTranslationQuota(studentID, pageCount)
}
