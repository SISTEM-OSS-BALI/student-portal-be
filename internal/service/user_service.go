package service

import (
	"github.com/username/gin-gorm-api/internal/models"
	"github.com/username/gin-gorm-api/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(name, email string) (models.User, error) {
	user := models.User{Name: name, Email: email}
	if err := s.repo.Create(&user); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *UserService) List() ([]models.User, error) {
	return s.repo.List()
}

func (s *UserService) GetByID(id string) (models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) Update(id string, name, email *string) (models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return models.User{}, err
	}
	if name != nil {
		user.Name = *name
	}
	if email != nil {
		user.Email = *email
	}
	if err := s.repo.Update(&user); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *UserService) Delete(id string) error {
	return s.repo.Delete(id)
}
