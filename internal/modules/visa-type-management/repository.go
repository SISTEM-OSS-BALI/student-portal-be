package visatype

import (
	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/schema"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(item *schema.VisaTypeManagement) error {
	return r.db.Create(item).Error
}

func (r *Repository) List(filter Filter) ([]schema.VisaTypeManagement, error) {
	q := r.db.Model(&schema.VisaTypeManagement{}).Preload("Country")
	if filter.CountryID != nil && *filter.CountryID != "" {
		q = q.Where("country_id = ?", *filter.CountryID)
	}
	var items []schema.VisaTypeManagement
	if err := q.Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) GetByID(id string) (schema.VisaTypeManagement, error) {
	var item schema.VisaTypeManagement
	err := r.db.Preload("Country").First(&item, "id = ?", id).Error
	return item, err
}

func (r *Repository) Update(id string, updates map[string]interface{}) (schema.VisaTypeManagement, error) {
	if err := r.db.Model(&schema.VisaTypeManagement{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return schema.VisaTypeManagement{}, err
	}
	return r.GetByID(id)
}

func (r *Repository) Delete(id string) error {
	return r.db.Delete(&schema.VisaTypeManagement{}, "id = ?", id).Error
}
