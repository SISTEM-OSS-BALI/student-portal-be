package documenttranslations

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

func (h *Handler) Create(c *gin.Context) {
	var input CreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	item, err := h.service.Create(input)
	if err != nil {
		if errors.Is(err, ErrTranslationQuotaExceeded) {
			httpx.RespondError(c, http.StatusBadRequest, "translation_quota_exceeded", err.Error(), nil)
			return
		}
		httpx.RespondError(c, http.StatusBadRequest, "create_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, NewResponseDTO(item))
}

func (h *Handler) List(c *gin.Context) {
	var filter Filter
	if value := c.Query("student_id"); value != "" {
		filter.StudentID = &value
	}
	if value := c.Query("document_id"); value != "" {
		filter.DocumentID = &value
	}
	if value := c.Query("uploader_id"); value != "" {
		filter.UploaderID = &value
	}
	if value := c.Query("answer_document_id"); value != "" {
		filter.AnswerDocumentID = &value
	}
	if value := c.Query("status"); value != "" {
		filter.Status = &value
	}

	items, err := h.service.List(filter)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewResponseListDTO(items))
}

func (h *Handler) GetByID(c *gin.Context) {
	item, err := h.service.GetByID(c.Param("id"))
	if err != nil {
		httpx.RespondError(c, http.StatusNotFound, "not_found", "document translation not found", nil)
		return
	}
	c.JSON(http.StatusOK, NewResponseDTO(item))
}

func (h *Handler) Update(c *gin.Context) {
	var input UpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	item, err := h.service.Update(c.Param("id"), input)
	if err != nil {
		if errors.Is(err, ErrTranslationQuotaExceeded) {
			httpx.RespondError(c, http.StatusBadRequest, "translation_quota_exceeded", err.Error(), nil)
			return
		}
		httpx.RespondError(c, http.StatusBadRequest, "update_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewResponseDTO(item))
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "delete_failed", err.Error(), nil)
		return
	}
	c.Status(http.StatusNoContent)
}
