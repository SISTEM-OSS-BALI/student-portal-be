package country

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

func (h *Handler) Create(c *gin.Context) {
	var input CreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	country, err := h.service.Create(input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "create_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, NewResponseDTO(country, 0, 0))
}

func (h *Handler) List(c *gin.Context) {
	countries, totals, stepTotals, err := h.service.ListWithDocumentAndStepCounts()
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewResponseListDTO(countries, totals, stepTotals))
}

func (h *Handler) GetByID(c *gin.Context) {
	country, total, stepTotal, err := h.service.GetByIDWithDocumentAndStepCount(c.Param("id"))
	if err != nil {
		httpx.RespondError(c, http.StatusNotFound, "not_found", "country not found", nil)
		return
	}
	c.JSON(http.StatusOK, NewResponseDTO(country, total, stepTotal))
}

func (h *Handler) Update(c *gin.Context) {
	var input UpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	country, err := h.service.Update(c.Param("id"), input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "update_failed", err.Error(), nil)
		return
	}

	total, err := h.service.DocumentCount(country.ID)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "count_failed", err.Error(), nil)
		return
	}
	stepTotal, err := h.service.StepCount(country.ID)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "count_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewResponseDTO(country, total, stepTotal))
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "delete_failed", err.Error(), nil)
		return
	}
	c.Status(http.StatusNoContent)
}
