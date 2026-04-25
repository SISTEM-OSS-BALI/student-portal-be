package sponsorletteraiapprovals

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/httpx"
	"github.com/username/gin-gorm-api/internal/modules/auth"
	"github.com/username/gin-gorm-api/internal/schema"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }
func getAuthClaims(c *gin.Context) (*auth.Claims, bool) {
	value, ok := c.Get("auth")
	if !ok {
		return nil, false
	}
	claims, ok := value.(*auth.Claims)
	return claims, ok
}
func (h *Handler) CreateOrUpdate(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		httpx.RespondError(c, http.StatusUnauthorized, "unauthorized", "missing auth claims", nil)
		return
	}
	var input CreateOrUpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}
	item, err := h.service.CreateOrUpdate(input, claims.UserID, claims.Role)
	if err != nil {
		statusCode := http.StatusBadRequest
		switch {
		case errors.Is(err, ErrDirectorRoleRequired), errors.Is(err, ErrApprovalAssignedToAnotherDirector):
			statusCode = http.StatusForbidden
		case errors.Is(err, gorm.ErrRecordNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, ErrSponsorLetterNotSubmitted):
			statusCode = http.StatusConflict
		}
		httpx.RespondError(c, statusCode, "sponsor_letter_ai_approval_upsert_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewResponseDTO(item))
}
func (h *Handler) List(c *gin.Context) {
	var filter Filter
	if value := c.Query("document_id"); value != "" {
		filter.DocumentID = &value
	}
	if value := c.Query("student_id"); value != "" {
		filter.StudentID = &value
	}
	if value := c.Query("reviewer_id"); value != "" {
		filter.ReviewerID = &value
	}
	if value := c.Query("status"); value != "" {
		status, err := normalizeApprovalStatus(value)
		if err != nil {
			httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
			return
		}
		filter.Status = &status
	}
	items, err := h.service.List(filter)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "sponsor_letter_ai_approval_list_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewResponseListDTO(items))
}
func (h *Handler) GetByID(c *gin.Context) {
	item, err := h.service.GetByID(c.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpx.RespondError(c, http.StatusNotFound, "sponsor_letter_ai_approval_not_found", "sponsor letter ai approval not found", nil)
			return
		}
		httpx.RespondError(c, http.StatusInternalServerError, "sponsor_letter_ai_approval_get_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewResponseDTO(item))
}
func (h *Handler) Update(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		httpx.RespondError(c, http.StatusUnauthorized, "unauthorized", "missing auth claims", nil)
		return
	}
	var input UpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}
	item, err := h.service.Update(c.Param("id"), input, claims.UserID, claims.Role)
	if err != nil {
		statusCode := http.StatusBadRequest
		switch {
		case errors.Is(err, ErrDirectorRoleRequired), errors.Is(err, ErrApprovalAssignedToAnotherDirector):
			statusCode = http.StatusForbidden
		case errors.Is(err, gorm.ErrRecordNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, ErrSponsorLetterNotSubmitted):
			statusCode = http.StatusConflict
		}
		httpx.RespondError(c, statusCode, "sponsor_letter_ai_approval_update_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewResponseDTO(item))
}

var _ schema.UserRole = schema.UserRoleDirector
