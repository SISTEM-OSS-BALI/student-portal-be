package user

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"github.com/username/gin-gorm-api/internal/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func getAuthUserID(c *gin.Context) *string {
	value, ok := c.Get("auth")
	if !ok {
		return nil
	}

	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return nil
	}
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil
	}

	field := rv.FieldByName("UserID")
	if !field.IsValid() || field.Kind() != reflect.String {
		return nil
	}

	userID := field.String()
	if userID == "" {
		return nil
	}
	return &userID
}

func (h *Handler) Create(c *gin.Context) {
	var input CreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	user, err := h.service.Create(
		input.Name,
		input.Email,
		input.Password,
		input.StageID,
		input.CurrentStepID,
		input.VisaStatus,
		input.StudentStatus,
		input.NameConsultant,
		input.NoPhone,
		input.NameCampus,
		input.Degree,
		input.NameDegree,
		input.VisaType,
		input.Source,
		input.TranslationQuota,
	)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "create_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, NewResponseDTO(user))
}

func (h *Handler) List(c *gin.Context) {
	users, err := h.service.List()
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewResponseListDTO(users))
}

func (h *Handler) GetByID(c *gin.Context) {
	user, err := h.service.GetByID(c.Param("id"))
	if err != nil {
		httpx.RespondError(c, http.StatusNotFound, "not_found", "user not found", nil)
		return
	}
	c.JSON(http.StatusOK, NewResponseDTO(user))
}

func (h *Handler) Update(c *gin.Context) {
	var input UpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	actorID := getAuthUserID(c)

	user, err := h.service.Update(
		c.Param("id"),
		input.Name,
		input.Email,
		input.StageID,
		input.CurrentStepID,
		input.VisaStatus,
		input.StudentStatus,
		input.NameConsultant,
		input.NameCampus,
		input.NoPhone,
		input.Degree,
		input.NameDegree,
		input.VisaType,
		input.Source,
		input.TranslationQuota,
		actorID,
	)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "update_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewResponseDTO(user))
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "delete_failed", err.Error(), nil)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ListStudents(c *gin.Context) {
	users, err := h.service.ListStudents()
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewResponseListDTO(users))
}

func (h *Handler) PatchQuotaTranslation(c *gin.Context) {
	var input PatchQuotaTranslationDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	user, err := h.service.PatchQuotaTranslation(c.Param("id"), input.TranslationQuota)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "update_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewResponseDTO(user))
}

func (h *Handler) PatchVisaStatus(c *gin.Context) {
	var input PatchVisaStatusDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	user, err := h.service.PatchVisaStatus(c.Param("id"), input.VisaStatus)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "update_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewResponseDTO(user))
}

func (h *Handler) PatchStudentStatus(c *gin.Context) {
	var input PatchStudentStatusDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	actorID := getAuthUserID(c)

	user, err := h.service.PatchStudentStatus(c.Param("id"), input.StudentStatus, actorID)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "update_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewResponseDTO(user))
}

func (h *Handler) PatchDocumentConsent(c *gin.Context) {
	var input PatchDocumentConsentDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	user, err := h.service.PatchDocumentConsent(c.Param("id"), input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "update_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewResponseDTO(user))
}
