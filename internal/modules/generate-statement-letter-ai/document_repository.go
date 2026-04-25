package generatestatementletterai

import (
	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type GeneratedDocumentRepository interface {
	GetByID(id string) (schema.GeneratedStatementLetterAIDocument, error)
	GetByStudentID(studentID string) (schema.GeneratedStatementLetterAIDocument, error)
	ListByStudentID(studentID string) ([]schema.GeneratedStatementLetterAIDocument, error)
	Create(doc *schema.GeneratedStatementLetterAIDocument) error
	Update(doc *schema.GeneratedStatementLetterAIDocument) error
	GetCurrentApprovalByDocumentID(documentID string) (schema.StatementLetterAIApproval, error)
	CreateApproval(item *schema.StatementLetterAIApproval) error
	UpdateApproval(item *schema.StatementLetterAIApproval) error
	CreateApprovalLog(item *schema.StatementLetterAIApprovalLog) error
	GetDirector() (schema.User, error)
}

type GormGeneratedDocumentRepository struct {
	db *gorm.DB
}

func NewGeneratedDocumentRepository(db *gorm.DB) *GormGeneratedDocumentRepository {
	return &GormGeneratedDocumentRepository{db: db}
}

func (r *GormGeneratedDocumentRepository) baseQuery() *gorm.DB {
	return r.db.
		Preload("Student").
		Preload("CurrentApproval").
		Preload("CurrentApproval.Reviewer").
		Preload("CurrentApproval.Logs").
		Preload("CurrentApproval.Logs.Actor").
		Preload("Approvals").
		Preload("Approvals.Reviewer").
		Preload("Approvals.Logs").
		Preload("Approvals.Logs.Actor")
}

func (r *GormGeneratedDocumentRepository) GetByID(id string) (schema.GeneratedStatementLetterAIDocument, error) {
	var doc schema.GeneratedStatementLetterAIDocument
	if err := r.baseQuery().Where("id = ?", id).First(&doc).Error; err != nil {
		return schema.GeneratedStatementLetterAIDocument{}, err
	}
	return doc, nil
}

func (r *GormGeneratedDocumentRepository) GetByStudentID(studentID string) (schema.GeneratedStatementLetterAIDocument, error) {
	var doc schema.GeneratedStatementLetterAIDocument
	if err := r.baseQuery().Where("student_id = ?", studentID).First(&doc).Error; err != nil {
		return schema.GeneratedStatementLetterAIDocument{}, err
	}
	return doc, nil
}

func (r *GormGeneratedDocumentRepository) ListByStudentID(studentID string) ([]schema.GeneratedStatementLetterAIDocument, error) {
	var docs []schema.GeneratedStatementLetterAIDocument
	query := r.baseQuery().Model(&schema.GeneratedStatementLetterAIDocument{}).Order("updated_at desc")
	if studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	if err := query.Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *GormGeneratedDocumentRepository) Create(doc *schema.GeneratedStatementLetterAIDocument) error {
	return r.db.Create(doc).Error
}

func (r *GormGeneratedDocumentRepository) Update(doc *schema.GeneratedStatementLetterAIDocument) error {
	return r.db.Save(doc).Error
}

func (r *GormGeneratedDocumentRepository) GetCurrentApprovalByDocumentID(documentID string) (schema.StatementLetterAIApproval, error) {
	var item schema.StatementLetterAIApproval
	if err := r.db.
		Preload("Reviewer").
		Preload("Logs").
		Preload("Logs.Actor").
		Where("document_id = ?", documentID).
		Order("updated_at desc").
		First(&item).Error; err != nil {
		return schema.StatementLetterAIApproval{}, err
	}
	return item, nil
}

func (r *GormGeneratedDocumentRepository) CreateApproval(item *schema.StatementLetterAIApproval) error {
	return r.db.Create(item).Error
}

func (r *GormGeneratedDocumentRepository) UpdateApproval(item *schema.StatementLetterAIApproval) error {
	return r.db.Save(item).Error
}

func (r *GormGeneratedDocumentRepository) CreateApprovalLog(item *schema.StatementLetterAIApprovalLog) error {
	return r.db.Create(item).Error
}

func (r *GormGeneratedDocumentRepository) GetDirector() (schema.User, error) {
	var user schema.User
	if err := r.db.Where("role = ?", schema.UserRoleDirector).Order("created_at asc").First(&user).Error; err != nil {
		return schema.User{}, err
	}
	return user, nil
}
