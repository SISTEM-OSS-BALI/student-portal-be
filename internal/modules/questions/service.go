package questions

import (
	"errors"
	"strings"

	"github.com/username/gin-gorm-api/internal/schema"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateQuestionBase(input QuestionBaseCreateDTO) (schema.QuestionBase, error) {
	base := schema.QuestionBase{
		Name:        input.Name,
		Desc:        input.Desc,
		TypeCountry: input.TypeCountry,
		CountryID:   input.CountryID,
	}
	if input.AllowMultipleSubmissions != nil {
		base.AllowMultipleSubmissions = *input.AllowMultipleSubmissions
	}
	if input.Active != nil {
		base.Active = *input.Active
	}
	if input.Version != nil {
		base.Version = *input.Version
	}
	if err := s.repo.CreateQuestionBase(&base); err != nil {
		return schema.QuestionBase{}, err
	}
	return base, nil
}

func (s *Service) ListQuestionBases(filter QuestionBaseFilter) ([]schema.QuestionBase, error) {
	return s.repo.ListQuestionBases(filter)
}

func (s *Service) GetQuestionBaseByID(id string) (schema.QuestionBase, error) {
	return s.repo.GetQuestionBaseByID(id)
}

func (s *Service) UpdateQuestionBase(id string, input QuestionBaseUpdateDTO) (schema.QuestionBase, error) {
	base, err := s.repo.GetQuestionBaseByID(id)
	if err != nil {
		return schema.QuestionBase{}, err
	}
	if input.Name != nil {
		base.Name = *input.Name
	}
	if input.Desc != nil {
		base.Desc = input.Desc
	}
	if input.TypeCountry != nil {
		base.TypeCountry = *input.TypeCountry
	}
	if input.CountryID != nil {
		base.CountryID = input.CountryID
	}
	if input.AllowMultipleSubmissions != nil {
		base.AllowMultipleSubmissions = *input.AllowMultipleSubmissions
	}
	if input.Active != nil {
		base.Active = *input.Active
	}
	if input.Version != nil {
		base.Version = *input.Version
	}
	if err := s.repo.UpdateQuestionBase(&base); err != nil {
		return schema.QuestionBase{}, err
	}
	return base, nil
}

func (s *Service) DeleteQuestionBase(id string) error {
	return s.repo.DeleteQuestionBase(id)
}

func (s *Service) CreateQuestion(input QuestionCreateDTO) (schema.Question, error) {
	question := schema.Question{
		BaseID:      input.BaseID,
		Text:        input.Text,
		InputType:   input.InputType,
		HelpText:    input.HelpText,
		Placeholder: input.Placeholder,
		MinLength:   input.MinLength,
		MaxLength:   input.MaxLength,
	}
	if input.Required != nil {
		question.Required = *input.Required
	}
	if input.Order != nil {
		question.Order = *input.Order
	}
	if input.Active != nil {
		question.Active = *input.Active
	}
	if err := s.repo.CreateQuestion(&question); err != nil {
		return schema.Question{}, err
	}
	return s.repo.GetQuestionByID(question.ID)
}

func (s *Service) ListQuestions(filter QuestionFilter) ([]schema.Question, error) {
	return s.repo.ListQuestions(filter)
}

func (s *Service) GetQuestionByID(id string) (schema.Question, error) {
	return s.repo.GetQuestionByID(id)
}

func (s *Service) UpdateQuestion(id string, input QuestionUpdateDTO) (schema.Question, error) {
	question, err := s.repo.GetQuestionByID(id)
	if err != nil {
		return schema.Question{}, err
	}
	if input.BaseID != nil {
		question.BaseID = *input.BaseID
	}
	if input.Text != nil {
		question.Text = *input.Text
	}
	if input.InputType != nil {
		question.InputType = *input.InputType
	}
	if input.Required != nil {
		question.Required = *input.Required
	}
	if input.Order != nil {
		question.Order = *input.Order
	}
	if input.HelpText != nil {
		question.HelpText = input.HelpText
	}
	if input.Placeholder != nil {
		question.Placeholder = input.Placeholder
	}
	if input.MinLength != nil {
		question.MinLength = input.MinLength
	}
	if input.MaxLength != nil {
		question.MaxLength = input.MaxLength
	}
	if input.Active != nil {
		question.Active = *input.Active
	}
	if err := s.repo.UpdateQuestion(&question); err != nil {
		return schema.Question{}, err
	}
	return s.repo.GetQuestionByID(question.ID)
}

func (s *Service) DeleteQuestion(id string) error {
	return s.repo.DeleteQuestion(id)
}

func (s *Service) CreateQuestionOption(input QuestionOptionCreateDTO) (schema.QuestionOption, error) {
	option := schema.QuestionOption{
		QuestionID: input.QuestionID,
		Label:      input.Label,
		Value:      input.Value,
	}
	if input.Order != nil {
		option.Order = *input.Order
	}
	if input.Active != nil {
		option.Active = *input.Active
	}
	if err := s.repo.CreateQuestionOption(&option); err != nil {
		return schema.QuestionOption{}, err
	}
	return option, nil
}

func (s *Service) ListQuestionOptions(filter QuestionOptionFilter) ([]schema.QuestionOption, error) {
	return s.repo.ListQuestionOptions(filter)
}

func (s *Service) GetQuestionOptionByID(id string) (schema.QuestionOption, error) {
	return s.repo.GetQuestionOptionByID(id)
}

func (s *Service) UpdateQuestionOption(id string, input QuestionOptionUpdateDTO) (schema.QuestionOption, error) {
	option, err := s.repo.GetQuestionOptionByID(id)
	if err != nil {
		return schema.QuestionOption{}, err
	}
	if input.QuestionID != nil {
		option.QuestionID = *input.QuestionID
	}
	if input.Label != nil {
		option.Label = *input.Label
	}
	if input.Value != nil {
		option.Value = *input.Value
	}
	if input.Order != nil {
		option.Order = *input.Order
	}
	if input.Active != nil {
		option.Active = *input.Active
	}
	if err := s.repo.UpdateQuestionOption(&option); err != nil {
		return schema.QuestionOption{}, err
	}
	return option, nil
}

func (s *Service) DeleteQuestionOption(id string) error {
	return s.repo.DeleteQuestionOption(id)
}

func (s *Service) CreateAnswerQuestion(input AnswerQuestionCreateDTO) (schema.AnswerQuestion, error) {
	answer := schema.AnswerQuestion{
		SubmissionID: input.SubmissionID,
		QuestionID:   input.QuestionID,
		AnswerText:   input.AnswerText,
		StudentID:  input.StudentID,
	}
	if err := s.repo.CreateAnswerQuestion(&answer); err != nil {
		return schema.AnswerQuestion{}, err
	}
	if len(input.SelectedOptionIDs) > 0 {
		if err := s.repo.ReplaceAnswerSelectedOptions(answer.ID, input.SelectedOptionIDs); err != nil {
			return schema.AnswerQuestion{}, err
		}
	}
	return s.repo.GetAnswerQuestionByID(answer.ID)
}

func (s *Service) ListAnswerQuestions(filter AnswerQuestionFilter) ([]schema.AnswerQuestion, error) {
	return s.repo.ListAnswerQuestions(filter)
}

func (s *Service) GetAnswerQuestionByID(id string) (schema.AnswerQuestion, error) {
	return s.repo.GetAnswerQuestionByID(id)
}

func (s *Service) UpdateAnswerQuestion(id string, input AnswerQuestionUpdateDTO) (schema.AnswerQuestion, error) {
	answer, err := s.repo.GetAnswerQuestionByID(id)
	if err != nil {
		return schema.AnswerQuestion{}, err
	}
	if input.SubmissionID != nil {
		answer.SubmissionID = input.SubmissionID
	}
	if input.AnswerText != nil {
		answer.AnswerText = input.AnswerText
	}
	if input.StudentID != nil {
		answer.StudentID = input.StudentID
	}
	if err := s.repo.UpdateAnswerQuestion(&answer); err != nil {
		return schema.AnswerQuestion{}, err
	}
	if input.SelectedOptionIDs != nil {
		if err := s.repo.ReplaceAnswerSelectedOptions(answer.ID, *input.SelectedOptionIDs); err != nil {
			return schema.AnswerQuestion{}, err
		}
	}
	return s.repo.GetAnswerQuestionByID(answer.ID)
}

func (s *Service) DeleteAnswerQuestion(id string) error {
	return s.repo.DeleteAnswerQuestion(id)
}

func (s *Service) CreateAnswerSubmission(input AnswerSubmissionCreateDTO) (schema.AnswerSubmission, error) {
	submission := schema.AnswerSubmission{
		BaseID:      input.BaseID,
		StudentID:   input.StudentID,
		SubmittedAt: input.SubmittedAt,
	}
	if input.Status != nil {
		submission.Status = *input.Status
	}
	if input.Version != nil {
		submission.Version = *input.Version
	}
	if err := s.repo.CreateAnswerSubmission(&submission); err != nil {
		return schema.AnswerSubmission{}, err
	}
	return submission, nil
}

func (s *Service) ListAnswerSubmissions(filter AnswerSubmissionFilter) ([]schema.AnswerSubmission, error) {
	return s.repo.ListAnswerSubmissions(filter)
}

func (s *Service) GetAnswerSubmissionByID(id string) (schema.AnswerSubmission, error) {
	return s.repo.GetAnswerSubmissionByID(id)
}

func (s *Service) UpdateAnswerSubmission(id string, input AnswerSubmissionUpdateDTO) (schema.AnswerSubmission, error) {
	submission, err := s.repo.GetAnswerSubmissionByID(id)
	if err != nil {
		return schema.AnswerSubmission{}, err
	}
	if input.BaseID != nil {
		submission.BaseID = *input.BaseID
	}
	if input.StudentID != nil {
		submission.StudentID = *input.StudentID
	}
	if input.Status != nil {
		submission.Status = *input.Status
	}
	if input.Version != nil {
		submission.Version = *input.Version
	}
	if input.SubmittedAt != nil {
		submission.SubmittedAt = input.SubmittedAt
	}
	if err := s.repo.UpdateAnswerSubmission(&submission); err != nil {
		return schema.AnswerSubmission{}, err
	}
	return submission, nil
}

func (s *Service) DeleteAnswerSubmission(id string) error {
	return s.repo.DeleteAnswerSubmission(id)
}

func (s *Service) CreateAnswerSelectedOption(input AnswerSelectedOptionCreateDTO) (schema.AnswerSelectedOption, error) {
	selected := schema.AnswerSelectedOption{
		AnswerID: input.AnswerID,
		OptionID: input.OptionID,
	}
	if err := s.repo.CreateAnswerSelectedOption(&selected); err != nil {
		return schema.AnswerSelectedOption{}, err
	}
	return selected, nil
}

func (s *Service) ListAnswerSelectedOptions(filter AnswerSelectedOptionFilter) ([]schema.AnswerSelectedOption, error) {
	return s.repo.ListAnswerSelectedOptions(filter)
}

func (s *Service) GetAnswerSelectedOption(answerID, optionID string) (schema.AnswerSelectedOption, error) {
	return s.repo.GetAnswerSelectedOption(answerID, optionID)
}

func (s *Service) DeleteAnswerSelectedOption(answerID, optionID string) error {
	return s.repo.DeleteAnswerSelectedOption(answerID, optionID)
}

func (s *Service) CreateAnswerDocuments(inputs []AnswerDocumentCreateDTO) ([]schema.AnswerDocument, error) {
	docs := make([]schema.AnswerDocument, 0, len(inputs))
	for _, input := range inputs {
		if input.StudentID == nil || strings.TrimSpace(*input.StudentID) == "" {
			return nil, errors.New("student_id is required")
		}
		doc := schema.AnswerDocument{
			SubmissionID: input.SubmissionID,
			StudentID:    *input.StudentID,
			DocumentID:   input.DocumentID,
			FileURL:      input.FileURL,
			FilePath:     input.FilePath,
			FileName:     input.FileName,
			FileType:     input.FileType,
			Status:       input.Status,
		}
		docs = append(docs, doc)
	}
	if err := s.repo.CreateAnswerDocuments(docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func (s *Service) ListAnswerDocuments(filter AnswerDocumentFilter) ([]schema.AnswerDocument, error) {
	return s.repo.ListAnswerDocuments(filter)
}

func (s *Service) GetAnswerDocumentByID(id string) (schema.AnswerDocument, error) {
	return s.repo.GetAnswerDocumentByID(id)
}

func (s *Service) UpdateAnswerDocument(id string, input AnswerDocumentUpdateDTO) (schema.AnswerDocument, error) {
	doc, err := s.repo.GetAnswerDocumentByID(id)
	if err != nil {
		return schema.AnswerDocument{}, err
	}
	if input.SubmissionID != nil {
		doc.SubmissionID = input.SubmissionID
	}
	if input.StudentID != nil {
		doc.StudentID = *input.StudentID
	}
	if input.DocumentID != nil {
		doc.DocumentID = *input.DocumentID
	}
	if input.FileURL != nil {
		doc.FileURL = *input.FileURL
	}
	if input.FilePath != nil {
		doc.FilePath = input.FilePath
	}
	if input.FileName != nil {
		doc.FileName = input.FileName
	}
	if input.FileType != nil {
		doc.FileType = input.FileType
	}
	if input.Status != nil {
		doc.Status = input.Status
	}
	if err := s.repo.UpdateAnswerDocument(&doc); err != nil {
		return schema.AnswerDocument{}, err
	}
	return doc, nil
}

func (s *Service) DeleteAnswerDocument(id string) error {
	return s.repo.DeleteAnswerDocument(id)
}
