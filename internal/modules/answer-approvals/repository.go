package answerapprovals

import (
	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type Filter struct {
	StudentID  *string
	AnswerID   *string
	ReviewerID *string
	Status     *string
}

type Repository interface {
	Create(approval *schema.AnswerApproval) error
	List(filter Filter) ([]schema.AnswerApproval, error)
	GetByID(id string) (schema.AnswerApproval, error)
	GetByAnswerID(answerID string) (schema.AnswerApproval, error)
	Update(approval *schema.AnswerApproval) error
	Delete(id string) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(approval *schema.AnswerApproval) error {
	return r.db.Create(approval).Error
}

func (r *GormRepository) List(filter Filter) ([]schema.AnswerApproval, error) {
	var approvals []schema.AnswerApproval
	query := r.db.Model(&schema.AnswerApproval{})
	if filter.StudentID != nil {
		query = query.Where("student_id = ?", *filter.StudentID)
	}
	if filter.AnswerID != nil {
		query = query.Where("answer_id = ?", *filter.AnswerID)
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

func (r *GormRepository) GetByID(id string) (schema.AnswerApproval, error) {
	var approval schema.AnswerApproval
	if err := r.db.Where("id = ?", id).First(&approval).Error; err != nil {
		return schema.AnswerApproval{}, err
	}
	return approval, nil
}

func (r *GormRepository) GetByAnswerID(answerID string) (schema.AnswerApproval, error) {
	var approval schema.AnswerApproval
	if err := r.db.Where("answer_id = ?", answerID).First(&approval).Error; err != nil {
		return schema.AnswerApproval{}, err
	}
	return approval, nil
}

func (r *GormRepository) Update(approval *schema.AnswerApproval) error {
	return r.db.Save(approval).Error
}

func (r *GormRepository) Delete(id string) error {
	return r.db.Delete(&schema.AnswerApproval{}, "id = ?", id).Error
}
