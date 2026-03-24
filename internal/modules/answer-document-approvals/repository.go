package answerdocumentapprovals

import (
	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type Filter struct {
	StudentID        *string
	AnswerDocumentID *string
	ReviewerID       *string
	Status           *string
}

type Repository interface {
	Create(approval *schema.AnswerDocumentApproval) error
	List(filter Filter) ([]schema.AnswerDocumentApproval, error)
	GetByID(id string) (schema.AnswerDocumentApproval, error)
	GetByAnswerDocumentID(answerDocumentID string) (schema.AnswerDocumentApproval, error)
	Update(approval *schema.AnswerDocumentApproval) error
	Delete(id string) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(approval *schema.AnswerDocumentApproval) error {
	return r.db.Create(approval).Error
}

func (r *GormRepository) List(filter Filter) ([]schema.AnswerDocumentApproval, error) {
	var approvals []schema.AnswerDocumentApproval
	query := r.db.Model(&schema.AnswerDocumentApproval{})
	if filter.StudentID != nil {
		query = query.Where("student_id = ?", *filter.StudentID)
	}
	if filter.AnswerDocumentID != nil {
		query = query.Where("answer_document_id = ?", *filter.AnswerDocumentID)
	}
	if filter.ReviewerID != nil {
		query = query.Where("reviewer_id = ?", *filter.ReviewerID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if err := query.Order("updated_at desc").Find(&approvals).Error; err != nil {
		return nil, err
	}
	return approvals, nil
}

func (r *GormRepository) GetByID(id string) (schema.AnswerDocumentApproval, error) {
	var approval schema.AnswerDocumentApproval
	if err := r.db.Where("id = ?", id).First(&approval).Error; err != nil {
		return schema.AnswerDocumentApproval{}, err
	}
	return approval, nil
}

func (r *GormRepository) GetByAnswerDocumentID(answerDocumentID string) (schema.AnswerDocumentApproval, error) {
	var approval schema.AnswerDocumentApproval
	if err := r.db.Where("answer_document_id = ?", answerDocumentID).First(&approval).Error; err != nil {
		return schema.AnswerDocumentApproval{}, err
	}
	return approval, nil
}

func (r *GormRepository) Update(approval *schema.AnswerDocumentApproval) error {
	return r.db.Save(approval).Error
}

func (r *GormRepository) Delete(id string) error {
	return r.db.Delete(&schema.AnswerDocumentApproval{}, "id = ?", id).Error
}
