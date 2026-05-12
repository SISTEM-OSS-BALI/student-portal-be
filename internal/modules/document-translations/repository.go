package documenttranslations

import (
	"errors"

	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

var ErrTranslationQuotaExceeded = errors.New("translation quota exceeded")

type Filter struct {
	StudentID        *string
	DocumentID       *string
	UploaderID       *string
	AnswerDocumentID *string
	Status           *string
}

type Repository interface {
	Create(item *schema.DocumentTranslation) error
	List(filter Filter) ([]schema.DocumentTranslation, error)
	GetByID(id string) (schema.DocumentTranslation, error)
	Update(item *schema.DocumentTranslation) error
	Delete(id string) error
	UpdateUserTranslationQuota(studentID string, pageCount int) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(item *schema.DocumentTranslation) error {
	var existing schema.DocumentTranslation
	err := r.db.
		Where("student_id = ? AND document_id = ?", item.StudentID, item.DocumentID).
		First(&existing).Error
	if err == nil {
		item.ID = existing.ID
		return r.Update(item)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if item.StudentID != "" && item.PageCount > 0 && !item.IsExistingTranslation {
			result := tx.Model(&schema.User{}).
				Where("id = ? AND translation_quota >= ?", item.StudentID, item.PageCount).
				Update("translation_quota", gorm.Expr("translation_quota - ?", item.PageCount))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return ErrTranslationQuotaExceeded
			}
		}

		if err := tx.Create(item).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *GormRepository) List(filter Filter) ([]schema.DocumentTranslation, error) {
	var items []schema.DocumentTranslation
	query := r.db.Model(&schema.DocumentTranslation{})
	if filter.StudentID != nil {
		query = query.Where("student_id = ?", *filter.StudentID)
	}
	if filter.DocumentID != nil {
		query = query.Where("document_id = ?", *filter.DocumentID)
	}
	if filter.UploaderID != nil {
		query = query.Where("uploader_id = ?", *filter.UploaderID)
	}
	if filter.AnswerDocumentID != nil {
		query = query.Where("answer_document_id = ?", *filter.AnswerDocumentID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if err := query.Order("updated_at desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GormRepository) GetByID(id string) (schema.DocumentTranslation, error) {
	var item schema.DocumentTranslation
	if err := r.db.Where("id = ?", id).First(&item).Error; err != nil {
		return schema.DocumentTranslation{}, err
	}
	return item, nil
}

func (r *GormRepository) Update(item *schema.DocumentTranslation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing schema.DocumentTranslation
		if err := tx.Where("id = ?", item.ID).First(&existing).Error; err != nil {
			return err
		}

		oldStudentID := existing.StudentID
		newStudentID := item.StudentID
		oldPages := existing.PageCount
		newPages := item.PageCount
		oldIsExistingTranslation := existing.IsExistingTranslation
		newIsExistingTranslation := item.IsExistingTranslation

		if oldStudentID == newStudentID {
			oldChargedPages := 0
			if !oldIsExistingTranslation {
				oldChargedPages = oldPages
			}
			newChargedPages := 0
			if !newIsExistingTranslation {
				newChargedPages = newPages
			}

			delta := newChargedPages - oldChargedPages
			if delta > 0 && newStudentID != "" {
				result := tx.Model(&schema.User{}).
					Where("id = ? AND translation_quota >= ?", newStudentID, delta).
					Update("translation_quota", gorm.Expr("translation_quota - ?", delta))
				if result.Error != nil {
					return result.Error
				}
				if result.RowsAffected == 0 {
					return ErrTranslationQuotaExceeded
				}
			}
			if delta < 0 && newStudentID != "" {
				if err := tx.Model(&schema.User{}).
					Where("id = ?", newStudentID).
					Update("translation_quota", gorm.Expr("translation_quota + ?", -delta)).
					Error; err != nil {
					return err
				}
			}
		} else {
			if oldStudentID != "" && oldPages > 0 && !oldIsExistingTranslation {
				if err := tx.Model(&schema.User{}).
					Where("id = ?", oldStudentID).
					Update("translation_quota", gorm.Expr("translation_quota + ?", oldPages)).
					Error; err != nil {
					return err
				}
			}
			if newStudentID != "" && newPages > 0 && !newIsExistingTranslation {
				result := tx.Model(&schema.User{}).
					Where("id = ? AND translation_quota >= ?", newStudentID, newPages).
					Update("translation_quota", gorm.Expr("translation_quota - ?", newPages))
				if result.Error != nil {
					return result.Error
				}
				if result.RowsAffected == 0 {
					return ErrTranslationQuotaExceeded
				}
			}
		}

		return tx.Save(item).Error
	})
}

func (r *GormRepository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var item schema.DocumentTranslation
		if err := tx.Where("id = ?", id).First(&item).Error; err != nil {
			return err
		}

		if err := tx.Delete(&schema.DocumentTranslation{}, "id = ?", id).Error; err != nil {
			return err
		}

		if item.StudentID == "" || item.PageCount <= 0 || item.IsExistingTranslation {
			return nil
		}
		return tx.Model(&schema.User{}).
			Where("id = ?", item.StudentID).
			Update("translation_quota", gorm.Expr("translation_quota + ?", item.PageCount)).
			Error
	})
}

func (r *GormRepository) UpdateUserTranslationQuota(studentID string, pageCount int) error {
	if studentID == "" || pageCount <= 0 {
		return nil
	}
	result := r.db.Model(&schema.User{}).
		Where("id = ? AND translation_quota >= ?", studentID, pageCount).
		Update("translation_quota", gorm.Expr("translation_quota - ?", pageCount))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTranslationQuotaExceeded
	}
	return nil
}
