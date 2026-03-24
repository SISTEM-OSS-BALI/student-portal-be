package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/modules/user"
	"github.com/username/gin-gorm-api/internal/schema"
)

type Service struct {
	repo user.Repository
}

func NewService(repo user.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Authenticate(email, password string) (schema.User, string, error) {
	u, err := s.repo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schema.User{}, "", ErrInvalidCredentials
		}
		return schema.User{}, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return schema.User{}, "", ErrInvalidCredentials
	}

	token, err := GenerateToken(u)
	if err != nil {
		return schema.User{}, "", err
	}

	return u, token, nil
}

func (s *Service) Register(name, email, password string) (schema.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return schema.User{}, err
	}

	u := schema.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
		Role:     schema.UserRoleStudent,
	}

	if err := s.repo.Create(&u); err != nil {
		return schema.User{}, err
	}

	return u, nil
}

func (s *Service) GetByID(id string) (schema.User, error) {
	return s.repo.GetByID(id)
}
