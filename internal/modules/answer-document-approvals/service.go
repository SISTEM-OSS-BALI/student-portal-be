package answerdocumentapprovals

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/schema"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func normalizeStatus(value *string) (string, error) {
	if value == nil {
		return "pending", nil
	}
	normalized := strings.ToLower(strings.TrimSpace(*value))
	if normalized == "" {
		return "pending", nil
	}
	switch normalized {
	case "pending", "approved", "rejected":
		return normalized, nil
	default:
		return "", errors.New("invalid status (use: pending, approved, rejected)")
	}
}

func (s *Service) CreateOrUpdate(input CreateDTO) (schema.AnswerDocumentApproval, error) {
	status, err := normalizeStatus(input.Status)
	if err != nil {
		return schema.AnswerDocumentApproval{}, err
	}

	approval, err := s.repo.GetByAnswerDocumentID(input.AnswerDocumentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newApproval := schema.AnswerDocumentApproval{
				AnswerDocumentID: input.AnswerDocumentID,
				StudentID:        input.StudentID,
				ReviewerID:       input.ReviewerID,
				Status:           status,
				Note:             input.Note,
				ReviewedAt:       input.ReviewedAt,
			}
			if (status == "approved" || status == "rejected") && newApproval.ReviewedAt == nil {
				now := time.Now()
				newApproval.ReviewedAt = &now
			}
			if err := s.repo.Create(&newApproval); err != nil {
				return schema.AnswerDocumentApproval{}, err
			}
			return newApproval, nil
		}
		return schema.AnswerDocumentApproval{}, err
	}

	approval.StudentID = input.StudentID
	approval.ReviewerID = input.ReviewerID
	approval.Status = status
	approval.Note = input.Note
	approval.ReviewedAt = input.ReviewedAt
	if (status == "approved" || status == "rejected") && approval.ReviewedAt == nil {
		now := time.Now()
		approval.ReviewedAt = &now
	}

	if err := s.repo.Update(&approval); err != nil {
		return schema.AnswerDocumentApproval{}, err
	}
	return approval, nil
}

func (s *Service) List(filter Filter) ([]schema.AnswerDocumentApproval, error) {
	return s.repo.List(filter)
}

func (s *Service) GetByID(id string) (schema.AnswerDocumentApproval, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(id string, input UpdateDTO) (schema.AnswerDocumentApproval, error) {
	approval, err := s.repo.GetByID(id)
	if err != nil {
		return schema.AnswerDocumentApproval{}, err
	}

	if input.AnswerDocumentID != nil {
		approval.AnswerDocumentID = *input.AnswerDocumentID
	}
	if input.StudentID != nil {
		approval.StudentID = *input.StudentID
	}
	if input.ReviewerID != nil {
		approval.ReviewerID = *input.ReviewerID
	}
	if input.Status != nil {
		status, err := normalizeStatus(input.Status)
		if err != nil {
			return schema.AnswerDocumentApproval{}, err
		}
		approval.Status = status
	}
	if input.Note != nil {
		approval.Note = input.Note
	}
	if input.ReviewedAt != nil {
		approval.ReviewedAt = input.ReviewedAt
	}
	if (approval.Status == "approved" || approval.Status == "rejected") && approval.ReviewedAt == nil {
		now := time.Now()
		approval.ReviewedAt = &now
	}

	if err := s.repo.Update(&approval); err != nil {
		return schema.AnswerDocumentApproval{}, err
	}
	return approval, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
