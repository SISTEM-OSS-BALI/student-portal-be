package sponsorletteraiapprovals

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type CreateOrUpdateDTO struct {
	DocumentID string                             `json:"document_id" binding:"required"`
	Status     schema.SponsorLetterApprovalStatus `json:"status" binding:"required"`
	Note       *string                            `json:"note"`
}

type UpdateDTO struct {
	Status *schema.SponsorLetterApprovalStatus `json:"status"`
	Note   *string                             `json:"note"`
}

type ReviewerDTO struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Email string          `json:"email"`
	Role  schema.UserRole `json:"role"`
}

type DocumentSummaryDTO struct {
	ID                    string                             `json:"id"`
	StudentID             string                             `json:"student_id"`
	FileURL               string                             `json:"file_url"`
	FilePath              *string                            `json:"file_path,omitempty"`
	FileName              *string                            `json:"file_name,omitempty"`
	FileType              *string                            `json:"file_type,omitempty"`
	WordFileURL           *string                            `json:"word_file_url,omitempty"`
	WordFilePath          *string                            `json:"word_file_path,omitempty"`
	WordFileName          *string                            `json:"word_file_name,omitempty"`
	WordFileType          *string                            `json:"word_file_type,omitempty"`
	Status                schema.SponsorLetterDocumentStatus `json:"status"`
	SubmittedToDirectorAt *time.Time                         `json:"submitted_to_director_at,omitempty"`
	ApprovedAt            *time.Time                         `json:"approved_at,omitempty"`
	RevisionRequestedAt   *time.Time                         `json:"revision_requested_at,omitempty"`
	CurrentApprovalID     *string                            `json:"current_approval_id,omitempty"`
}

type LogResponseDTO struct {
	ID         string                              `json:"id"`
	ApprovalID string                              `json:"approval_id"`
	ActorID    string                              `json:"actor_id"`
	Actor      *ReviewerDTO                        `json:"actor,omitempty"`
	FromStatus *schema.SponsorLetterApprovalStatus `json:"from_status,omitempty"`
	ToStatus   schema.SponsorLetterApprovalStatus  `json:"to_status"`
	Note       *string                             `json:"note,omitempty"`
	CreatedAt  time.Time                           `json:"created_at"`
}

type ResponseDTO struct {
	ID         string                             `json:"id"`
	DocumentID string                             `json:"document_id"`
	ReviewerID string                             `json:"reviewer_id"`
	Reviewer   *ReviewerDTO                       `json:"reviewer,omitempty"`
	Document   *DocumentSummaryDTO                `json:"document,omitempty"`
	Status     schema.SponsorLetterApprovalStatus `json:"status"`
	Note       *string                            `json:"note,omitempty"`
	ReviewedAt *time.Time                         `json:"reviewed_at,omitempty"`
	Logs       []LogResponseDTO                   `json:"logs,omitempty"`
	CreatedAt  time.Time                          `json:"created_at"`
	UpdatedAt  time.Time                          `json:"updated_at"`
}

func newReviewerDTO(user *schema.User) *ReviewerDTO {
	if user == nil {
		return nil
	}
	return &ReviewerDTO{ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role}
}
func newDocumentSummaryDTO(doc *schema.GeneratedSponsorLetterAIDocument) *DocumentSummaryDTO {
	if doc == nil {
		return nil
	}
	return &DocumentSummaryDTO{ID: doc.ID, StudentID: doc.StudentID, FileURL: doc.FileURL, FilePath: doc.FilePath, FileName: doc.FileName, FileType: doc.FileType, WordFileURL: doc.WordFileURL, WordFilePath: doc.WordFilePath, WordFileName: doc.WordFileName, WordFileType: doc.WordFileType, Status: doc.Status, SubmittedToDirectorAt: doc.SubmittedToDirectorAt, ApprovedAt: doc.ApprovedAt, RevisionRequestedAt: doc.RevisionRequestedAt, CurrentApprovalID: doc.CurrentApprovalID}
}
func newLogResponseDTO(log schema.SponsorLetterAIApprovalLog) LogResponseDTO {
	return LogResponseDTO{ID: log.ID, ApprovalID: log.ApprovalID, ActorID: log.ActorID, Actor: newReviewerDTO(log.Actor), FromStatus: log.FromStatus, ToStatus: log.ToStatus, Note: log.Note, CreatedAt: log.CreatedAt}
}
func NewResponseDTO(item schema.SponsorLetterAIApproval) ResponseDTO {
	logs := make([]LogResponseDTO, 0, len(item.Logs))
	for _, log := range item.Logs {
		logs = append(logs, newLogResponseDTO(log))
	}
	return ResponseDTO{ID: item.ID, DocumentID: item.DocumentID, ReviewerID: item.ReviewerID, Reviewer: newReviewerDTO(item.Reviewer), Document: newDocumentSummaryDTO(item.Document), Status: item.Status, Note: item.Note, ReviewedAt: item.ReviewedAt, Logs: logs, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}
}
func NewResponseListDTO(items []schema.SponsorLetterAIApproval) []ResponseDTO {
	out := make([]ResponseDTO, 0, len(items))
	for _, item := range items {
		out = append(out, NewResponseDTO(item))
	}
	return out
}
