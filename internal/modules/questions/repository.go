package questions

import (
	"strings"

	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type QuestionBaseFilter struct {
	TypeCountry *string
	CountryID   *string
	Active      *bool
}

type QuestionFilter struct {
	BaseID *string
	Active *bool
}

type QuestionOptionFilter struct {
	QuestionID *string
	Active     *bool
}

type AnswerQuestionFilter struct {
	SubmissionID *string
	QuestionID   *string
	UserID  *string
}

type AnswerSubmissionFilter struct {
	BaseID    *string
	StudentID *string
	Status    *string
}

type AnswerSelectedOptionFilter struct {
	AnswerID  *string
	AnswerIDs []string
	OptionID  *string
}

type AnswerDocumentFilter struct {
	SubmissionID *string
	StudentID    *string
	DocumentID   *string
}

type Repository interface {
	CreateQuestionBase(base *schema.QuestionBase) error
	ListQuestionBases(filter QuestionBaseFilter) ([]schema.QuestionBase, error)
	GetQuestionBaseByID(id string) (schema.QuestionBase, error)
	UpdateQuestionBase(base *schema.QuestionBase) error
	DeleteQuestionBase(id string) error

	CreateQuestion(question *schema.Question) error
	ListQuestions(filter QuestionFilter) ([]schema.Question, error)
	GetQuestionByID(id string) (schema.Question, error)
	UpdateQuestion(question *schema.Question) error
	DeleteQuestion(id string) error

	CreateQuestionOption(option *schema.QuestionOption) error
	ListQuestionOptions(filter QuestionOptionFilter) ([]schema.QuestionOption, error)
	GetQuestionOptionByID(id string) (schema.QuestionOption, error)
	UpdateQuestionOption(option *schema.QuestionOption) error
	DeleteQuestionOption(id string) error

	CreateAnswerQuestion(answer *schema.AnswerQuestion) error
	ListAnswerQuestions(filter AnswerQuestionFilter) ([]schema.AnswerQuestion, error)
	GetAnswerQuestionByID(id string) (schema.AnswerQuestion, error)
	UpdateAnswerQuestion(answer *schema.AnswerQuestion) error
	DeleteAnswerQuestion(id string) error
	ReplaceAnswerSelectedOptions(answerID string, optionIDs []string) error

	CreateAnswerSubmission(submission *schema.AnswerSubmission) error
	ListAnswerSubmissions(filter AnswerSubmissionFilter) ([]schema.AnswerSubmission, error)
	GetAnswerSubmissionByID(id string) (schema.AnswerSubmission, error)
	UpdateAnswerSubmission(submission *schema.AnswerSubmission) error
	DeleteAnswerSubmission(id string) error

	CreateAnswerSelectedOption(selected *schema.AnswerSelectedOption) error
	ListAnswerSelectedOptions(filter AnswerSelectedOptionFilter) ([]schema.AnswerSelectedOption, error)
	GetAnswerSelectedOption(answerID, optionID string) (schema.AnswerSelectedOption, error)
	DeleteAnswerSelectedOption(answerID, optionID string) error

	CreateAnswerDocuments(docs []schema.AnswerDocument) error
	ListAnswerDocuments(filter AnswerDocumentFilter) ([]schema.AnswerDocument, error)
	GetAnswerDocumentByID(id string) (schema.AnswerDocument, error)
	UpdateAnswerDocument(doc *schema.AnswerDocument) error
	DeleteAnswerDocument(id string) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) CreateQuestionBase(base *schema.QuestionBase) error {
	return r.db.Create(base).Error
}

func (r *GormRepository) ListQuestionBases(filter QuestionBaseFilter) ([]schema.QuestionBase, error) {
	var bases []schema.QuestionBase
	query := r.db.Model(&schema.QuestionBase{})
	if filter.TypeCountry != nil && strings.TrimSpace(*filter.TypeCountry) != "" {
		query = query.Where("type_country = ?", *filter.TypeCountry)
	}
	if filter.CountryID != nil && strings.TrimSpace(*filter.CountryID) != "" {
		query = query.Where("country_id = ?", *filter.CountryID)
	}
	if filter.Active != nil {
		query = query.Where("active = ?", *filter.Active)
	}
	if err := query.Order("id desc").Find(&bases).Error; err != nil {
		return nil, err
	}
	return bases, nil
}

func (r *GormRepository) GetQuestionBaseByID(id string) (schema.QuestionBase, error) {
	var base schema.QuestionBase
	if err := r.db.First(&base, "id = ?", id).Error; err != nil {
		return schema.QuestionBase{}, err
	}
	return base, nil
}

func (r *GormRepository) UpdateQuestionBase(base *schema.QuestionBase) error {
	return r.db.Save(base).Error
}

func (r *GormRepository) DeleteQuestionBase(id string) error {
	return r.db.Delete(&schema.QuestionBase{}, "id = ?", id).Error
}

func (r *GormRepository) CreateQuestion(question *schema.Question) error {
	return r.db.Create(question).Error
}

func (r *GormRepository) ListQuestions(filter QuestionFilter) ([]schema.Question, error) {
	var questions []schema.Question
	query := r.db.Preload("Options", func(db *gorm.DB) *gorm.DB {
		return db.Order("`order` asc")
	})
	if filter.BaseID != nil && strings.TrimSpace(*filter.BaseID) != "" {
		query = query.Where("base_id = ?", *filter.BaseID)
	}
	if filter.Active != nil {
		query = query.Where("active = ?", *filter.Active)
	}
	if err := query.Order("`order` asc").Order("id desc").Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *GormRepository) GetQuestionByID(id string) (schema.Question, error) {
	var question schema.Question
	if err := r.db.Preload("Options", func(db *gorm.DB) *gorm.DB {
		return db.Order("`order` asc")
	}).First(&question, "id = ?", id).Error; err != nil {
		return schema.Question{}, err
	}
	return question, nil
}

func (r *GormRepository) UpdateQuestion(question *schema.Question) error {
	return r.db.Save(question).Error
}

func (r *GormRepository) DeleteQuestion(id string) error {
	return r.db.Delete(&schema.Question{}, "id = ?", id).Error
}

func (r *GormRepository) CreateQuestionOption(option *schema.QuestionOption) error {
	return r.db.Create(option).Error
}

func (r *GormRepository) ListQuestionOptions(filter QuestionOptionFilter) ([]schema.QuestionOption, error) {
	var options []schema.QuestionOption
	query := r.db.Model(&schema.QuestionOption{})
	if filter.QuestionID != nil && strings.TrimSpace(*filter.QuestionID) != "" {
		query = query.Where("question_id = ?", *filter.QuestionID)
	}
	if filter.Active != nil {
		query = query.Where("active = ?", *filter.Active)
	}
	if err := query.Order("`order` asc").Order("id desc").Find(&options).Error; err != nil {
		return nil, err
	}
	return options, nil
}

func (r *GormRepository) GetQuestionOptionByID(id string) (schema.QuestionOption, error) {
	var option schema.QuestionOption
	if err := r.db.First(&option, "id = ?", id).Error; err != nil {
		return schema.QuestionOption{}, err
	}
	return option, nil
}

func (r *GormRepository) UpdateQuestionOption(option *schema.QuestionOption) error {
	return r.db.Save(option).Error
}

func (r *GormRepository) DeleteQuestionOption(id string) error {
	return r.db.Delete(&schema.QuestionOption{}, "id = ?", id).Error
}

func (r *GormRepository) CreateAnswerQuestion(answer *schema.AnswerQuestion) error {
	return r.db.Create(answer).Error
}

func (r *GormRepository) ListAnswerQuestions(filter AnswerQuestionFilter) ([]schema.AnswerQuestion, error) {
	var answers []schema.AnswerQuestion
	query := r.db.Model(&schema.AnswerQuestion{})
	if filter.SubmissionID != nil && strings.TrimSpace(*filter.SubmissionID) != "" {
		query = query.Where("submission_id = ?", *filter.SubmissionID)
	}
	if filter.QuestionID != nil && strings.TrimSpace(*filter.QuestionID) != "" {
		query = query.Where("question_id = ?", *filter.QuestionID)
	}
	if filter.UserID != nil && strings.TrimSpace(*filter.UserID) != "" {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if err := query.Order("created_at desc").Find(&answers).Error; err != nil {
		return nil, err
	}
	return answers, nil
}

func (r *GormRepository) GetAnswerQuestionByID(id string) (schema.AnswerQuestion, error) {
	var answer schema.AnswerQuestion
	if err := r.db.First(&answer, "id = ?", id).Error; err != nil {
		return schema.AnswerQuestion{}, err
	}
	return answer, nil
}

func (r *GormRepository) UpdateAnswerQuestion(answer *schema.AnswerQuestion) error {
	return r.db.Save(answer).Error
}

func (r *GormRepository) DeleteAnswerQuestion(id string) error {
	return r.db.Delete(&schema.AnswerQuestion{}, "id = ?", id).Error
}

func (r *GormRepository) CreateAnswerSubmission(submission *schema.AnswerSubmission) error {
	return r.db.Create(submission).Error
}

func (r *GormRepository) ListAnswerSubmissions(filter AnswerSubmissionFilter) ([]schema.AnswerSubmission, error) {
	var submissions []schema.AnswerSubmission
	query := r.db.Model(&schema.AnswerSubmission{})
	if filter.BaseID != nil && strings.TrimSpace(*filter.BaseID) != "" {
		query = query.Where("base_id = ?", *filter.BaseID)
	}
	if filter.StudentID != nil && strings.TrimSpace(*filter.StudentID) != "" {
		query = query.Where("student_id = ?", *filter.StudentID)
	}
	if filter.Status != nil && strings.TrimSpace(*filter.Status) != "" {
		query = query.Where("status = ?", *filter.Status)
	}
	if err := query.Order("created_at desc").Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *GormRepository) GetAnswerSubmissionByID(id string) (schema.AnswerSubmission, error) {
	var submission schema.AnswerSubmission
	if err := r.db.First(&submission, "id = ?", id).Error; err != nil {
		return schema.AnswerSubmission{}, err
	}
	return submission, nil
}

func (r *GormRepository) UpdateAnswerSubmission(submission *schema.AnswerSubmission) error {
	return r.db.Save(submission).Error
}

func (r *GormRepository) DeleteAnswerSubmission(id string) error {
	return r.db.Delete(&schema.AnswerSubmission{}, "id = ?", id).Error
}

func (r *GormRepository) ReplaceAnswerSelectedOptions(answerID string, optionIDs []string) error {
	answerID = strings.TrimSpace(answerID)
	if answerID == "" {
		return nil
	}
	unique := make(map[string]struct{}, len(optionIDs))
	selected := make([]schema.AnswerSelectedOption, 0, len(optionIDs))
	for _, id := range optionIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, ok := unique[id]; ok {
			continue
		}
		unique[id] = struct{}{}
		selected = append(selected, schema.AnswerSelectedOption{
			AnswerID: answerID,
			OptionID: id,
		})
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("answer_id = ?", answerID).Delete(&schema.AnswerSelectedOption{}).Error; err != nil {
			return err
		}
		if len(selected) == 0 {
			return nil
		}
		return tx.Create(&selected).Error
	})
}

func (r *GormRepository) CreateAnswerSelectedOption(selected *schema.AnswerSelectedOption) error {
	return r.db.Create(selected).Error
}

func (r *GormRepository) ListAnswerSelectedOptions(filter AnswerSelectedOptionFilter) ([]schema.AnswerSelectedOption, error) {
	var selected []schema.AnswerSelectedOption
	query := r.db.Model(&schema.AnswerSelectedOption{})
	if len(filter.AnswerIDs) > 0 {
		query = query.Where("answer_id IN ?", filter.AnswerIDs)
	}
	if filter.AnswerID != nil && strings.TrimSpace(*filter.AnswerID) != "" {
		query = query.Where("answer_id = ?", *filter.AnswerID)
	}
	if filter.OptionID != nil && strings.TrimSpace(*filter.OptionID) != "" {
		query = query.Where("option_id = ?", *filter.OptionID)
	}
	if err := query.Order("answer_id asc").Order("option_id asc").Find(&selected).Error; err != nil {
		return nil, err
	}
	return selected, nil
}

func (r *GormRepository) GetAnswerSelectedOption(answerID, optionID string) (schema.AnswerSelectedOption, error) {
	var selected schema.AnswerSelectedOption
	if err := r.db.First(&selected, "answer_id = ? AND option_id = ?", answerID, optionID).Error; err != nil {
		return schema.AnswerSelectedOption{}, err
	}
	return selected, nil
}

func (r *GormRepository) DeleteAnswerSelectedOption(answerID, optionID string) error {
	return r.db.Delete(&schema.AnswerSelectedOption{}, "answer_id = ? AND option_id = ?", answerID, optionID).Error
}

func (r *GormRepository) CreateAnswerDocuments(docs []schema.AnswerDocument) error {
	if len(docs) == 0 {
		return nil
	}
	return r.db.Create(&docs).Error
}

func (r *GormRepository) ListAnswerDocuments(filter AnswerDocumentFilter) ([]schema.AnswerDocument, error) {
	var docs []schema.AnswerDocument
	query := r.db.Model(&schema.AnswerDocument{})
	if filter.SubmissionID != nil && strings.TrimSpace(*filter.SubmissionID) != "" {
		query = query.Where("submission_id = ?", *filter.SubmissionID)
	}
	if filter.StudentID != nil && strings.TrimSpace(*filter.StudentID) != "" {
		query = query.Where("student_id = ?", *filter.StudentID)
	}
	if filter.DocumentID != nil && strings.TrimSpace(*filter.DocumentID) != "" {
		query = query.Where("document_id = ?", *filter.DocumentID)
	}
	if err := query.Order("id desc").Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *GormRepository) GetAnswerDocumentByID(id string) (schema.AnswerDocument, error) {
	var doc schema.AnswerDocument
	if err := r.db.First(&doc, "id = ?", id).Error; err != nil {
		return schema.AnswerDocument{}, err
	}
	return doc, nil
}

func (r *GormRepository) UpdateAnswerDocument(doc *schema.AnswerDocument) error {
	return r.db.Save(doc).Error
}

func (r *GormRepository) DeleteAnswerDocument(id string) error {
	return r.db.Delete(&schema.AnswerDocument{}, "id = ?", id).Error
}
