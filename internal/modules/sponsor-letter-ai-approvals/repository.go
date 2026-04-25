package sponsorletteraiapprovals

import (
	"fmt"

	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type Filter struct {
	DocumentID *string
	StudentID  *string
	ReviewerID *string
	Status     *schema.SponsorLetterApprovalStatus
}

type Repository interface {
	GetDocumentByID(id string) (schema.GeneratedSponsorLetterAIDocument, error)
	UpdateDocument(doc *schema.GeneratedSponsorLetterAIDocument) error
	GetByID(id string) (schema.SponsorLetterAIApproval, error)
	GetCurrentByDocumentID(documentID string) (schema.SponsorLetterAIApproval, error)
	CreateApproval(item *schema.SponsorLetterAIApproval) error
	UpdateApproval(item *schema.SponsorLetterAIApproval) error
	List(filter Filter) ([]schema.SponsorLetterAIApproval, error)
	CreateLog(item *schema.SponsorLetterAIApprovalLog) error
}

type GormRepository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *GormRepository { return &GormRepository{db: db} }
func (r *GormRepository) tableName(model interface{}) (string, error) {
	stmt := &gorm.Statement{DB: r.db}
	if err := stmt.Parse(model); err != nil {
		return "", err
	}
	return stmt.Schema.Table, nil
}
func (r *GormRepository) approvalQuery() *gorm.DB {
	return r.db.Preload("Reviewer").Preload("Document").Preload("Document.Student").Preload("Logs").Preload("Logs.Actor")
}
func (r *GormRepository) documentQuery() *gorm.DB {
	return r.db.Preload("Student").Preload("CurrentApproval").Preload("CurrentApproval.Reviewer").Preload("Approvals").Preload("Approvals.Reviewer")
}
func (r *GormRepository) GetDocumentByID(id string) (schema.GeneratedSponsorLetterAIDocument, error) {
	var doc schema.GeneratedSponsorLetterAIDocument
	if err := r.documentQuery().Where("id = ?", id).First(&doc).Error; err != nil {
		return schema.GeneratedSponsorLetterAIDocument{}, err
	}
	return doc, nil
}
func (r *GormRepository) UpdateDocument(doc *schema.GeneratedSponsorLetterAIDocument) error {
	return r.db.Save(doc).Error
}
func (r *GormRepository) GetByID(id string) (schema.SponsorLetterAIApproval, error) {
	var item schema.SponsorLetterAIApproval
	if err := r.approvalQuery().Where("id = ?", id).First(&item).Error; err != nil {
		return schema.SponsorLetterAIApproval{}, err
	}
	return item, nil
}
func (r *GormRepository) GetCurrentByDocumentID(documentID string) (schema.SponsorLetterAIApproval, error) {
	var item schema.SponsorLetterAIApproval
	if err := r.approvalQuery().Where("document_id = ?", documentID).Order("updated_at desc").First(&item).Error; err != nil {
		return schema.SponsorLetterAIApproval{}, err
	}
	return item, nil
}
func (r *GormRepository) CreateApproval(item *schema.SponsorLetterAIApproval) error {
	return r.db.Create(item).Error
}
func (r *GormRepository) UpdateApproval(item *schema.SponsorLetterAIApproval) error {
	return r.db.Save(item).Error
}
func (r *GormRepository) List(filter Filter) ([]schema.SponsorLetterAIApproval, error) {
	var items []schema.SponsorLetterAIApproval
	query := r.approvalQuery().Model(&schema.SponsorLetterAIApproval{}).Order("updated_at desc")
	approvalTable, err := r.tableName(&schema.SponsorLetterAIApproval{})
	if err != nil {
		return nil, err
	}
	if filter.DocumentID != nil {
		query = query.Where(fmt.Sprintf("%s.document_id = ?", approvalTable), *filter.DocumentID)
	}
	if filter.ReviewerID != nil {
		query = query.Where(fmt.Sprintf("%s.reviewer_id = ?", approvalTable), *filter.ReviewerID)
	}
	if filter.Status != nil {
		query = query.Where(fmt.Sprintf("%s.status = ?", approvalTable), *filter.Status)
	}
	if filter.StudentID != nil {
		documentTable, err := r.tableName(&schema.GeneratedSponsorLetterAIDocument{})
		if err != nil {
			return nil, err
		}
		query = query.Joins(fmt.Sprintf("JOIN %s d ON d.id = %s.document_id", documentTable, approvalTable)).Where("d.student_id = ?", *filter.StudentID)
	}
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
func (r *GormRepository) CreateLog(item *schema.SponsorLetterAIApprovalLog) error {
	return r.db.Create(item).Error
}
