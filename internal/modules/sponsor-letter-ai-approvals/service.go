package sponsorletteraiapprovals

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/schema"
)

var ErrDirectorRoleRequired = errors.New("only director can review sponsor letter")
var ErrApprovalDecisionRequired = errors.New("director review status must be APPROVED, REVISION_REQUESTED, or REJECTED")
var ErrSponsorLetterNotSubmitted = errors.New("sponsor letter has not been submitted to director")
var ErrApprovalAssignedToAnotherDirector = errors.New("sponsor letter review is assigned to another director")

func normalizeApprovalStatus(value string) (schema.SponsorLetterApprovalStatus, error) {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	if normalized == "" {
		return "", errors.New("status is required")
	}
	switch schema.SponsorLetterApprovalStatus(normalized) {
	case schema.SponsorLetterApprovalStatusPending, schema.SponsorLetterApprovalStatusApproved, schema.SponsorLetterApprovalStatusRevisionRequested, schema.SponsorLetterApprovalStatusRejected:
		return schema.SponsorLetterApprovalStatus(normalized), nil
	default:
		return "", errors.New("invalid status (use: PENDING, APPROVED, REVISION_REQUESTED, REJECTED)")
	}
}
func normalizeApprovalDecisionStatus(value string) (schema.SponsorLetterApprovalStatus, error) {
	status, err := normalizeApprovalStatus(value)
	if err != nil {
		return "", err
	}
	switch status {
	case schema.SponsorLetterApprovalStatusApproved, schema.SponsorLetterApprovalStatusRevisionRequested, schema.SponsorLetterApprovalStatusRejected:
		return status, nil
	default:
		return "", ErrApprovalDecisionRequired
	}
}

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }
func (s *Service) CreateOrUpdate(input CreateOrUpdateDTO, reviewerID string, reviewerRole schema.UserRole) (schema.SponsorLetterAIApproval, error) {
	if reviewerRole != schema.UserRoleDirector {
		return schema.SponsorLetterAIApproval{}, ErrDirectorRoleRequired
	}
	status, err := normalizeApprovalDecisionStatus(string(input.Status))
	if err != nil {
		return schema.SponsorLetterAIApproval{}, err
	}
	doc, err := s.repo.GetDocumentByID(strings.TrimSpace(input.DocumentID))
	if err != nil {
		return schema.SponsorLetterAIApproval{}, err
	}
	if doc.Status != schema.SponsorLetterDocumentStatusSubmittedDirector {
		return schema.SponsorLetterAIApproval{}, ErrSponsorLetterNotSubmitted
	}
	currentApproval := doc.CurrentApproval
	if currentApproval == nil {
		current, currentErr := s.repo.GetCurrentByDocumentID(doc.ID)
		if currentErr != nil {
			if errors.Is(currentErr, gorm.ErrRecordNotFound) {
				return schema.SponsorLetterAIApproval{}, ErrSponsorLetterNotSubmitted
			}
			return schema.SponsorLetterAIApproval{}, currentErr
		}
		currentApproval = &current
	}
	if currentApproval.ReviewerID != reviewerID {
		return schema.SponsorLetterAIApproval{}, ErrApprovalAssignedToAnotherDirector
	}
	now := time.Now()
	fromStatus := currentApproval.Status
	currentApproval.Status = status
	currentApproval.Note = input.Note
	currentApproval.ReviewedAt = &now
	if err := s.repo.UpdateApproval(currentApproval); err != nil {
		return schema.SponsorLetterAIApproval{}, err
	}
	doc.CurrentApprovalID = &currentApproval.ID
	switch status {
	case schema.SponsorLetterApprovalStatusApproved:
		doc.Status = schema.SponsorLetterDocumentStatusApproved
		doc.ApprovedAt = &now
		doc.RevisionRequestedAt = nil
	case schema.SponsorLetterApprovalStatusRevisionRequested, schema.SponsorLetterApprovalStatusRejected:
		doc.Status = schema.SponsorLetterDocumentStatusRevisionRequested
		doc.RevisionRequestedAt = &now
		doc.ApprovedAt = nil
	}
	if doc.SubmittedToDirectorAt == nil {
		doc.SubmittedToDirectorAt = &now
	}
	if err := s.repo.UpdateDocument(&doc); err != nil {
		return schema.SponsorLetterAIApproval{}, err
	}
	fromStatusCopy := fromStatus
	log := schema.SponsorLetterAIApprovalLog{ApprovalID: currentApproval.ID, ActorID: reviewerID, FromStatus: &fromStatusCopy, ToStatus: status, Note: input.Note}
	if err := s.repo.CreateLog(&log); err != nil {
		return schema.SponsorLetterAIApproval{}, err
	}
	return s.repo.GetByID(currentApproval.ID)
}
func (s *Service) List(filter Filter) ([]schema.SponsorLetterAIApproval, error) {
	return s.repo.List(filter)
}
func (s *Service) GetByID(id string) (schema.SponsorLetterAIApproval, error) {
	return s.repo.GetByID(strings.TrimSpace(id))
}
func (s *Service) Update(id string, input UpdateDTO, reviewerID string, reviewerRole schema.UserRole) (schema.SponsorLetterAIApproval, error) {
	if reviewerRole != schema.UserRoleDirector {
		return schema.SponsorLetterAIApproval{}, ErrDirectorRoleRequired
	}
	approval, err := s.repo.GetByID(strings.TrimSpace(id))
	if err != nil {
		return schema.SponsorLetterAIApproval{}, err
	}
	statusValue := approval.Status
	if input.Status != nil {
		statusValue = *input.Status
	}
	return s.CreateOrUpdate(CreateOrUpdateDTO{DocumentID: approval.DocumentID, Status: statusValue, Note: firstNonNilString(input.Note, approval.Note)}, reviewerID, reviewerRole)
}
func firstNonNilString(primary, fallback *string) *string {
	if primary != nil {
		return primary
	}
	return fallback
}
