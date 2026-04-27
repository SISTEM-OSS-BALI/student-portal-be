package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type CreateDTO struct {
	Name             string  `json:"name" binding:"required"`
	Email            string  `json:"email" binding:"required,email"`
	Role             string  `json:"role" binding:"omitempty,oneof=student admission director" default:"student"`
	Password         string  `json:"password" binding:"required,min=8"`
	StageID          *string `json:"stage_id"`
	CurrentStepID    *string `json:"current_step_id"`
	VisaStatus       *string `json:"visa_status"`
	StudentStatus    *string `json:"student_status"`
	NameConsultant   *string `json:"name_consultant"`
	NoPhone          *string `json:"no_phone"`
	NameCampus       *string `json:"name_campus"`
	Degree           *string `json:"degree"`
	NameDegree       *string `json:"name_degree"`
	DocumentConsentSignatureURL *string `json:"document_consent_signature_url"`
	DocumentConsentProofPhotoURL *string `json:"document_consent_proof_photo_url"`
	DocumentConsentSignedAt *time.Time `json:"document_consent_signed_at"`
	DocumentConsentSigned *bool `json:"document_consent_signed"`
	VisaType         *string `json:"visa_type"`
	TranslationQuota int     `json:"translation_quota"`
}

type UpdateDTO struct {
	Name             *string `json:"name"`
	Email            *string `json:"email"`
	Role             *string `json:"role"`
	NoPhone          *string `json:"no_phone"`
	StageID          *string `json:"stage_id"`
	CurrentStepID    *string `json:"current_step_id"`
	VisaStatus       *string `json:"visa_status"`
	StudentStatus    *string `json:"student_status"`
	NameConsultant   *string `json:"name_consultant"`
	NameCampus       *string `json:"name_campus"`
	NameDegree       *string `json:"name_degree"`
	Degree           *string `json:"degree"`
	DocumentConsentSignatureURL *string `json:"document_consent_signature_url"`
	DocumentConsentProofPhotoURL *string `json:"document_consent_proof_photo_url"`
	DocumentConsentSignedAt *time.Time `json:"document_consent_signed_at"`
	DocumentConsentSigned *bool `json:"document_consent_signed"`
	VisaType         *string `json:"visa_type"`
	TranslationQuota *int    `json:"translation_quota"`
}

type PatchQuotaTranslationDTO struct {
	TranslationQuota int `json:"translation_quota" binding:"required,gte=0"`
}

type PatchVisaStatusDTO struct {
	VisaStatus *string `json:"visa_status" binding:"required"`
}

type PatchStudentStatusDTO struct {
	StudentStatus *string `json:"student_status" binding:"required"`
}

type PatchDocumentConsentDTO struct {
	DocumentConsentSignatureURL  *string    `json:"document_consent_signature_url"`
	DocumentConsentProofPhotoURL *string    `json:"document_consent_proof_photo_url"`
	DocumentConsentSignedAt      *time.Time `json:"document_consent_signed_at"`
	DocumentConsentSigned        *bool      `json:"document_consent_signed"`
}

type ResponseDTO struct {
	ID                     string           `json:"id"`
	Name                   string           `json:"name"`
	Email                  string           `json:"email"`
	Role                   string           `json:"role"`
	NoPhone                *string          `json:"no_phone,omitempty"`
	StageID                *string          `json:"stage_id,omitempty"`
	CurrentStepID          *string          `json:"current_step_id,omitempty"`
	VisaStatus             *string          `json:"visa_status,omitempty"`
	VisaGrantedAt          *time.Time       `json:"visa_granted_at,omitempty"`
	VisaGrantDurationDays  *int             `json:"visa_grant_duration_days,omitempty"`
	VisaGrantDurationLabel *string          `json:"visa_grant_duration_label,omitempty"`
	StudentStatus          string           `json:"student_status"`
	StudentStatusUpdatedByID    *string          `json:"student_status_updated_by_id,omitempty"`
	StudentStatusUpdatedByName  *string          `json:"student_status_updated_by_name,omitempty"`
	StudentStatusUpdatedAt      *time.Time       `json:"student_status_updated_at,omitempty"`
	StudentStatusUpdatedAtLabel *string          `json:"student_status_updated_at_label,omitempty"`
	NameConsultant         *string          `json:"name_consultant,omitempty"`
	Stage                  *StageDTO        `json:"stage,omitempty"`
	NameCampus             *string          `json:"name_campus,omitempty"`
	VisaType               *string          `json:"visa_type,omitempty"`
	Degree                 *string          `json:"degree,omitempty"`
	NameDegree             *string          `json:"name_degree,omitempty"`
	DocumentConsentSignatureURL *string `json:"document_consent_signature_url"`
	DocumentConsentProofPhotoURL *string `json:"document_consent_proof_photo_url"`
	DocumentConsentSignedAt *time.Time `json:"document_consent_signed_at"`
	DocumentConsentSigned *bool `json:"document_consent_signed"`
	TranslationQuota       int              `json:"translation_quota"`
	NotesStudent           []NoteStudentDTO `json:"notes,omitempty"`
	JoinedAt               time.Time        `json:"joined_at"`
	CreatedAt              time.Time        `json:"created_at"`
	UpdatedAt              time.Time        `json:"updated_at"`
}

type StageDTO struct {
	ID         string       `json:"id"`
	CountryID  string       `json:"country_id"`
	DocumentID string       `json:"document_id"`
	Country    *CountryDTO  `json:"country,omitempty"`
	Document   *DocumentDTO `json:"document,omitempty"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

type CountryDTO struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	Steps []StepDTO `json:"steps,omitempty"`
}

type DocumentDTO struct {
	ID                string                   `json:"id"`
	Label             string                   `json:"label"`
	InternalCode      string                   `json:"internal_code"`
	FileType          string                   `json:"file_type"`
	Category          string                   `json:"category"`
	TranslationNeeded schema.TranslationNeeded `json:"translation_needed"`
	Required          bool                     `json:"required"`
	AutoRenamePattern schema.AutoRenamePattern `json:"auto_rename_pattern"`
	Notes             string                   `json:"notes"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
}

type NoteStudentDTO struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StepDTO struct {
	ID        string     `json:"id"`
	Label     string     `json:"label"`
	ChildIDs  []string   `json:"child_ids,omitempty"`
	Children  []ChildDTO `json:"children,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type ChildDTO struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}



func boolPtr(v bool) *bool {
	return &v
}

func NewResponseDTO(user schema.User) ResponseDTO {
	visaGrantDurationDays, visaGrantDurationLabel := buildVisaGrantDuration(user.CreatedAt, user.VisaGrantedAt)
	studentStatusUpdatedByName := optionalUserName(user.StudentStatusUpdatedBy)
	studentStatusUpdatedAtLabel := formatStudentStatusUpdatedAtLabel(user.StudentStatusUpdatedAt)
	documentConsentSigned := boolPtr(user.DocumentConsentSigned)

	return ResponseDTO{
		ID:                            user.ID,
		Name:                          user.Name,
		Email:                         user.Email,
		Role:                          string(user.Role),
		NoPhone:                       user.NoPhone,
		StageID:                       user.StageID,
		CurrentStepID:                 user.CurrentStepID,
		VisaStatus:                    user.VisaStatus,
		VisaGrantedAt:                 user.VisaGrantedAt,
		VisaGrantDurationDays:         visaGrantDurationDays,
		VisaGrantDurationLabel:        visaGrantDurationLabel,
		StudentStatus:                 string(user.StudentStatus),
		StudentStatusUpdatedByID:      user.StudentStatusUpdatedByID,
		StudentStatusUpdatedByName:    studentStatusUpdatedByName,
		StudentStatusUpdatedAt:        user.StudentStatusUpdatedAt,
		StudentStatusUpdatedAtLabel:   studentStatusUpdatedAtLabel,
		NameConsultant:                user.NameConsultant,
		NameCampus:                    user.NameCampus,
		Degree:                        user.Degree,
		NameDegree:                    user.NameDegree,
		DocumentConsentSignatureURL:   user.DocumentConsentSignatureURL,
		DocumentConsentProofPhotoURL:  user.DocumentConsentProofPhotoURL,
		DocumentConsentSignedAt:       user.DocumentConsentSignedAt,
		DocumentConsentSigned:         documentConsentSigned,
		VisaType:                      user.VisaType,
		JoinedAt:                      user.CreatedAt,
		CreatedAt:                     user.CreatedAt,
		UpdatedAt:                     user.UpdatedAt,
		Stage:                         newStageDTO(user.Stage),
		NotesStudent:                  newNoteStudentListDTO(user.NotesStudent),
		TranslationQuota:              user.TranslationQuota,
	}
}

func NewResponseListDTO(users []schema.User) []ResponseDTO {
	out := make([]ResponseDTO, 0, len(users))
	for _, u := range users {
		out = append(out, NewResponseDTO(u))
	}
	return out
}

func newStageDTO(stage *schema.StageManagement) *StageDTO {
	if stage == nil {
		return nil
	}
	return &StageDTO{
		ID:         stage.ID,
		CountryID:  stage.CountryID,
		DocumentID: stage.DocumentID,
		Country:    newCountryDTO(stage.Country),
		Document:   newDocumentDTO(stage.Document),
		CreatedAt:  stage.CreatedAt,
		UpdatedAt:  stage.UpdatedAt,
	}
}

func newCountryDTO(country *schema.CountryManagement) *CountryDTO {
	if country == nil {
		return nil
	}
	return &CountryDTO{
		ID:    country.ID,
		Name:  country.NameCountry,
		Steps: newStepListDTO(country.CountrySteps),
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

func newNoteStudentListDTO(notes []schema.NoteStudent) []NoteStudentDTO {
	if len(notes) == 0 {
		return nil
	}
	out := make([]NoteStudentDTO, 0, len(notes))
	for _, note := range notes {
		out = append(out, NoteStudentDTO{
			ID:        note.ID,
			UserID:    note.UserID,
			Content:   note.Content,
			CreatedAt: note.CreatedAt,
			UpdatedAt: note.UpdatedAt,
		})
	}
	return out
}

func newStepListDTO(items []schema.CountryStepsManagement) []StepDTO {
	if len(items) == 0 {
		return nil
	}
	out := make([]StepDTO, 0, len(items))
	for _, item := range items {
		if item.Step == nil {
			continue
		}
		childIDs := make([]string, 0, len(item.Step.Children))
		children := make([]ChildDTO, 0, len(item.Step.Children))
		for _, child := range item.Step.Children {
			childIDs = append(childIDs, child.ID)
			children = append(children, ChildDTO{ID: child.ID, Label: child.Label})
		}
		out = append(out, StepDTO{
			ID:        item.Step.ID,
			Label:     item.Step.Label,
			ChildIDs:  childIDs,
			Children:  children,
			CreatedAt: item.Step.CreatedAt,
			UpdatedAt: item.Step.UpdatedAt,
		})
	}
	return out
}

func buildVisaGrantDuration(joinedAt time.Time, visaGrantedAt *time.Time) (*int, *string) {
	if joinedAt.IsZero() || visaGrantedAt == nil || visaGrantedAt.IsZero() {
		return nil, nil
	}

	totalDays := int(visaGrantedAt.Sub(joinedAt).Hours() / 24)
	if totalDays < 0 {
		totalDays = 0
	}

	label := formatDurationLabel(totalDays)
	return &totalDays, &label
}

func formatDurationLabel(totalDays int) string {
	if totalDays <= 0 {
		return "0 hari"
	}

	years := totalDays / 365
	remainingDays := totalDays % 365
	months := remainingDays / 30
	days := remainingDays % 30

	parts := make([]string, 0, 3)
	if years > 0 {
		parts = append(parts, fmt.Sprintf("%d tahun", years))
	}
	if months > 0 {
		parts = append(parts, fmt.Sprintf("%d bulan", months))
	}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d hari", days))
	}
	if len(parts) == 0 {
		return "0 hari"
	}

	return strings.Join(parts, " ")
}

func optionalUserName(user *schema.User) *string {
	if user == nil || strings.TrimSpace(user.Name) == "" {
		return nil
	}

	name := user.Name
	return &name
}

func formatStudentStatusUpdatedAtLabel(value *time.Time) *string {
	if value == nil || value.IsZero() {
		return nil
	}

	weekdayNames := map[time.Weekday]string{
		time.Sunday:    "Minggu",
		time.Monday:    "Senin",
		time.Tuesday:   "Selasa",
		time.Wednesday: "Rabu",
		time.Thursday:  "Kamis",
		time.Friday:    "Jumat",
		time.Saturday:  "Sabtu",
	}
	monthNames := map[time.Month]string{
		time.January:   "Januari",
		time.February:  "Februari",
		time.March:     "Maret",
		time.April:     "April",
		time.May:       "Mei",
		time.June:      "Juni",
		time.July:      "Juli",
		time.August:    "Agustus",
		time.September: "September",
		time.October:   "Oktober",
		time.November:  "November",
		time.December:  "Desember",
	}

	localTime := value.Local()
	label := fmt.Sprintf(
		"%s, %02d %s %d %02d:%02d",
		weekdayNames[localTime.Weekday()],
		localTime.Day(),
		monthNames[localTime.Month()],
		localTime.Year(),
		localTime.Hour(),
		localTime.Minute(),
	)
	return &label
}
