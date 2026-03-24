package auth

import "github.com/username/gin-gorm-api/internal/schema"

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
	Token string `json:"token"`
	User  UserDTO `json:"user"`
}

type RegisterResponseDTO struct {
	User UserDTO `json:"user"`
}

func NewUserDTO(user schema.User) UserDTO {
	return UserDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}
