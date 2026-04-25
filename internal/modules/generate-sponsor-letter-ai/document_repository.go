package generatesponsorletterai

import (
	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type GeneratedDocumentRepository interface {
	GetByID(id string) (schema.GeneratedSponsorLetterAIDocument, error)
	GetByStudentID(studentID string) (schema.GeneratedSponsorLetterAIDocument, error)
	ListByStudentID(studentID string) ([]schema.GeneratedSponsorLetterAIDocument, error)
	Create(doc *schema.GeneratedSponsorLetterAIDocument) error
	Update(doc *schema.GeneratedSponsorLetterAIDocument) error
	GetCurrentApprovalByDocumentID(documentID string) (schema.SponsorLetterAIApproval, error)
	CreateApproval(item *schema.SponsorLetterAIApproval) error
	UpdateApproval(item *schema.SponsorLetterAIApproval) error
	CreateApprovalLog(item *schema.SponsorLetterAIApprovalLog) error
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

func (r *GormGeneratedDocumentRepository) GetByID(id string) (schema.GeneratedSponsorLetterAIDocument, error) {
	var doc schema.GeneratedSponsorLetterAIDocument
	if err := r.baseQuery().Where("id = ?", id).First(&doc).Error; err != nil {
		return schema.GeneratedSponsorLetterAIDocument{}, err
	}
	return doc, nil
}

func (r *GormGeneratedDocumentRepository) GetByStudentID(studentID string) (schema.GeneratedSponsorLetterAIDocument, error) {
	var doc schema.GeneratedSponsorLetterAIDocument
	if err := r.baseQuery().Where("student_id = ?", studentID).First(&doc).Error; err != nil {
		return schema.GeneratedSponsorLetterAIDocument{}, err
	}
	return doc, nil
}

func (r *GormGeneratedDocumentRepository) ListByStudentID(studentID string) ([]schema.GeneratedSponsorLetterAIDocument, error) {
	var docs []schema.GeneratedSponsorLetterAIDocument
	query := r.baseQuery().Model(&schema.GeneratedSponsorLetterAIDocument{}).Order("updated_at desc")
	if studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	if err := query.Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *GormGeneratedDocumentRepository) Create(doc *schema.GeneratedSponsorLetterAIDocument) error {
	return r.db.Create(doc).Error
}

func (r *GormGeneratedDocumentRepository) Update(doc *schema.GeneratedSponsorLetterAIDocument) error {
	return r.db.Save(doc).Error
}

func (r *GormGeneratedDocumentRepository) GetCurrentApprovalByDocumentID(documentID string) (schema.SponsorLetterAIApproval, error) {
	var item schema.SponsorLetterAIApproval
	if err := r.db.
		Preload("Reviewer").
		Preload("Logs").
		Preload("Logs.Actor").
		Where("document_id = ?", documentID).
		Order("updated_at desc").
		First(&item).Error; err != nil {
		return schema.SponsorLetterAIApproval{}, err
	}
	return item, nil
}

func (r *GormGeneratedDocumentRepository) CreateApproval(item *schema.SponsorLetterAIApproval) error {
	return r.db.Create(item).Error
}

func (r *GormGeneratedDocumentRepository) UpdateApproval(item *schema.SponsorLetterAIApproval) error {
	return r.db.Save(item).Error
}

func (r *GormGeneratedDocumentRepository) CreateApprovalLog(item *schema.SponsorLetterAIApprovalLog) error {
	return r.db.Create(item).Error
}

func (r *GormGeneratedDocumentRepository) GetDirector() (schema.User, error) {
	var user schema.User
	if err := r.db.Where("role = ?", schema.UserRoleDirector).Order("created_at asc").First(&user).Error; err != nil {
		return schema.User{}, err
	}
	return user, nil
}
