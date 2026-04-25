package statementletteraiapprovals

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/schema"
)

var ErrDirectorRoleRequired = errors.New("only director can review statement letter")
var ErrApprovalDecisionRequired = errors.New("director review status must be APPROVED, REVISION_REQUESTED, or REJECTED")
var ErrStatementLetterNotSubmitted = errors.New("statement letter has not been submitted to director")
var ErrApprovalAssignedToAnotherDirector = errors.New("statement letter review is assigned to another director")

func normalizeApprovalStatus(value string) (schema.StatementLetterApprovalStatus, error) {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	if normalized == "" {
		return "", errors.New("status is required")
	}
	switch schema.StatementLetterApprovalStatus(normalized) {
	case schema.StatementLetterApprovalStatusPending,
		schema.StatementLetterApprovalStatusApproved,
		schema.StatementLetterApprovalStatusRevisionRequested,
		schema.StatementLetterApprovalStatusRejected:
		return schema.StatementLetterApprovalStatus(normalized), nil
	default:
		return "", errors.New("invalid status (use: PENDING, APPROVED, REVISION_REQUESTED, REJECTED)")
	}
}

func normalizeApprovalDecisionStatus(value string) (schema.StatementLetterApprovalStatus, error) {
	status, err := normalizeApprovalStatus(value)
	if err != nil {
		return "", err
	}

	switch status {
	case schema.StatementLetterApprovalStatusApproved,
		schema.StatementLetterApprovalStatusRevisionRequested,
		schema.StatementLetterApprovalStatusRejected:
		return status, nil
	default:
		return "", ErrApprovalDecisionRequired
	}
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOrUpdate(input CreateOrUpdateDTO, reviewerID string, reviewerRole schema.UserRole) (schema.StatementLetterAIApproval, error) {
	if reviewerRole != schema.UserRoleDirector {
		return schema.StatementLetterAIApproval{}, ErrDirectorRoleRequired
	}

	status, err := normalizeApprovalDecisionStatus(string(input.Status))
	if err != nil {
		return schema.StatementLetterAIApproval{}, err
	}

	doc, err := s.repo.GetDocumentByID(strings.TrimSpace(input.DocumentID))
	if err != nil {
		return schema.StatementLetterAIApproval{}, err
	}
	if doc.Status != schema.StatementLetterDocumentStatusSubmittedDirector {
		return schema.StatementLetterAIApproval{}, ErrStatementLetterNotSubmitted
	}

	currentApproval := doc.CurrentApproval
	if currentApproval == nil {
		current, currentErr := s.repo.GetCurrentByDocumentID(doc.ID)
		if currentErr != nil {
			if errors.Is(currentErr, gorm.ErrRecordNotFound) {
				return schema.StatementLetterAIApproval{}, ErrStatementLetterNotSubmitted
			}
			return schema.StatementLetterAIApproval{}, currentErr
		}
		currentApproval = &current
	}

	if currentApproval.ReviewerID != reviewerID {
		return schema.StatementLetterAIApproval{}, ErrApprovalAssignedToAnotherDirector
	}

	now := time.Now()
	fromStatus := currentApproval.Status
	currentApproval.Status = status
	currentApproval.Note = input.Note
	currentApproval.ReviewedAt = &now
	if err := s.repo.UpdateApproval(currentApproval); err != nil {
		return schema.StatementLetterAIApproval{}, err
	}

	doc.CurrentApprovalID = &currentApproval.ID
	switch status {
	case schema.StatementLetterApprovalStatusApproved:
		doc.Status = schema.StatementLetterDocumentStatusApproved
		doc.ApprovedAt = &now
		doc.RevisionRequestedAt = nil
	case schema.StatementLetterApprovalStatusRevisionRequested, schema.StatementLetterApprovalStatusRejected:
		doc.Status = schema.StatementLetterDocumentStatusRevisionRequested
		doc.RevisionRequestedAt = &now
		doc.ApprovedAt = nil
	}
	if doc.SubmittedToDirectorAt == nil {
		doc.SubmittedToDirectorAt = &now
	}
	if err := s.repo.UpdateDocument(&doc); err != nil {
		return schema.StatementLetterAIApproval{}, err
	}

	fromStatusCopy := fromStatus
	log := schema.StatementLetterAIApprovalLog{
		ApprovalID: currentApproval.ID,
		ActorID:    reviewerID,
		FromStatus: &fromStatusCopy,
		ToStatus:   status,
		Note:       input.Note,
	}
	if err := s.repo.CreateLog(&log); err != nil {
		return schema.StatementLetterAIApproval{}, err
	}

	return s.repo.GetByID(currentApproval.ID)
}

func (s *Service) List(filter Filter) ([]schema.StatementLetterAIApproval, error) {
	return s.repo.List(filter)
}

func (s *Service) GetByID(id string) (schema.StatementLetterAIApproval, error) {
	return s.repo.GetByID(strings.TrimSpace(id))
}

func (s *Service) Update(id string, input UpdateDTO, reviewerID string, reviewerRole schema.UserRole) (schema.StatementLetterAIApproval, error) {
	if reviewerRole != schema.UserRoleDirector {
		return schema.StatementLetterAIApproval{}, ErrDirectorRoleRequired
	}

	approval, err := s.repo.GetByID(strings.TrimSpace(id))
	if err != nil {
		return schema.StatementLetterAIApproval{}, err
	}

	statusValue := approval.Status
	if input.Status != nil {
		statusValue = *input.Status
	}

	return s.CreateOrUpdate(CreateOrUpdateDTO{
		DocumentID: approval.DocumentID,
		Status:     statusValue,
		Note:       firstNonNilString(input.Note, approval.Note),
	}, reviewerID, reviewerRole)
}

func firstNonNilString(primary, fallback *string) *string {
	if primary != nil {
		return primary
	}
	return fallback
}
