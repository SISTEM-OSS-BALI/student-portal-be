package generatesponsorletterai

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/schema"
)

var ErrSponsorLetterDirectorNotFound = errors.New("director not found")
var ErrSponsorLetterAlreadySubmitted = errors.New("sponsor letter already submitted to director")
var ErrSponsorLetterSubmissionNotCancelable = errors.New("sponsor letter submission can only be canceled while waiting for director review")
var ErrSponsorLetterDownloadNotApproved = errors.New("sponsor letter can only be downloaded after approval")
var ErrSponsorLetterDownloadUnavailable = errors.New("sponsor letter pdf is unavailable")

type GeneratedDocumentDownload struct {
	FilePath *string
	FileURL  string
	FileName string
	FileType string
}

type generatedDocumentStatusResolution struct {
	Status        schema.SponsorLetterDocumentStatus
	ResetApproval bool
}

func normalizeGeneratedDocumentStatus(input *schema.SponsorLetterDocumentStatus, current schema.SponsorLetterDocumentStatus) (generatedDocumentStatusResolution, error) {
	if input == nil {
		if current != "" {
			return generatedDocumentStatusResolution{Status: current}, nil
		}
		return generatedDocumentStatusResolution{Status: schema.SponsorLetterDocumentStatusDraft}, nil
	}

	normalized := strings.ToUpper(strings.TrimSpace(string(*input)))
	if normalized == "" {
		if current != "" {
			return generatedDocumentStatusResolution{Status: current}, nil
		}
		return generatedDocumentStatusResolution{Status: schema.SponsorLetterDocumentStatusDraft}, nil
	}

	switch normalized {
	case "GENERATED":
		return generatedDocumentStatusResolution{Status: schema.SponsorLetterDocumentStatusDraft, ResetApproval: true}, nil
	case "FINALIZED_PDF":
		if current != "" {
			return generatedDocumentStatusResolution{Status: current}, nil
		}
		return generatedDocumentStatusResolution{Status: schema.SponsorLetterDocumentStatusDraft}, nil
	}

	status := schema.SponsorLetterDocumentStatus(normalized)
	switch status {
	case schema.SponsorLetterDocumentStatusDraft:
		return generatedDocumentStatusResolution{Status: status, ResetApproval: true}, nil
	case schema.SponsorLetterDocumentStatusSubmittedDirector, schema.SponsorLetterDocumentStatusRevisionRequested, schema.SponsorLetterDocumentStatusApproved:
		return generatedDocumentStatusResolution{Status: status}, nil
	default:
		return generatedDocumentStatusResolution{}, errors.New("invalid sponsor letter document status")
	}
}

func resetDocumentApprovalState(doc *schema.GeneratedSponsorLetterAIDocument) {
	doc.Status = schema.SponsorLetterDocumentStatusDraft
	doc.SubmittedToDirectorAt = nil
	doc.ApprovedAt = nil
	doc.RevisionRequestedAt = nil
	doc.CurrentApprovalID = nil
	doc.CurrentApproval = nil
	doc.Approvals = nil
}

type GeneratedDocumentService struct {
	repo GeneratedDocumentRepository
}

func normalizeGeneratedDocumentSource(input *schema.GeneratedDocumentSource, current schema.GeneratedDocumentSource) (schema.GeneratedDocumentSource, error) {
	if input == nil {
		if current != "" {
			return current, nil
		}
		return schema.GeneratedDocumentSourceAI, nil
	}

	normalized := schema.GeneratedDocumentSource(strings.ToUpper(strings.TrimSpace(string(*input))))
	switch normalized {
	case schema.GeneratedDocumentSourceAI, schema.GeneratedDocumentSourceManual:
		return normalized, nil
	case "":
		if current != "" {
			return current, nil
		}
		return schema.GeneratedDocumentSourceAI, nil
	default:
		return "", errors.New("invalid document source (use: AI or MANUAL)")
	}
}

func NewGeneratedDocumentService(repo GeneratedDocumentRepository) *GeneratedDocumentService {
	return &GeneratedDocumentService{repo: repo}
}

func (s *GeneratedDocumentService) Upsert(input GeneratedDocumentUpsertDTO) (schema.GeneratedSponsorLetterAIDocument, error) {
	studentID := strings.TrimSpace(input.StudentID)
	fileURL := strings.TrimSpace(input.FileURL)
	if studentID == "" {
		return schema.GeneratedSponsorLetterAIDocument{}, errors.New("student_id is required")
	}
	if fileURL == "" {
		return schema.GeneratedSponsorLetterAIDocument{}, errors.New("file_url is required")
	}

	doc, err := s.repo.GetByStudentID(studentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			statusResult, statusErr := normalizeGeneratedDocumentStatus(input.Status, "")
			if statusErr != nil {
				return schema.GeneratedSponsorLetterAIDocument{}, statusErr
			}
			source, sourceErr := normalizeGeneratedDocumentSource(input.Source, "")
			if sourceErr != nil {
				return schema.GeneratedSponsorLetterAIDocument{}, sourceErr
			}
			newDoc := schema.GeneratedSponsorLetterAIDocument{StudentID: studentID, FileURL: fileURL, FilePath: input.FilePath, FileName: input.FileName, FileType: input.FileType, WordFileURL: input.WordFileURL, WordFilePath: input.WordFilePath, WordFileName: input.WordFileName, WordFileType: input.WordFileType, Status: statusResult.Status, Source: source}
			if !statusResult.ResetApproval {
				newDoc.SubmittedToDirectorAt = input.SubmittedToDirectorAt
				newDoc.ApprovedAt = input.ApprovedAt
				newDoc.RevisionRequestedAt = input.RevisionRequestedAt
				newDoc.CurrentApprovalID = input.CurrentApprovalID
			}
			if err := s.repo.Create(&newDoc); err != nil {
				return schema.GeneratedSponsorLetterAIDocument{}, err
			}
			return s.repo.GetByStudentID(studentID)
		}
		return schema.GeneratedSponsorLetterAIDocument{}, err
	}

	statusResult, err := normalizeGeneratedDocumentStatus(input.Status, doc.Status)
	if err != nil {
		return schema.GeneratedSponsorLetterAIDocument{}, err
	}
	source, err := normalizeGeneratedDocumentSource(input.Source, doc.Source)
	if err != nil {
		return schema.GeneratedSponsorLetterAIDocument{}, err
	}

	doc.FileURL = fileURL
	doc.FilePath = input.FilePath
	doc.FileName = input.FileName
	doc.FileType = input.FileType
	doc.WordFileURL = input.WordFileURL
	doc.WordFilePath = input.WordFilePath
	doc.WordFileName = input.WordFileName
	doc.WordFileType = input.WordFileType
	doc.Source = source
	if statusResult.ResetApproval {
		resetDocumentApprovalState(&doc)
	} else {
		doc.Status = statusResult.Status
		if input.SubmittedToDirectorAt != nil {
			doc.SubmittedToDirectorAt = input.SubmittedToDirectorAt
		}
		if input.ApprovedAt != nil {
			doc.ApprovedAt = input.ApprovedAt
		}
		if input.RevisionRequestedAt != nil {
			doc.RevisionRequestedAt = input.RevisionRequestedAt
		}
		if input.CurrentApprovalID != nil {
			doc.CurrentApprovalID = input.CurrentApprovalID
		}
	}
	if err := s.repo.Update(&doc); err != nil {
		return schema.GeneratedSponsorLetterAIDocument{}, err
	}
	return s.repo.GetByStudentID(studentID)
}

func (s *GeneratedDocumentService) SubmitToDirector(documentID string, actorID string, actorRole schema.UserRole, note *string) (schema.GeneratedSponsorLetterAIDocument, error) {
	doc, err := s.repo.GetByID(strings.TrimSpace(documentID))
	if err != nil {
		return schema.GeneratedSponsorLetterAIDocument{}, err
	}
	if strings.TrimSpace(doc.FileURL) == "" {
		return schema.GeneratedSponsorLetterAIDocument{}, errors.New("sponsor letter file is required before submit")
	}

	currentApproval := doc.CurrentApproval
	if currentApproval == nil {
		current, currentErr := s.repo.GetCurrentApprovalByDocumentID(doc.ID)
		if currentErr == nil {
			currentApproval = &current
		} else if !errors.Is(currentErr, gorm.ErrRecordNotFound) {
			return schema.GeneratedSponsorLetterAIDocument{}, currentErr
		}
	}

	if doc.Status == schema.SponsorLetterDocumentStatusSubmittedDirector && currentApproval != nil && currentApproval.Status == schema.SponsorLetterApprovalStatusPending {
		return schema.GeneratedSponsorLetterAIDocument{}, ErrSponsorLetterAlreadySubmitted
	}

	reviewerID := actorID
	if actorRole != schema.UserRoleDirector {
		director, directorErr := s.repo.GetDirector()
		if directorErr != nil {
			if errors.Is(directorErr, gorm.ErrRecordNotFound) {
				return schema.GeneratedSponsorLetterAIDocument{}, ErrSponsorLetterDirectorNotFound
			}
			return schema.GeneratedSponsorLetterAIDocument{}, directorErr
		}
		reviewerID = director.ID
	}

	now := time.Now()
	var fromStatus *schema.SponsorLetterApprovalStatus
	var approval schema.SponsorLetterAIApproval

	if currentApproval != nil && currentApproval.ReviewerID == reviewerID {
		statusCopy := currentApproval.Status
		fromStatus = &statusCopy
		currentApproval.Status = schema.SponsorLetterApprovalStatusPending
		currentApproval.Note = note
		currentApproval.ReviewedAt = nil
		if err := s.repo.UpdateApproval(currentApproval); err != nil {
			return schema.GeneratedSponsorLetterAIDocument{}, err
		}
		approval = *currentApproval
	} else {
		approval = schema.SponsorLetterAIApproval{DocumentID: doc.ID, ReviewerID: reviewerID, Status: schema.SponsorLetterApprovalStatusPending, Note: note}
		if currentApproval != nil {
			statusCopy := currentApproval.Status
			fromStatus = &statusCopy
		}
		if err := s.repo.CreateApproval(&approval); err != nil {
			return schema.GeneratedSponsorLetterAIDocument{}, err
		}
	}

	doc.Status = schema.SponsorLetterDocumentStatusSubmittedDirector
	doc.SubmittedToDirectorAt = &now
	doc.ApprovedAt = nil
	doc.RevisionRequestedAt = nil
	doc.CurrentApprovalID = &approval.ID
	if err := s.repo.Update(&doc); err != nil {
		return schema.GeneratedSponsorLetterAIDocument{}, err
	}

	log := schema.SponsorLetterAIApprovalLog{ApprovalID: approval.ID, ActorID: actorID, FromStatus: fromStatus, ToStatus: schema.SponsorLetterApprovalStatusPending, Note: note}
	if err := s.repo.CreateApprovalLog(&log); err != nil {
		return schema.GeneratedSponsorLetterAIDocument{}, err
	}
	return s.repo.GetByID(doc.ID)
}

func (s *GeneratedDocumentService) CancelSubmitToDirector(documentID string, actorID string, note *string) (schema.GeneratedSponsorLetterAIDocument, error) {
	doc, err := s.repo.GetByID(strings.TrimSpace(documentID))
	if err != nil {
		return schema.GeneratedSponsorLetterAIDocument{}, err
	}

	currentApproval := doc.CurrentApproval
	if currentApproval == nil {
		current, currentErr := s.repo.GetCurrentApprovalByDocumentID(doc.ID)
		if currentErr == nil {
			currentApproval = &current
		} else if !errors.Is(currentErr, gorm.ErrRecordNotFound) {
			return schema.GeneratedSponsorLetterAIDocument{}, currentErr
		}
	}

	if doc.Status != schema.SponsorLetterDocumentStatusSubmittedDirector ||
		currentApproval == nil ||
		currentApproval.Status != schema.SponsorLetterApprovalStatusPending {
		return schema.GeneratedSponsorLetterAIDocument{}, ErrSponsorLetterSubmissionNotCancelable
	}

	fromStatus := currentApproval.Status
	currentApproval.Status = schema.SponsorLetterApprovalStatusCanceled
	currentApproval.Note = note
	currentApproval.ReviewedAt = nil
	if err := s.repo.UpdateApproval(currentApproval); err != nil {
		return schema.GeneratedSponsorLetterAIDocument{}, err
	}

	doc.Status = schema.SponsorLetterDocumentStatusDraft
	doc.SubmittedToDirectorAt = nil
	doc.ApprovedAt = nil
	doc.RevisionRequestedAt = nil
	doc.CurrentApprovalID = nil
	if err := s.repo.Update(&doc); err != nil {
		return schema.GeneratedSponsorLetterAIDocument{}, err
	}

	fromStatusCopy := fromStatus
	log := schema.SponsorLetterAIApprovalLog{
		ApprovalID: currentApproval.ID,
		ActorID:    actorID,
		FromStatus: &fromStatusCopy,
		ToStatus:   schema.SponsorLetterApprovalStatusCanceled,
		Note:       note,
	}
	if err := s.repo.CreateApprovalLog(&log); err != nil {
		return schema.GeneratedSponsorLetterAIDocument{}, err
	}

	return s.repo.GetByID(doc.ID)
}

func (s *GeneratedDocumentService) ListByStudentID(studentID string) ([]schema.GeneratedSponsorLetterAIDocument, error) {
	return s.repo.ListByStudentID(strings.TrimSpace(studentID))
}
func (s *GeneratedDocumentService) GetByStudentID(studentID string) (schema.GeneratedSponsorLetterAIDocument, error) {
	return s.repo.GetByStudentID(strings.TrimSpace(studentID))
}

func (s *GeneratedDocumentService) GetApprovedDownload(documentID string) (GeneratedDocumentDownload, error) {
	doc, err := s.repo.GetByID(strings.TrimSpace(documentID))
	if err != nil {
		return GeneratedDocumentDownload{}, err
	}
	if doc.Status != schema.SponsorLetterDocumentStatusApproved {
		return GeneratedDocumentDownload{}, ErrSponsorLetterDownloadNotApproved
	}

	fileName := "sponsor-letter-" + doc.ID + ".pdf"
	if doc.FileName != nil && strings.TrimSpace(*doc.FileName) != "" {
		fileName = strings.TrimSpace(*doc.FileName)
	}

	fileType := "application/pdf"
	if doc.FileType != nil && strings.TrimSpace(*doc.FileType) != "" {
		fileType = strings.TrimSpace(*doc.FileType)
	}

	if doc.FilePath != nil {
		filePath := strings.TrimSpace(*doc.FilePath)
		if filePath != "" {
			return GeneratedDocumentDownload{
				FilePath: &filePath,
				FileURL:  strings.TrimSpace(doc.FileURL),
				FileName: fileName,
				FileType: fileType,
			}, nil
		}
	}

	fileURL := strings.TrimSpace(doc.FileURL)
	if fileURL != "" {
		return GeneratedDocumentDownload{
			FileURL:  fileURL,
			FileName: fileName,
			FileType: fileType,
		}, nil
	}

	return GeneratedDocumentDownload{}, ErrSponsorLetterDownloadUnavailable
}
