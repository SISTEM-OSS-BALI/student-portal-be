package informationcountrymanagement

import (
	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type Repository interface {
	Create(data *schema.InformationCountryManagement) error
	List() ([]schema.InformationCountryManagement, error)
	GetByID(id string) (schema.InformationCountryManagement, error)
	GetBySlug(slug string) (schema.InformationCountryManagement, error)
	ListByCountryID(countryID string) ([]schema.InformationCountryManagement, error)
	Update(data *schema.InformationCountryManagement) error
	Delete(id string) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(data *schema.InformationCountryManagement) error {
	return r.db.Create(data).Error
}

func (r *GormRepository) List() ([]schema.InformationCountryManagement, error) {
	var result []schema.InformationCountryManagement

	if err := r.db.
		Preload("Country").
		Order("CASE priority WHEN 'high' THEN 1 WHEN 'medium' THEN 2 WHEN 'normal' THEN 3 WHEN 'low' THEN 4 ELSE 5 END").
		Order("created_at desc").
		Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r *GormRepository) GetByID(id string) (schema.InformationCountryManagement, error) {
	var result schema.InformationCountryManagement

	if err := r.db.
		Preload("Country").
		Where("id = ?", id).
		First(&result).Error; err != nil {
		return schema.InformationCountryManagement{}, err
	}

	return result, nil
}

func (r *GormRepository) GetBySlug(slug string) (schema.InformationCountryManagement, error) {
	var result schema.InformationCountryManagement

	if err := r.db.
		Preload("Country").
		Where("slug = ?", slug).
		First(&result).Error; err != nil {
		return schema.InformationCountryManagement{}, err
	}

	return result, nil
}

func (r *GormRepository) ListByCountryID(countryID string) ([]schema.InformationCountryManagement, error) {
	var result []schema.InformationCountryManagement

	if err := r.db.
		Preload("Country").
		Where("country_id = ?", countryID).
		Order("CASE priority WHEN 'high' THEN 1 WHEN 'medium' THEN 2 WHEN 'normal' THEN 3 WHEN 'low' THEN 4 ELSE 5 END").
		Order("created_at desc").
		Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r *GormRepository) Update(data *schema.InformationCountryManagement) error {
	return r.db.Save(data).Error
}

func (r *GormRepository) Delete(id string) error {
	return r.db.Delete(&schema.InformationCountryManagement{}, "id = ?", id).Error
}
