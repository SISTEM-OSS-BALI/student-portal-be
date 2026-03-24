package answerapprovals

import (
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

func (h *Handler) CreateOrUpdate(c *gin.Context) {
	var input CreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	approval, err := h.service.CreateOrUpdate(input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "create_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, NewResponseDTO(approval))
}

func (h *Handler) List(c *gin.Context) {
	var filter Filter
	if value := c.Query("student_id"); value != "" {
		filter.StudentID = &value
	}
	if value := c.Query("answer_id"); value != "" {
		filter.AnswerID = &value
	}
	if value := c.Query("reviewer_id"); value != "" {
		filter.ReviewerID = &value
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
	approval, err := h.service.GetByID(c.Param("id"))
	if err != nil {
		httpx.RespondError(c, http.StatusNotFound, "not_found", "answer approval not found", nil)
		return
	}
	c.JSON(http.StatusOK, NewResponseDTO(approval))
}

func (h *Handler) Update(c *gin.Context) {
	var input UpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	approval, err := h.service.Update(c.Param("id"), input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "update_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewResponseDTO(approval))
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "delete_failed", err.Error(), nil)
		return
	}
	c.Status(http.StatusNoContent)
}
