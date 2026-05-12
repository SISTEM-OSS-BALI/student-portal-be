package promo

import (
	"time"

	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/schema"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(item *schema.Promo) error {
	return r.db.Create(item).Error
}

func (r *Repository) List(filter Filter) ([]schema.Promo, error) {
	q := r.db.Model(&schema.Promo{})
	if filter.ActiveOnly {
		now := time.Now()
		q = q.Where("is_active = ? AND valid_from <= ? AND valid_to >= ?", true, now, now)
	}

	var items []schema.Promo
	if err := q.Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) GetByID(id string) (schema.Promo, error) {
	var item schema.Promo
	err := r.db.First(&item, "id = ?", id).Error
	return item, err
}

func (r *Repository) Update(id string, updates map[string]interface{}) (schema.Promo, error) {
	if err := r.db.Model(&schema.Promo{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return schema.Promo{}, err
	}
	return r.GetByID(id)
}

func (r *Repository) Delete(id string) error {
	return r.db.Delete(&schema.Promo{}, "id = ?", id).Error
}
