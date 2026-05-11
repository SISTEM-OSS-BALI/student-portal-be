package user

import (
	"errors"
	"strings"
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func isVisaGrantedStatus(value *string) bool {
	if value == nil {
		return false
	}

	switch strings.ToUpper(strings.TrimSpace(*value)) {
	case "GRANT", "GRANTED":
		return true
	default:
		return false
	}
}

func resolveVisaGrantedAt(visaStatus *string, existing *time.Time) *time.Time {
	if !isVisaGrantedStatus(visaStatus) {
		return nil
	}
	if existing != nil {
		return existing
	}

	now := time.Now()
	return &now
}

func resolveStudentStatus(value *string, fallback schema.StatusStudent) schema.StatusStudent {
	value = normalizeOptionalString(value)
	if value == nil || *value == "" {
		return fallback
	}

	return schema.StatusStudent(*value)
}

func resolveAuditActorID(actorID *string) *string {
	actorID = normalizeOptionalString(actorID)
	if actorID == nil || *actorID == "" {
		return nil
	}
	return actorID
}

func buildStudentStatusAudit(actorID *string) (*string, *time.Time) {
	now := time.Now()
	return resolveAuditActorID(actorID), &now
}

func (s *Service) Create(
	name, email, password string,
	stageID, currentStepID, visaStatus, studentStatus, nameConsultant, noPhone, nameCampus, degree, nameDegree, visaType, source *string,
	translationQuota int,
) (schema.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return schema.User{}, err
	}

	visaStatus = normalizeOptionalString(visaStatus)
	visaType = normalizeOptionalString(visaType)
	if visaType != nil {
		ok, err := s.repo.VisaTypeExists(*visaType)
		if err != nil {
			return schema.User{}, err
		}
		if !ok {
			return schema.User{}, errors.New("invalid visa_type: not found")
		}
	}

	user := schema.User{
		Name:             name,
		Email:            email,
		Password:         string(hashed),
		Role:             schema.UserRoleStudent,
		StageID:          stageID,
		CurrentStepID:    currentStepID,
		VisaStatus:       visaStatus,
		VisaGrantedAt:    resolveVisaGrantedAt(visaStatus, nil),
		StudentStatus:    resolveStudentStatus(studentStatus, schema.StatusStudentOnGoing),
		NameConsultant:   nameConsultant,
		NoPhone:          noPhone,
		NameCampus:       nameCampus,
		Degree:           degree,
		NameDegree:       nameDegree,
		VisaType:         visaType,
		TranslationQuota: translationQuota,
		Source:           source,
	}
	if err := s.repo.Create(&user); err != nil {
		return schema.User{}, err
	}
	return user, nil
}

func (s *Service) List() ([]schema.User, error) {
	return s.repo.List()
}

func (s *Service) GetByID(id string) (schema.User, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(
	id string,
	name, email, stageID, currentStepID, visaStatus, studentStatus, nameConsultant, nameCampus, noPhone, degree, nameDegree, visaType, source *string,
	translationQuota *int,
	actorID *string,
) (schema.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return schema.User{}, err
	}
	if name != nil {
		user.Name = *name
	}
	if email != nil {
		user.Email = *email
	}
	if stageID != nil {
		user.StageID = stageID
	}
	if currentStepID != nil {
		user.CurrentStepID = currentStepID
	}
	if visaStatus != nil {
		visaStatus = normalizeOptionalString(visaStatus)
		user.VisaStatus = visaStatus
		user.VisaGrantedAt = resolveVisaGrantedAt(visaStatus, user.VisaGrantedAt)
	}
	if studentStatus != nil {
		nextStatus := resolveStudentStatus(studentStatus, user.StudentStatus)
		if nextStatus != user.StudentStatus {
			user.StudentStatus = nextStatus
			user.StudentStatusUpdatedByID, user.StudentStatusUpdatedAt = buildStudentStatusAudit(actorID)
		}
	}

	if nameConsultant != nil {
		user.NameConsultant = nameConsultant
	}
	if nameCampus != nil {
		user.NameCampus = nameCampus
	}
	if noPhone != nil {
		user.NoPhone = noPhone
	}
	if degree != nil {
		user.Degree = degree
	}
	if translationQuota != nil {
		user.TranslationQuota = *translationQuota
	}

	if nameDegree != nil {
		user.NameDegree = nameDegree
	}
	if visaType != nil {
		visaType = normalizeOptionalString(visaType)
		if visaType != nil {
			ok, err := s.repo.VisaTypeExists(*visaType)
			if err != nil {
				return schema.User{}, err
			}
			if !ok {
				return schema.User{}, errors.New("invalid visa_type: not found")
			}
		}
		user.VisaType = visaType
	}
	if source != nil {
		user.Source = normalizeOptionalString(source)
	}
	if err := s.repo.Update(&user); err != nil {
		return schema.User{}, err
	}
	return user, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *Service) ListStudents() ([]schema.User, error) {
	return s.repo.ListStudents()
}

func (s *Service) PatchQuotaTranslation(id string, quota int) (schema.User, error) {
	return s.repo.PatchQuotaTranslation(id, quota)
}

func (s *Service) PatchVisaStatus(id string, visaStatus *string) (schema.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return schema.User{}, err
	}

	visaStatus = normalizeOptionalString(visaStatus)
	visaGrantedAt := resolveVisaGrantedAt(visaStatus, user.VisaGrantedAt)
	return s.repo.PatchVisaStatus(id, visaStatus, visaGrantedAt)
}

func (s *Service) PatchStudentStatus(id string, studentStatus *string, actorID *string) (schema.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return schema.User{}, err
	}

	nextStatus := resolveStudentStatus(studentStatus, user.StudentStatus)
	if nextStatus == user.StudentStatus {
		return user, nil
	}

	updatedByID, updatedAt := buildStudentStatusAudit(actorID)
	return s.repo.PatchStudentStatus(id, nextStatus, updatedByID, updatedAt)
}

func (s *Service) PatchDocumentConsent(id string, payload PatchDocumentConsentDTO) (schema.User, error) {
	return s.repo.PatchDocumentConsent(id, payload)
}
