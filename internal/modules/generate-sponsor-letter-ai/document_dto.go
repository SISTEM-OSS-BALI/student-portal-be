package generatesponsorletterai

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type GeneratedDocumentUpsertDTO struct {
	StudentID             string                              `json:"student_id" binding:"required"`
	FileURL               string                              `json:"file_url" binding:"required"`
	FilePath              *string                             `json:"file_path"`
	FileName              *string                             `json:"file_name"`
	FileType              *string                             `json:"file_type"`
	WordFileURL           *string                             `json:"word_file_url"`
	WordFilePath          *string                             `json:"word_file_path"`
	WordFileName          *string                             `json:"word_file_name"`
	WordFileType          *string                             `json:"word_file_type"`
	Status                *schema.SponsorLetterDocumentStatus `json:"status"`
	Source                *schema.GeneratedDocumentSource     `json:"source"`
	SubmittedToDirectorAt *time.Time                          `json:"submitted_to_director_at"`
	ApprovedAt            *time.Time                          `json:"approved_at"`
	RevisionRequestedAt   *time.Time                          `json:"revision_requested_at"`
	CurrentApprovalID     *string                             `json:"current_approval_id"`
}

type GeneratedDocumentTemplateDTO struct {
	DocumentType           string                         `json:"document_type"`
	DefaultSource          schema.GeneratedDocumentSource `json:"default_source"`
	SupportsManualCreation bool                           `json:"supports_manual_creation"`
	SupportsAIGeneration   bool                           `json:"supports_ai_generation"`
	ChecklistVersion       string                         `json:"checklist_version,omitempty"`
	ChecklistItems         []string                       `json:"checklist_items,omitempty"`
	ChecklistSource        string                         `json:"checklist_source,omitempty"`
}

type SubmitToDirectorDTO struct {
	Note *string `json:"note"`
}

type GeneratedDocumentStudentDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GeneratedDocumentApprovalReviewerDTO struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Email string          `json:"email"`
	Role  schema.UserRole `json:"role"`
}

type GeneratedDocumentApprovalLogDTO struct {
	ID         string                                `json:"id"`
	ApprovalID string                                `json:"approval_id"`
	ActorID    string                                `json:"actor_id"`
	Actor      *GeneratedDocumentApprovalReviewerDTO `json:"actor,omitempty"`
	FromStatus *schema.SponsorLetterApprovalStatus   `json:"from_status,omitempty"`
	ToStatus   schema.SponsorLetterApprovalStatus    `json:"to_status"`
	Note       *string                               `json:"note,omitempty"`
	CreatedAt  time.Time                             `json:"created_at"`
}

type GeneratedDocumentApprovalDTO struct {
	ID         string                                `json:"id"`
	DocumentID string                                `json:"document_id"`
	ReviewerID string                                `json:"reviewer_id"`
	Reviewer   *GeneratedDocumentApprovalReviewerDTO `json:"reviewer,omitempty"`
	Status     schema.SponsorLetterApprovalStatus    `json:"status"`
	Note       *string                               `json:"note,omitempty"`
	ReviewedAt *time.Time                            `json:"reviewed_at,omitempty"`
	Logs       []GeneratedDocumentApprovalLogDTO     `json:"logs,omitempty"`
	CreatedAt  time.Time                             `json:"created_at"`
	UpdatedAt  time.Time                             `json:"updated_at"`
}

type GeneratedDocumentResponseDTO struct {
	ID                     string                             `json:"id"`
	StudentID              string                             `json:"student_id"`
	Student                *GeneratedDocumentStudentDTO       `json:"student,omitempty"`
	FileURL                string                             `json:"file_url"`
	FilePath               *string                            `json:"file_path,omitempty"`
	FileName               *string                            `json:"file_name,omitempty"`
	FileType               *string                            `json:"file_type,omitempty"`
	WordFileURL            *string                            `json:"word_file_url,omitempty"`
	WordFilePath           *string                            `json:"word_file_path,omitempty"`
	WordFileName           *string                            `json:"word_file_name,omitempty"`
	WordFileType           *string                            `json:"word_file_type,omitempty"`
	Status                 schema.SponsorLetterDocumentStatus `json:"status"`
	Source                 schema.GeneratedDocumentSource     `json:"source"`
	SupportsManualCreation bool                               `json:"supports_manual_creation"`
	SupportsAIGeneration   bool                               `json:"supports_ai_generation"`
	ChecklistVersion       string                             `json:"checklist_version,omitempty"`
	ChecklistItems         []string                           `json:"checklist_items,omitempty"`
	ChecklistSource        string                             `json:"checklist_source,omitempty"`
	SubmittedToDirectorAt  *time.Time                         `json:"submitted_to_director_at,omitempty"`
	ApprovedAt             *time.Time                         `json:"approved_at,omitempty"`
	RevisionRequestedAt    *time.Time                         `json:"revision_requested_at,omitempty"`
	CurrentApprovalID      *string                            `json:"current_approval_id,omitempty"`
	CanDownloadPDF         bool                               `json:"can_download_pdf"`
	DownloadPDFURL         *string                            `json:"download_pdf_url,omitempty"`
	CurrentApproval        *GeneratedDocumentApprovalDTO      `json:"current_approval,omitempty"`
	Approvals              []GeneratedDocumentApprovalDTO     `json:"approvals,omitempty"`
	CreatedAt              time.Time                          `json:"created_at"`
	UpdatedAt              time.Time                          `json:"updated_at"`
}

func newGeneratedDocumentApprovalReviewerDTO(user *schema.User) *GeneratedDocumentApprovalReviewerDTO {
	if user == nil {
		return nil
	}
	return &GeneratedDocumentApprovalReviewerDTO{ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role}
}

func newGeneratedDocumentApprovalLogDTO(item *schema.SponsorLetterAIApprovalLog) *GeneratedDocumentApprovalLogDTO {
	if item == nil {
		return nil
	}
	return &GeneratedDocumentApprovalLogDTO{ID: item.ID, ApprovalID: item.ApprovalID, ActorID: item.ActorID, Actor: newGeneratedDocumentApprovalReviewerDTO(item.Actor), FromStatus: item.FromStatus, ToStatus: item.ToStatus, Note: item.Note, CreatedAt: item.CreatedAt}
}

func newGeneratedDocumentApprovalDTO(item *schema.SponsorLetterAIApproval) *GeneratedDocumentApprovalDTO {
	if item == nil {
		return nil
	}
	response := &GeneratedDocumentApprovalDTO{ID: item.ID, DocumentID: item.DocumentID, ReviewerID: item.ReviewerID, Reviewer: newGeneratedDocumentApprovalReviewerDTO(item.Reviewer), Status: item.Status, Note: item.Note, ReviewedAt: item.ReviewedAt, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}
	if len(item.Logs) > 0 {
		response.Logs = make([]GeneratedDocumentApprovalLogDTO, 0, len(item.Logs))
		for _, logItem := range item.Logs {
			if dto := newGeneratedDocumentApprovalLogDTO(&logItem); dto != nil {
				response.Logs = append(response.Logs, *dto)
			}
		}
	}
	return response
}

func NewGeneratedDocumentResponseDTO(doc schema.GeneratedSponsorLetterAIDocument) GeneratedDocumentResponseDTO {
	response := GeneratedDocumentResponseDTO{ID: doc.ID, StudentID: doc.StudentID, FileURL: doc.FileURL, FilePath: doc.FilePath, FileName: doc.FileName, FileType: doc.FileType, WordFileURL: doc.WordFileURL, WordFilePath: doc.WordFilePath, WordFileName: doc.WordFileName, WordFileType: doc.WordFileType, Status: doc.Status, Source: doc.Source, SupportsManualCreation: true, SupportsAIGeneration: true, ChecklistVersion: sponsorChecklistVersion, ChecklistItems: sponsorChecklistItems(), ChecklistSource: sponsorChecklistSource(), SubmittedToDirectorAt: doc.SubmittedToDirectorAt, ApprovedAt: doc.ApprovedAt, RevisionRequestedAt: doc.RevisionRequestedAt, CurrentApprovalID: doc.CurrentApprovalID, CreatedAt: doc.CreatedAt, UpdatedAt: doc.UpdatedAt}
	if doc.Status == schema.SponsorLetterDocumentStatusApproved {
		response.CanDownloadPDF = true
		downloadURL := "/api/generate-sponsor-letter-ai/documents/" + doc.ID + "/download-pdf"
		response.DownloadPDFURL = &downloadURL
	}
	if doc.Student != nil {
		response.Student = &GeneratedDocumentStudentDTO{ID: doc.Student.ID, Name: doc.Student.Name, Email: doc.Student.Email}
	}
	response.CurrentApproval = newGeneratedDocumentApprovalDTO(doc.CurrentApproval)
	if len(doc.Approvals) > 0 {
		response.Approvals = make([]GeneratedDocumentApprovalDTO, 0, len(doc.Approvals))
		for _, item := range doc.Approvals {
			if dto := newGeneratedDocumentApprovalDTO(&item); dto != nil {
				response.Approvals = append(response.Approvals, *dto)
			}
		}
	}
	return response
}

func NewGeneratedDocumentResponseDTOs(docs []schema.GeneratedSponsorLetterAIDocument) []GeneratedDocumentResponseDTO {
	out := make([]GeneratedDocumentResponseDTO, 0, len(docs))
	for _, doc := range docs {
		out = append(out, NewGeneratedDocumentResponseDTO(doc))
	}
	return out
}

func NewGeneratedDocumentTemplateDTO() GeneratedDocumentTemplateDTO {
	return GeneratedDocumentTemplateDTO{
		DocumentType:           "sponsor_letter",
		DefaultSource:          schema.GeneratedDocumentSourceManual,
		SupportsManualCreation: true,
		SupportsAIGeneration:   true,
		ChecklistVersion:       sponsorChecklistVersion,
		ChecklistItems:         sponsorChecklistItems(),
		ChecklistSource:        sponsorChecklistSource(),
	}
}
