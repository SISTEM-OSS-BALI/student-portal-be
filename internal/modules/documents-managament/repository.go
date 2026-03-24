package documents

import (
	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type Repository interface {
	Create(doc *schema.DocumentsManagement) error
	List() ([]schema.DocumentsManagement, error)
	GetByID(id string) (schema.DocumentsManagement, error)
	Update(doc *schema.DocumentsManagement) error
	Delete(id string) error
	DocumentTranslationRequired() ([]schema.DocumentsManagement, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(doc *schema.DocumentsManagement) error {
	return r.db.Create(doc).Error
}

func (r *GormRepository) List() ([]schema.DocumentsManagement, error) {
	var docs []schema.DocumentsManagement
	if err := r.db.Order("id desc").Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *GormRepository) GetByID(id string) (schema.DocumentsManagement, error) {
	var doc schema.DocumentsManagement
	if err := r.db.Where("id = ?", id).First(&doc).Error; err != nil {
		return schema.DocumentsManagement{}, err
	}
	return doc, nil
}

func (r *GormRepository) Update(doc *schema.DocumentsManagement) error {
	return r.db.Save(doc).Error
}

func (r *GormRepository) Delete(id string) error {
	return r.db.Delete(&schema.DocumentsManagement{}, "id = ?", id).Error
}

func (r *GormRepository) DocumentTranslationRequired() ([]schema.DocumentsManagement, error) {
	var docs []schema.DocumentsManagement
	if err := r.db.
		Model(&schema.DocumentsManagement{}).
		Where("translation_needed = ?", schema.TranslationNeededYes).
		Order("id desc").
		Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}
