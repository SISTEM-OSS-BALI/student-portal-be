package questions

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type QuestionBaseCreateDTO struct {
	Name                     string  `json:"name" binding:"required"`
	Desc                     *string `json:"desc"`
	TypeCountry              string  `json:"type_country" binding:"required"`
	CountryID                *string `json:"country_id"`
	AllowMultipleSubmissions *bool   `json:"allow_multiple_submissions"`
	Active                   *bool   `json:"active"`
	Version                  *int    `json:"version"`
}

type QuestionBaseUpdateDTO struct {
	Name                     *string `json:"name"`
	Desc                     *string `json:"desc"`
	TypeCountry              *string `json:"type_country"`
	CountryID                *string `json:"country_id"`
	AllowMultipleSubmissions *bool   `json:"allow_multiple_submissions"`
	Active                   *bool   `json:"active"`
	Version                  *int    `json:"version"`
}

type QuestionBaseResponseDTO struct {
	ID                       string    `json:"id"`
	Name                     string    `json:"name"`
	Desc                     *string   `json:"desc,omitempty"`
	TypeCountry              string    `json:"type_country"`
	CountryID                *string   `json:"country_id,omitempty"`
	AllowMultipleSubmissions bool      `json:"allow_multiple_submissions"` 
	Active                   bool      `json:"active"`
	Version                  int       `json:"version"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type QuestionCreateDTO struct {
	BaseID      string              `json:"base_id" binding:"required"`
	Text        string              `json:"text" binding:"required"`
	InputType   schema.QuestionType `json:"input_type" binding:"required"`
	Required    *bool               `json:"required"`
	Order       *int                `json:"order"`
	HelpText    *string             `json:"help_text"`
	Placeholder *string             `json:"placeholder"`
	MinLength   *int                `json:"min_length"`
	MaxLength   *int                `json:"max_length"`
	Active      *bool               `json:"active"`
}

type QuestionUpdateDTO struct {
	BaseID      *string              `json:"base_id"`
	Text        *string              `json:"text"`
	InputType   *schema.QuestionType `json:"input_type"`
	Required    *bool                `json:"required"`
	Order       *int                 `json:"order"`
	HelpText    *string              `json:"help_text"`
	Placeholder *string              `json:"placeholder"`
	MinLength   *int                 `json:"min_length"`
	MaxLength   *int                 `json:"max_length"`
	Active      *bool                `json:"active"`
}

type QuestionOptionCreateDTO struct {
	QuestionID string `json:"question_id" binding:"required"`
	Label      string `json:"label" binding:"required"`
	Value      string `json:"value" binding:"required"`
	Order      *int   `json:"order"`
	Active     *bool  `json:"active"`
}

type QuestionOptionUpdateDTO struct {
	QuestionID *string `json:"question_id"`
	Label      *string `json:"label"`
	Value      *string `json:"value"`
	Order      *int    `json:"order"`
	Active     *bool   `json:"active"`
}

type AnswerQuestionCreateDTO struct {
	SubmissionID      *string  `json:"submission_id"`
	QuestionID        string   `json:"question_id" binding:"required"`
	AnswerText        *string  `json:"answer_text"`
	StudentID       *string  `json:"student_id"`
	SelectedOptionIDs []string `json:"selected_option_ids"`
}

type AnswerQuestionUpdateDTO struct {
	SubmissionID      *string   `json:"submission_id"`
	AnswerText        *string   `json:"answer_text"`
	StudentID       *string   `json:"student_id"`
	SelectedOptionIDs *[]string `json:"selected_option_ids"`
}

type AnswerSelectedOptionCreateDTO struct {
	AnswerID string `json:"answer_id" binding:"required"`
	OptionID string `json:"option_id" binding:"required"`
}

type AnswerDocumentCreateDTO struct {
	SubmissionID *string `json:"submission_id"`
	StudentID    *string `json:"student_id"`
	DocumentID   string  `json:"document_id" binding:"required"`
	FileURL      string  `json:"file_url" binding:"required"`
	FilePath     *string `json:"file_path"`
	FileName     *string `json:"file_name"`
	FileType     *string `json:"file_type"`
	Status       *string `json:"status"`
}

type AnswerDocumentUpdateDTO struct {
	SubmissionID *string `json:"submission_id"`
	StudentID    *string `json:"student_id"`
	DocumentID   *string `json:"document_id"`
	FileURL      *string `json:"file_url"`
	FilePath     *string `json:"file_path"`
	FileName     *string `json:"file_name"`
	FileType     *string `json:"file_type"`
	Status       *string `json:"status"`
}

type QuestionOptionItemDTO struct {
	ID         string    `json:"id"`
	QuestionID string    `json:"question_id"`
	Label      string    `json:"label"`
	Value      string    `json:"value"`
	Order      int       `json:"order"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type QuestionResponseDTO struct {
	ID          string                  `json:"id"`
	BaseID      string                  `json:"base_id"`
	Text        string                  `json:"text"`
	InputType   schema.QuestionType     `json:"input_type"`
	Required    bool                    `json:"required"`
	Order       int                     `json:"order"`
	HelpText    *string                 `json:"help_text,omitempty"`
	Placeholder *string                 `json:"placeholder,omitempty"`
	MinLength   *int                    `json:"min_length,omitempty"`
	MaxLength   *int                    `json:"max_length,omitempty"`
	Active      bool                    `json:"active"`
	OptionIDs   []string                `json:"option_ids,omitempty"`
	Options     []QuestionOptionItemDTO `json:"options,omitempty"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

type QuestionOptionResponseDTO struct {
	ID         string    `json:"id"`
	QuestionID string    `json:"question_id"`
	Label      string    `json:"label"`
	Value      string    `json:"value"`
	Order      int       `json:"order"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AnswerQuestionResponseDTO struct {
	ID                string    `json:"id"`
	SubmissionID      *string   `json:"submission_id,omitempty"`
	QuestionID        string    `json:"question_id"`
	AnswerText        *string   `json:"answer_text,omitempty"`
	StudentID       *string   `json:"student_id,omitempty"`
	SelectedOptionIDs []string  `json:"selected_option_ids,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type AnswerSelectedOptionResponseDTO struct {
	AnswerID string `json:"answer_id"`
	OptionID string `json:"option_id"`
}

type AnswerDocumentResponseDTO struct {
	ID           string     `json:"id"`
	SubmissionID *string    `json:"submission_id,omitempty"`
	StudentID    string     `json:"student_id"`
	DocumentID   string     `json:"document_id"`
	FileURL      string     `json:"file_url"`
	FilePath     *string    `json:"file_path,omitempty"`
	FileName     *string    `json:"file_name,omitempty"`
	FileType     *string    `json:"file_type,omitempty"`
	Status       *string    `json:"status,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type AnswerSubmissionCreateDTO struct {
	BaseID      string     `json:"base_id" binding:"required"`
	StudentID   string     `json:"student_id" binding:"required"`
	Status      *string    `json:"status"`
	Version     *int       `json:"version"`
	SubmittedAt *time.Time `json:"submitted_at"`
}

type AnswerSubmissionUpdateDTO struct {
	BaseID      *string    `json:"base_id"`
	StudentID   *string    `json:"student_id"`
	Status      *string    `json:"status"`
	Version     *int       `json:"version"`
	SubmittedAt *time.Time `json:"submitted_at"`
}

type AnswerSubmissionResponseDTO struct {
	ID          string     `json:"id"`
	BaseID      string     `json:"base_id"`
	StudentID   string     `json:"student_id"`
	Status      string     `json:"status"`
	Version     int        `json:"version"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func NewQuestionBaseResponseDTO(base schema.QuestionBase) QuestionBaseResponseDTO {
	return QuestionBaseResponseDTO{
		ID:                       base.ID,
		Name:                     base.Name,
		Desc:                     base.Desc,
		TypeCountry:              base.TypeCountry,
		CountryID:                base.CountryID,
		AllowMultipleSubmissions: base.AllowMultipleSubmissions,
		Active:                   base.Active,
		Version:                  base.Version,
		CreatedAt:                base.CreatedAt,
		UpdatedAt:                base.UpdatedAt,
	}
}

func NewQuestionBaseResponseListDTO(bases []schema.QuestionBase) []QuestionBaseResponseDTO {
	out := make([]QuestionBaseResponseDTO, 0, len(bases))
	for _, base := range bases {
		out = append(out, NewQuestionBaseResponseDTO(base))
	}
	return out
}

func NewQuestionResponseDTO(question schema.Question) QuestionResponseDTO {
	optionIDs := make([]string, 0, len(question.Options))
	options := make([]QuestionOptionItemDTO, 0, len(question.Options))
	for _, option := range question.Options {
		optionIDs = append(optionIDs, option.ID)
		options = append(options, QuestionOptionItemDTO{
			ID:         option.ID,
			QuestionID: option.QuestionID,
			Label:      option.Label,
			Value:      option.Value,
			Order:      option.Order,
			Active:     option.Active,
			CreatedAt:  option.CreatedAt,
			UpdatedAt:  option.UpdatedAt,
		})
	}

	return QuestionResponseDTO{
		ID:          question.ID,
		BaseID:      question.BaseID,
		Text:        question.Text,
		InputType:   question.InputType,
		Required:    question.Required,
		Order:       question.Order,
		HelpText:    question.HelpText,
		Placeholder: question.Placeholder,
		MinLength:   question.MinLength,
		MaxLength:   question.MaxLength,
		Active:      question.Active,
		OptionIDs:   optionIDs,
		Options:     options,
		CreatedAt:   question.CreatedAt,
		UpdatedAt:   question.UpdatedAt,
	}
}

func NewQuestionResponseListDTO(questions []schema.Question) []QuestionResponseDTO {
	out := make([]QuestionResponseDTO, 0, len(questions))
	for _, question := range questions {
		out = append(out, NewQuestionResponseDTO(question))
	}
	return out
}

func NewQuestionOptionResponseDTO(option schema.QuestionOption) QuestionOptionResponseDTO {
	return QuestionOptionResponseDTO{
		ID:         option.ID,
		QuestionID: option.QuestionID,
		Label:      option.Label,
		Value:      option.Value,
		Order:      option.Order,
		Active:     option.Active,
		CreatedAt:  option.CreatedAt,
		UpdatedAt:  option.UpdatedAt,
	}
}

func NewQuestionOptionResponseListDTO(options []schema.QuestionOption) []QuestionOptionResponseDTO {
	out := make([]QuestionOptionResponseDTO, 0, len(options))
	for _, option := range options {
		out = append(out, NewQuestionOptionResponseDTO(option))
	}
	return out
}

func NewAnswerQuestionResponseDTO(answer schema.AnswerQuestion, selectedIDs []string) AnswerQuestionResponseDTO {
	return AnswerQuestionResponseDTO{
		ID:                answer.ID,
		SubmissionID:      answer.SubmissionID,
		QuestionID:        answer.QuestionID,
		AnswerText:        answer.AnswerText,
		StudentID:       answer.StudentID,
		SelectedOptionIDs: selectedIDs,
		CreatedAt:         answer.CreatedAt,
	}
}

func NewAnswerQuestionResponseListDTO(answers []schema.AnswerQuestion, selected map[string][]string) []AnswerQuestionResponseDTO {
	out := make([]AnswerQuestionResponseDTO, 0, len(answers))
	for _, answer := range answers {
		out = append(out, NewAnswerQuestionResponseDTO(answer, selected[answer.ID]))
	}
	return out
}

func NewAnswerSelectedOptionResponseDTO(selected schema.AnswerSelectedOption) AnswerSelectedOptionResponseDTO {
	return AnswerSelectedOptionResponseDTO{
		AnswerID: selected.AnswerID,
		OptionID: selected.OptionID,
	}
}

func NewAnswerDocumentResponseDTO(doc schema.AnswerDocument) AnswerDocumentResponseDTO {
	return AnswerDocumentResponseDTO{
		ID:           doc.ID,
		SubmissionID: doc.SubmissionID,
		StudentID:    doc.StudentID,
		DocumentID:   doc.DocumentID,
		FileURL:      doc.FileURL,
		FilePath:     doc.FilePath,
		FileName:     doc.FileName,
		FileType:     doc.FileType,
		Status:       doc.Status,
		CreatedAt:    doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
	}
}

func NewAnswerDocumentResponseListDTO(docs []schema.AnswerDocument) []AnswerDocumentResponseDTO {
	out := make([]AnswerDocumentResponseDTO, 0, len(docs))
	for _, doc := range docs {
		out = append(out, NewAnswerDocumentResponseDTO(doc))
	}
	return out
}

func NewAnswerSelectedOptionResponseListDTO(selected []schema.AnswerSelectedOption) []AnswerSelectedOptionResponseDTO {
	out := make([]AnswerSelectedOptionResponseDTO, 0, len(selected))
	for _, item := range selected {
		out = append(out, NewAnswerSelectedOptionResponseDTO(item))
	}
	return out
}

func NewAnswerSubmissionResponseDTO(submission schema.AnswerSubmission) AnswerSubmissionResponseDTO {
	return AnswerSubmissionResponseDTO{
		ID:          submission.ID,
		BaseID:      submission.BaseID,
		StudentID:   submission.StudentID,
		Status:      submission.Status,
		Version:     submission.Version,
		SubmittedAt: submission.SubmittedAt,
		CreatedAt:   submission.CreatedAt,
		UpdatedAt:   submission.UpdatedAt,
	}
}

func NewAnswerSubmissionResponseListDTO(submissions []schema.AnswerSubmission) []AnswerSubmissionResponseDTO {
	out := make([]AnswerSubmissionResponseDTO, 0, len(submissions))
	for _, submission := range submissions {
		out = append(out, NewAnswerSubmissionResponseDTO(submission))
	}
	return out
}
