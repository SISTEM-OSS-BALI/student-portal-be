package stages

import (
	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type Repository interface {
	Create(doc *schema.StageManagement) error
	List() ([]schema.StageManagement, error)
	ListWithCountryAndDocument() ([]schema.StageManagement, error)
	GetByID(id string) (schema.StageManagement, error)
	Update(doc *schema.StageManagement) error
	Delete(id string) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(doc *schema.StageManagement) error {
	return r.db.Create(doc).Error
}

func (r *GormRepository) List() ([]schema.StageManagement, error) {
	var docs []schema.StageManagement
	if err := r.db.Preload("Document").Preload("Country").Order("id desc").Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *GormRepository) ListWithCountryAndDocument() ([]schema.StageManagement, error) {
	return r.List()
}

func (r *GormRepository) GetByID(id string) (schema.StageManagement, error) {
	var doc schema.StageManagement
	if err := r.db.Preload("Document").Preload("Country").Where("id = ?", id).First(&doc).Error; err != nil {
		return schema.StageManagement{}, err
	}
	return doc, nil
}

func (r *GormRepository) Update(doc *schema.StageManagement) error {
	return r.db.Save(doc).Error
}

func (r *GormRepository) Delete(id string) error {
	return r.db.Delete(&schema.StageManagement{}, "id = ?", id).Error
}
