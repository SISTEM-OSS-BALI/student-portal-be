package country

import (
	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type Repository interface {
	Create(country *schema.CountryManagement) error
	List() ([]schema.CountryManagement, error)
	GetByID(id string) (schema.CountryManagement, error)
	Update(country *schema.CountryManagement) error
	Delete(id string) error
	CountDocumentsByCountryIDs(countryIDs []string) (map[string]int64, error)
	CountDocumentsByCountryID(countryID string) (int64, error)
	CountStepsByCountryIDs(countryIDs []string) (map[string]int64, error)
	CountStepsByCountryID(countryID string) (int64, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(country *schema.CountryManagement) error {
	return r.db.Create(country).Error
}

func (r *GormRepository) List() ([]schema.CountryManagement, error) {
	var countries []schema.CountryManagement
	if err := r.db.Order("id desc").Find(&countries).Error; err != nil {
		return nil, err
	}
	return countries, nil
}

func (r *GormRepository) GetByID(id string) (schema.CountryManagement, error) {
	var country schema.CountryManagement
	if err := r.db.Where("id = ?", id).First(&country).Error; err != nil {
		return schema.CountryManagement{}, err
	}
	return country, nil
}

func (r *GormRepository) Update(country *schema.CountryManagement) error {
	return r.db.Save(country).Error
}

func (r *GormRepository) Delete(id string) error {
	return r.db.Delete(&schema.CountryManagement{}, "id = ?", id).Error
}

func (r *GormRepository) CountDocumentsByCountryIDs(countryIDs []string) (map[string]int64, error) {
	out := make(map[string]int64)
	if len(countryIDs) == 0 {
		return out, nil
	}

	type row struct {
		CountryID string `gorm:"column:country_id"`
		Total     int64  `gorm:"column:total"`
	}
	var rows []row
	if err := r.db.Model(&schema.StageManagement{}).
		Select("country_id, COUNT(DISTINCT document_id) as total").
		Where("country_id IN ?", countryIDs).
		Group("country_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, r := range rows {
		out[r.CountryID] = r.Total
	}
	return out, nil
}

func (r *GormRepository) CountDocumentsByCountryID(countryID string) (int64, error) {
	var total int64
	if err := r.db.Model(&schema.StageManagement{}).
		Where("country_id = ?", countryID).
		Distinct("document_id").
		Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *GormRepository) CountStepsByCountryIDs(countryIDs []string) (map[string]int64, error) {
	out := make(map[string]int64)
	if len(countryIDs) == 0 {
		return out, nil
	}

	type row struct {
		CountryID string `gorm:"column:country_id"`
		Total     int64  `gorm:"column:total"`
	}
	var rows []row
	if err := r.db.Model(&schema.CountryStepsManagement{}).
		Select("country_id, COUNT(DISTINCT step_id) as total").
		Where("country_id IN ?", countryIDs).
		Group("country_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, r := range rows {
		out[r.CountryID] = r.Total
	}
	return out, nil
}

func (r *GormRepository) CountStepsByCountryID(countryID string) (int64, error) {
	var total int64
	if err := r.db.Model(&schema.CountryStepsManagement{}).
		Where("country_id = ?", countryID).
		Distinct("step_id").
		Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}
