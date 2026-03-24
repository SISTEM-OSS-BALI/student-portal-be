package steps

import (
	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type Repository interface {
	Create(step *schema.StepsManagement) error
	List() ([]schema.StepsManagement, error)
	GetByID(id string) (schema.StepsManagement, error)
	Update(step *schema.StepsManagement) error
	Delete(id string) error
	ReplaceChildren(stepID string, childIDs []string) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(step *schema.StepsManagement) error {
	return r.db.Create(step).Error
}

func (r *GormRepository) List() ([]schema.StepsManagement, error) {
	var steps []schema.StepsManagement
	if err := r.db.Preload("Children").Order("id desc").Find(&steps).Error; err != nil {
		return nil, err
	}
	return steps, nil
}

func (r *GormRepository) GetByID(id string) (schema.StepsManagement, error) {
	var step schema.StepsManagement
	if err := r.db.Preload("Children").Where("id = ?", id).First(&step).Error; err != nil {
		return schema.StepsManagement{}, err
	}
	return step, nil
}

func (r *GormRepository) Update(step *schema.StepsManagement) error {
	return r.db.Save(step).Error
}

func (r *GormRepository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		step := schema.StepsManagement{ID: id}
		if err := tx.Model(&step).Association("Children").Clear(); err != nil {
			return err
		}
		return tx.Delete(&schema.StepsManagement{}, "id = ?", id).Error
	})
}

func (r *GormRepository) ReplaceChildren(stepID string, childIDs []string) error {
	step := schema.StepsManagement{ID: stepID}
	children := make([]schema.ChildStepsManagement, 0, len(childIDs))
	for _, id := range childIDs {
		if id == "" {
			continue
		}
		children = append(children, schema.ChildStepsManagement{ID: id})
	}
	return r.db.Model(&step).Association("Children").Replace(children)
}
