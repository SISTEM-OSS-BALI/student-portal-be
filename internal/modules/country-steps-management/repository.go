package countrysteps

import (
	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type Repository interface {
	Create(item *schema.CountryStepsManagement) error
	List() ([]schema.CountryStepsManagement, error)
	ListWithCountryAndStep() ([]schema.CountryStepsManagement, error)
	GetByID(id string) (schema.CountryStepsManagement, error)
	Update(item *schema.CountryStepsManagement) error
	Delete(id string) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(item *schema.CountryStepsManagement) error {
	return r.db.Create(item).Error
}

func (r *GormRepository) List() ([]schema.CountryStepsManagement, error) {
	var items []schema.CountryStepsManagement
	if err := r.db.Preload("Country").
		Preload("Step").
		Preload("Step.Children").
		Order("id desc").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GormRepository) ListWithCountryAndStep() ([]schema.CountryStepsManagement, error) {
	return r.List()
}

func (r *GormRepository) GetByID(id string) (schema.CountryStepsManagement, error) {
	var item schema.CountryStepsManagement
	if err := r.db.Preload("Country").
		Preload("Step").
		Preload("Step.Children").
		Where("id = ?", id).
		First(&item).Error; err != nil {
		return schema.CountryStepsManagement{}, err
	}
	return item, nil
}

func (r *GormRepository) Update(item *schema.CountryStepsManagement) error {
	return r.db.Save(item).Error
}

func (r *GormRepository) Delete(id string) error {
	return r.db.Delete(&schema.CountryStepsManagement{}, "id = ?", id).Error
}
