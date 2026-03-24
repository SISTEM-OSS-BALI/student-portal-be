package user

import (
	"github.com/username/gin-gorm-api/internal/schema"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(
	name, email, password string,
	stageID, noPhone, nameCampus, degree, nameDegree, visaType *string,
	translationQuota int,
) (schema.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return schema.User{}, err
	}

	user := schema.User{
		Name:             name,
		Email:            email,
		Password:         string(hashed),
		Role:             schema.UserRoleStudent,
		StageID:          stageID,
		NoPhone:          noPhone,
		NameCampus:       nameCampus,
		Degree:           degree,
		NameDegree:       nameDegree,
		VisaType:         visaType,
		TranslationQuota: translationQuota,
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
	name, email, stageID, nameCampus, noPhone, degree, nameDegree, visaType *string,
	translationQuota *int,
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
		user.VisaType = visaType
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
