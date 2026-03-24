package generatestatementletterai

import (
	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type GeneratedDocumentRepository interface {
	GetByStudentID(studentID string) (schema.GeneratedStatementLetterAIDocument, error)
	ListByStudentID(studentID string) ([]schema.GeneratedStatementLetterAIDocument, error)
	Create(doc *schema.GeneratedStatementLetterAIDocument) error
	Update(doc *schema.GeneratedStatementLetterAIDocument) error
}

type GormGeneratedDocumentRepository struct {
	db *gorm.DB
}

func NewGeneratedDocumentRepository(db *gorm.DB) *GormGeneratedDocumentRepository {
	return &GormGeneratedDocumentRepository{db: db}
}

func (r *GormGeneratedDocumentRepository) baseQuery() *gorm.DB {
	return r.db.Preload("Student")
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
