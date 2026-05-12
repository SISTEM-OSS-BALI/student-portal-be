package auth

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserDTO struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Email string          `json:"email"`
	Role  schema.UserRole `json:"role"`
}

type LoginResponseDTO struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

type RegisterResponseDTO struct {
	User UserDTO `json:"user"`
}

type ForgotPasswordDTO struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordResponseDTO struct {
	Message      string     `json:"message"`
	ResetOTP     *string    `json:"reset_otp,omitempty"`
	ResetExpires *time.Time `json:"reset_expires,omitempty"`
}

type ResetPasswordDTO struct {
	Email           string `json:"email" binding:"required,email"`
	OTP             string `json:"otp" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

func NewUserDTO(user schema.User) UserDTO {
	return UserDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}
