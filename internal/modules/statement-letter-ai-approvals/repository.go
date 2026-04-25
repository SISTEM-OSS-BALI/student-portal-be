package statementletteraiapprovals

import (
	"fmt"

	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type Filter struct {
	DocumentID *string
	StudentID  *string
	ReviewerID *string
	Status     *schema.StatementLetterApprovalStatus
}

type Repository interface {
	GetDocumentByID(id string) (schema.GeneratedStatementLetterAIDocument, error)
	UpdateDocument(doc *schema.GeneratedStatementLetterAIDocument) error
	GetByID(id string) (schema.StatementLetterAIApproval, error)
	GetCurrentByDocumentID(documentID string) (schema.StatementLetterAIApproval, error)
	CreateApproval(item *schema.StatementLetterAIApproval) error
	UpdateApproval(item *schema.StatementLetterAIApproval) error
	List(filter Filter) ([]schema.StatementLetterAIApproval, error)
	CreateLog(item *schema.StatementLetterAIApprovalLog) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) tableName(model interface{}) (string, error) {
	stmt := &gorm.Statement{DB: r.db}
	if err := stmt.Parse(model); err != nil {
		return "", err
	}
	return stmt.Schema.Table, nil
}

func (r *GormRepository) approvalQuery() *gorm.DB {
	return r.db.
		Preload("Reviewer").
		Preload("Document").
		Preload("Document.Student").
		Preload("Logs").
		Preload("Logs.Actor")
}

func (r *GormRepository) documentQuery() *gorm.DB {
	return r.db.
		Preload("Student").
		Preload("CurrentApproval").
		Preload("CurrentApproval.Reviewer").
		Preload("Approvals").
		Preload("Approvals.Reviewer")
}

func (r *GormRepository) GetDocumentByID(id string) (schema.GeneratedStatementLetterAIDocument, error) {
	var doc schema.GeneratedStatementLetterAIDocument
	if err := r.documentQuery().Where("id = ?", id).First(&doc).Error; err != nil {
		return schema.GeneratedStatementLetterAIDocument{}, err
	}
	return doc, nil
}

func (r *GormRepository) UpdateDocument(doc *schema.GeneratedStatementLetterAIDocument) error {
	return r.db.Save(doc).Error
}

func (r *GormRepository) GetByID(id string) (schema.StatementLetterAIApproval, error) {
	var item schema.StatementLetterAIApproval
	if err := r.approvalQuery().Where("id = ?", id).First(&item).Error; err != nil {
		return schema.StatementLetterAIApproval{}, err
	}
	return item, nil
}

func (r *GormRepository) GetCurrentByDocumentID(documentID string) (schema.StatementLetterAIApproval, error) {
	var item schema.StatementLetterAIApproval
	if err := r.approvalQuery().Where("document_id = ?", documentID).Order("updated_at desc").First(&item).Error; err != nil {
		return schema.StatementLetterAIApproval{}, err
	}
	return item, nil
}

func (r *GormRepository) CreateApproval(item *schema.StatementLetterAIApproval) error {
	return r.db.Create(item).Error
}

func (r *GormRepository) UpdateApproval(item *schema.StatementLetterAIApproval) error {
	return r.db.Save(item).Error
}

func (r *GormRepository) List(filter Filter) ([]schema.StatementLetterAIApproval, error) {
	var items []schema.StatementLetterAIApproval
	query := r.approvalQuery().Model(&schema.StatementLetterAIApproval{}).Order("updated_at desc")
	approvalTable, err := r.tableName(&schema.StatementLetterAIApproval{})
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
		documentTable, err := r.tableName(&schema.GeneratedStatementLetterAIDocument{})
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

func (r *GormRepository) CreateLog(item *schema.StatementLetterAIApprovalLog) error {
	return r.db.Create(item).Error
}
