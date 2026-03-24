package notesstudent

import (
	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type Repository interface {
	Create(note *schema.NoteStudent) error
	List() ([]schema.NoteStudent, error)
	ListByUserID(userID string) ([]schema.NoteStudent, error)
	GetByID(id string) (schema.NoteStudent, error)
	Update(note *schema.NoteStudent) error
	Delete(id string) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(note *schema.NoteStudent) error {
	return r.db.Create(note).Error
}

func (r *GormRepository) List() ([]schema.NoteStudent, error) {
	var notes []schema.NoteStudent
	if err := r.db.Order("id desc").Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *GormRepository) ListByUserID(userID string) ([]schema.NoteStudent, error) {
	var notes []schema.NoteStudent
	if err := r.db.Where("user_id = ?", userID).Order("id desc").Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *GormRepository) GetByID(id string) (schema.NoteStudent, error) {
	var note schema.NoteStudent
	if err := r.db.Where("id = ?", id).First(&note).Error; err != nil {
		return schema.NoteStudent{}, err
	}
	return note, nil
}

func (r *GormRepository) Update(note *schema.NoteStudent) error {
	return r.db.Save(note).Error
}

func (r *GormRepository) Delete(id string) error {
	return r.db.Delete(&schema.NoteStudent{}, "id = ?", id).Error
}
