package childsteps

import (
	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type Repository interface {
	Create(child *schema.ChildStepsManagement) error
	List() ([]schema.ChildStepsManagement, error)
	ListWithSteps() ([]schema.ChildStepsManagement, error)
	GetByID(id string) (schema.ChildStepsManagement, error)
	Update(child *schema.ChildStepsManagement) error
	Delete(id string) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(child *schema.ChildStepsManagement) error {
	return r.db.Create(child).Error
}

func (r *GormRepository) List() ([]schema.ChildStepsManagement, error) {
	var children []schema.ChildStepsManagement
	if err := r.db.Order("id desc").Find(&children).Error; err != nil {
		return nil, err
	}
	return children, nil
}

func (r *GormRepository) ListWithSteps() ([]schema.ChildStepsManagement, error) {
	var children []schema.ChildStepsManagement
	if err := r.db.Preload("Steps").Order("id desc").Find(&children).Error; err != nil {
		return nil, err
	}
	return children, nil
}

func (r *GormRepository) GetByID(id string) (schema.ChildStepsManagement, error) {
	var child schema.ChildStepsManagement
	if err := r.db.Preload("Steps").Where("id = ?", id).First(&child).Error; err != nil {
		return schema.ChildStepsManagement{}, err
	}
	return child, nil
}

func (r *GormRepository) Update(child *schema.ChildStepsManagement) error {
	return r.db.Save(child).Error
}

func (r *GormRepository) Delete(id string) error {
	return r.db.Delete(&schema.ChildStepsManagement{}, "id = ?", id).Error
}
