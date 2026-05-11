package documents

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

func (s *Service) Create(input CreateDTO) (schema.DocumentsManagement, error) {
	translation, err := parseTranslationNeeded(input.TranslationNeeded)
	if err != nil {
		return schema.DocumentsManagement{}, err
	}
	if input.Required == nil {
		return schema.DocumentsManagement{}, errors.New("required field is missing")
	}

	autoRename, err := parseAutoRenamePattern(input.AutoRenamePattern, input.Label)
	if err != nil {
		return schema.DocumentsManagement{}, err
	}

	doc := schema.DocumentsManagement{
		Label:             input.Label,
		InternalCode:      input.InternalCode,
		FileType:          input.FileType,
		Category:          input.Category,
		ExampleURL:        input.ExampleURL,
		TranslationNeeded: translation,
		Required:          *input.Required,
		AutoRenamePattern: autoRename,
		Notes:             input.Notes,
	}
	if err := s.repo.Create(&doc); err != nil {
		return schema.DocumentsManagement{}, err
	}
	return doc, nil
}

func (s *Service) List() ([]schema.DocumentsManagement, error) {
	return s.repo.List()
}

func (s *Service) GetByID(id string) (schema.DocumentsManagement, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(id string, input UpdateDTO) (schema.DocumentsManagement, error) {
	doc, err := s.repo.GetByID(id)
	if err != nil {
		return schema.DocumentsManagement{}, err
	}

	if input.Label != nil {
		doc.Label = *input.Label
	}
	if input.InternalCode != nil {
		doc.InternalCode = *input.InternalCode
	}
	if input.FileType != nil {
		doc.FileType = *input.FileType
	}
	if input.Category != nil {
		doc.Category = *input.Category
	}
	if input.ExampleURL != nil {
		doc.ExampleURL = input.ExampleURL
	}
	if input.TranslationNeeded != nil {
		translation, err := parseTranslationNeeded(*input.TranslationNeeded)
		if err != nil {
			return schema.DocumentsManagement{}, err
		}
		doc.TranslationNeeded = translation
	}
	if input.Required != nil {
		doc.Required = *input.Required
	}
	if input.AutoRenamePattern != nil {
		autoRename, err := parseAutoRenamePattern(*input.AutoRenamePattern, doc.Label)
		if err != nil {
			return schema.DocumentsManagement{}, err
		}
		doc.AutoRenamePattern = autoRename
	}
	if input.Notes != nil {
		doc.Notes = *input.Notes
	}

	if err := s.repo.Update(&doc); err != nil {
		return schema.DocumentsManagement{}, err
	}
	return doc, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

func parseTranslationNeeded(value string) (schema.TranslationNeeded, error) {
	v := strings.ToUpper(strings.TrimSpace(value))
	switch v {
	case string(schema.TranslationNeededYes):
		return schema.TranslationNeededYes, nil
	case string(schema.TranslationNeededNo):
		return schema.TranslationNeededNo, nil
	default:
		return "", errors.New("translation_needed must be YES or NO")
	}
}

func parseAutoRenamePattern(value schema.AutoRenamePattern, label string) (schema.AutoRenamePattern, error) {
	v := schema.AutoRenamePattern(strings.TrimSpace(string(value)))
	if v == "" {
		return schema.AutoRenamePatternNone, nil
	}

	switch v {
	case schema.AutoRenamePatternNone,
		schema.AutoRenamePatternDate,
		schema.AutoRenamePatternDocumentID:
		return v, nil
	default:
		return parseLabelBasedPattern(v, label)
	}
}

func parseLabelBasedPattern(value schema.AutoRenamePattern, label string) (schema.AutoRenamePattern, error) {
	token := normalizeLabelToken(label)
	if token == "" {
		return "", errors.New("label is required for label-based auto_rename_pattern")
	}

	expected1 := "{studentName}_" + token + ".pdf"
	expected2 := token + "_{studentName}.pdf"

	raw := strings.TrimSpace(string(value))
	if strings.EqualFold(raw, expected1) {
		return schema.AutoRenamePattern(expected1), nil
	}
	if strings.EqualFold(raw, expected2) {
		return schema.AutoRenamePattern(expected2), nil
	}

	return "", errors.New("auto_rename_pattern must match label-based options")
}

func normalizeLabelToken(label string) string {
	normalized := strings.ToUpper(strings.TrimSpace(label))
	if normalized == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(normalized))
	for _, r := range normalized {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (s *Service) DocumentTranslationRequired() ([]schema.DocumentsManagement, error) {
	return s.repo.DocumentTranslationRequired()
}
