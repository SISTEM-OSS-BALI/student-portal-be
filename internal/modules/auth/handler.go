package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/username/gin-gorm-api/internal/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(c *gin.Context) {
	var input LoginDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	user, token, err := h.service.Authenticate(input.Email, input.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			httpx.RespondError(c, http.StatusUnauthorized, "invalid_credentials", "email atau password salah", nil)
			return
		}
		httpx.RespondError(c, http.StatusInternalServerError, "login_failed", err.Error(), nil)
		return
	}

	setAuthCookie(c, token)

	c.JSON(http.StatusOK, LoginResponseDTO{
		Token: token,
		User:  NewUserDTO(user),
	})
}

func (h *Handler) Register(c *gin.Context) {
	var input RegisterDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	user, err := h.service.Register(input.Name, input.Email, input.Password)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "register_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, RegisterResponseDTO{
		User: NewUserDTO(user),
	})
}

func (h *Handler) Me(c *gin.Context) {
	claimsAny, ok := c.Get("auth")
	if !ok {
		httpx.RespondError(c, http.StatusUnauthorized, "unauthorized", "missing auth claims", nil)
		return
	}

	claims, ok := claimsAny.(*Claims)
	if !ok || claims == nil {
		httpx.RespondError(c, http.StatusUnauthorized, "unauthorized", "invalid auth claims", nil)
		return
	}

	user, err := h.service.GetByID(claims.UserID)
	if err != nil {
		httpx.RespondError(c, http.StatusNotFound, "not_found", "user not found", nil)
		return
	}

	c.JSON(http.StatusOK, NewUserDTO(user))
}

func (h *Handler) ForgotPassword(c *gin.Context) {
	var input ForgotPasswordDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	otp, expiresAt, err := h.service.ForgotPassword(input.Email)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "forgot_password_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, ForgotPasswordResponseDTO{
		Message:      "Jika email terdaftar sebagai student, instruksi reset password telah diproses.",
		ResetOTP:     otp,
		ResetExpires: expiresAt,
	})
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var input ResetPasswordDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	if err := h.service.ResetPassword(input.Email, input.OTP, input.NewPassword, input.ConfirmPassword); err != nil {
		if errors.Is(err, ErrInvalidResetToken) {
			httpx.RespondError(c, http.StatusBadRequest, "invalid_reset_token", "kode OTP tidak valid atau sudah kedaluwarsa", nil)
			return
		}
		httpx.RespondError(c, http.StatusBadRequest, "reset_password_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diperbarui."})
}
