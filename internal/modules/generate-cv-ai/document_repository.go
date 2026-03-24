package generatecvai

import (
	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type GeneratedDocumentRepository interface {
	GetByStudentID(studentID string) (schema.GeneratedCVAIDocument, error)
	ListByStudentID(studentID string) ([]schema.GeneratedCVAIDocument, error)
	Create(doc *schema.GeneratedCVAIDocument) error
	Update(doc *schema.GeneratedCVAIDocument) error
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

func (r *GormGeneratedDocumentRepository) GetByStudentID(studentID string) (schema.GeneratedCVAIDocument, error) {
	var doc schema.GeneratedCVAIDocument
	if err := r.baseQuery().Where("student_id = ?", studentID).First(&doc).Error; err != nil {
		return schema.GeneratedCVAIDocument{}, err
	}
	return doc, nil
}

func (r *GormGeneratedDocumentRepository) ListByStudentID(studentID string) ([]schema.GeneratedCVAIDocument, error) {
	var docs []schema.GeneratedCVAIDocument
	query := r.baseQuery().Model(&schema.GeneratedCVAIDocument{}).Order("updated_at desc")
	if studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	if err := query.Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *GormGeneratedDocumentRepository) Create(doc *schema.GeneratedCVAIDocument) error {
	return r.db.Create(doc).Error
}

func (r *GormGeneratedDocumentRepository) Update(doc *schema.GeneratedCVAIDocument) error {
	return r.db.Save(doc).Error
}
