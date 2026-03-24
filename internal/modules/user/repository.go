package user

import (
	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *schema.User) error
	List() ([]schema.User, error)
	GetByID(id string) (schema.User, error)
	GetByEmail(email string) (schema.User, error)
	ListStudents() ([]schema.User, error)
	Update(user *schema.User) error
	Delete(id string) error
	PatchQuotaTranslation(id string, quota int) (schema.User, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(user *schema.User) error {
	return r.db.Create(user).Error
}

func (r *GormRepository) List() ([]schema.User, error) {
	var users []schema.User
	if err := r.db.Preload("Stage").Preload("Stage.Country").Preload("Stage.Document").Preload("NotesStudent").
		Preload("Stage.Country.CountrySteps.Step").
		Preload("Stage.Country.CountrySteps.Step.Children").
		Order("id desc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *GormRepository) GetByID(id string) (schema.User, error) {
	var user schema.User
	if err := r.db.Preload("Stage").Preload("Stage.Country").Preload("Stage.Document").Preload("NotesStudent").
		Preload("Stage.Country.CountrySteps.Step").
		Preload("Stage.Country.CountrySteps.Step.Children").
		Where("id = ?", id).First(&user).Error; err != nil {
		return schema.User{}, err
	}
	return user, nil
}

func (r *GormRepository) GetByEmail(email string) (schema.User, error) {
	var user schema.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return schema.User{}, err
	}
	return user, nil
}

func (r *GormRepository) Update(user *schema.User) error {
	return r.db.Save(user).Error
}

func (r *GormRepository) Delete(id string) error {
	return r.db.Delete(&schema.User{}, "id = ?", id).Error
}

func (r *GormRepository) ListStudents() ([]schema.User, error) {
	var users []schema.User
	if err := r.db.Preload("Stage").Preload("Stage.Country").Preload("Stage.Document").Preload("NotesStudent").
		Preload("Stage.Country.CountrySteps.Step").
		Preload("Stage.Country.CountrySteps.Step.Children").
		Where("role = ?", schema.UserRoleStudent).Order("id desc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *GormRepository) PatchQuotaTranslation(id string, quota int) (schema.User, error) {
	if err := r.db.Model(&schema.User{}).
		Where("id = ?", id).
		Update("translation_quota", gorm.Expr("translation_quota + ?", quota)).
		Error; err != nil {
		return schema.User{}, err
	}
	return r.GetByID(id)
}
 