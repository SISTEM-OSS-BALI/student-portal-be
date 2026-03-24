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
