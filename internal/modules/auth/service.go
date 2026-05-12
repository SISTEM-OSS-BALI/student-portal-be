package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/username/gin-gorm-api/internal/notify"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/modules/user"
	"github.com/username/gin-gorm-api/internal/schema"
)

type Service struct {
	repo     user.Repository
	notifier *notify.Service
	db       *gorm.DB
}

func NewService(repo user.Repository, notifier *notify.Service, db *gorm.DB) *Service {
	return &Service{repo: repo, notifier: notifier, db: db}
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

func (s *Service) ForgotPassword(email string) (*string, *time.Time, error) {
	u, err := s.repo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	if u.Role != schema.UserRoleStudent {
		return nil, nil, nil
	}

	otp, err := generateOTP()
	if err != nil {
		return nil, nil, err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	expiresAt := time.Now().Add(15 * time.Minute)
	resetRecord := schema.PasswordResetOTP{
		UserID:    u.ID,
		Email:     u.Email,
		Code:      string(hashed),
		ExpiresAt: expiresAt,
	}
	if err := s.db.Create(&resetRecord).Error; err != nil {
		return nil, nil, err
	}

	if s.notifier != nil && s.notifier.Enabled() {
		_ = s.notifier.SendForgotPasswordEmail(u.Email, u.Name, otp, expiresAt)
	}

	// Return OTP for non-email test fallback if notifier disabled.
	if s.notifier == nil || !s.notifier.Enabled() {
		return &otp, &expiresAt, nil
	}
	return nil, &expiresAt, nil
}

func (s *Service) ResetPassword(email, otp, newPassword, confirmPassword string) error {
	if newPassword != confirmPassword {
		return errors.New("new_password and confirm_password do not match")
	}

	u, err := s.repo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidResetToken
		}
		return err
	}
	if u.Role != schema.UserRoleStudent {
		return ErrInvalidResetToken
	}

	var resetRecord schema.PasswordResetOTP
	if err := s.db.
		Where("user_id = ? AND email = ? AND used_at IS NULL", u.ID, u.Email).
		Order("created_at desc").
		First(&resetRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidResetToken
		}
		return err
	}
	if time.Now().After(resetRecord.ExpiresAt) {
		return ErrInvalidResetToken
	}
	if err := bcrypt.CompareHashAndPassword([]byte(resetRecord.Code), []byte(otp)); err != nil {
		return ErrInvalidResetToken
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashed)
	if err := s.repo.Update(&u); err != nil {
		return err
	}

	now := time.Now()
	return s.db.Model(&schema.PasswordResetOTP{}).
		Where("id = ?", resetRecord.ID).
		Update("used_at", &now).Error
}

func generateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
