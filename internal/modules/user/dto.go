package user

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type CreateDTO struct {
	Name             string  `json:"name" binding:"required"`
	Email            string  `json:"email" binding:"required,email"`
	Role             string  `json:"role" binding:"required,oneof=student admission director" default:"student"`
	Password         string  `json:"password" binding:"required,min=8"`
	StageID          *string `json:"stage_id"`
	NoPhone          *string `json:"no_phone"`
	NameCampus       *string `json:"name_campus"`
	Degree           *string `json:"degree"`
	NameDegree       *string `json:"name_degree"`
	VisaType         *string `json:"visa_type"`
	TranslationQuota int     `json:"translation_quota"`
}

type UpdateDTO struct {
	Name             *string `json:"name"`
	Email            *string `json:"email"`
	Role             *string `json:"role"`
	NoPhone          *string `json:"no_phone"`
	StageID          *string `json:"stage_id"`
	NameCampus       *string `json:"name_campus"`
	NameDegree       *string `json:"name_degree"`
	Degree           *string `json:"degree"`
	VisaType         *string `json:"visa_type"`
	TranslationQuota *int    `json:"translation_quota"`
}

type PatchQuotaTranslationDTO struct {
	TranslationQuota int `json:"translation_quota" binding:"required,gte=0"`
}

type ResponseDTO struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Email            string           `json:"email"`
	Role             string           `json:"role"`
	NoPhone          *string          `json:"no_phone,omitempty"`
	StageID          *string          `json:"stage_id,omitempty"`
	Stage            *StageDTO        `json:"stage,omitempty"`
	Status           string           `json:"status"`
	NameCampus       *string          `json:"name_campus,omitempty"`
	VisaType         *string          `json:"visa_type,omitempty"`
	Degree           *string          `json:"degree,omitempty"`
	NameDegree       *string          `json:"name_degree,omitempty"`
	TranslationQuota int              `json:"translation_quota"`
	NotesStudent     []NoteStudentDTO `json:"notes,omitempty"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
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

func NewResponseDTO(user schema.User) ResponseDTO {
	return ResponseDTO{
		ID:               user.ID,
		Name:             user.Name,
		Email:            user.Email,
		Role:             string(user.Role),
		NoPhone:          user.NoPhone,
		StageID:          user.StageID,
		Status:           string(user.Status),
		NameCampus:       user.NameCampus,
		Degree:           user.Degree,
		NameDegree:       user.NameDegree,
		VisaType:         user.VisaType,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		Stage:            newStageDTO(user.Stage),
		NotesStudent:     newNoteStudentListDTO(user.NotesStudent),
		TranslationQuota: user.TranslationQuota,
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
