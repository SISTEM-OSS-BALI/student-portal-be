package user

import (
	"time"

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
	PatchVisaStatus(id string, visaStatus *string, visaGrantedAt *time.Time) (schema.User, error)
	PatchStudentStatus(id string, studentStatus schema.StatusStudent, updatedByID *string, updatedAt *time.Time) (schema.User, error)
	PatchDocumentConsent(id string, payload PatchDocumentConsentDTO) (schema.User, error)
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
	if err := r.db.Preload("Stage").Preload("Stage.Country").Preload("Stage.Document").Preload("NotesStudent").Preload("StudentStatusUpdatedBy").
		Preload("Stage.Country.CountrySteps.Step").
		Preload("Stage.Country.CountrySteps.Step.Children").
		Order("id desc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *GormRepository) GetByID(id string) (schema.User, error) {
	var user schema.User
	if err := r.db.Preload("Stage").Preload("Stage.Country").Preload("Stage.Document").Preload("NotesStudent").Preload("StudentStatusUpdatedBy").
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
	if err := r.db.Preload("Stage").Preload("Stage.Country").Preload("Stage.Document").Preload("NotesStudent").Preload("StudentStatusUpdatedBy").
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

func (r *GormRepository) PatchVisaStatus(id string, visaStatus *string, visaGrantedAt *time.Time) (schema.User, error) {
	if err := r.db.Model(&schema.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"visa_status":     visaStatus,
			"visa_granted_at": visaGrantedAt,
		}).
		Error; err != nil {
		return schema.User{}, err
	}
	return r.GetByID(id)
}

func (r *GormRepository) PatchStudentStatus(id string, studentStatus schema.StatusStudent, updatedByID *string, updatedAt *time.Time) (schema.User, error) {
	if err := r.db.Model(&schema.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"student_status":               studentStatus,
			"student_status_updated_by_id": updatedByID,
			"student_status_updated_at":    updatedAt,
		}).
		Error; err != nil {
		return schema.User{}, err
	}
	return r.GetByID(id)
}

func (r *GormRepository) PatchDocumentConsent(id string, payload PatchDocumentConsentDTO) (schema.User, error) {
	updates := map[string]interface{}{}

	if payload.DocumentConsentSignatureURL != nil {
		updates["document_consent_signature_url"] = *payload.DocumentConsentSignatureURL
	}

	if payload.DocumentConsentProofPhotoURL != nil {
		updates["document_consent_proof_photo_url"] = *payload.DocumentConsentProofPhotoURL
	}

	if payload.DocumentConsentSignedAt != nil {
		updates["document_consent_signed_at"] = *payload.DocumentConsentSignedAt
	}

	if payload.DocumentConsentSigned != nil {
		updates["document_consent_signed"] = *payload.DocumentConsentSigned
	}

	if len(updates) == 0 {
		return r.GetByID(id)
	}

	if err := r.db.Model(&schema.User{}).
		Where("id = ?", id).
		Updates(updates).
		Error; err != nil {
		return schema.User{}, err
	}

	return r.GetByID(id)
}
