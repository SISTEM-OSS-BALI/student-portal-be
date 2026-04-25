package informationcountrymanagement

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

	data, err := h.service.Create(
		input.Slug,
		input.Title,
		input.Description,
		input.Priority,
		input.CountryID,
	)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "create_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, NewResponseDTO(data))
}

func (h *Handler) List(c *gin.Context) {
	if countryID := c.Query("country_id"); countryID != "" {
		data, err := h.service.ListByCountryID(countryID)
		if err != nil {
			httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
			return
		}

		c.JSON(http.StatusOK, NewResponseListDTO(data))
		return
	}

	if slug := c.Query("slug"); slug != "" {
		data, err := h.service.GetBySlug(slug)
		if err != nil {
			httpx.RespondError(c, http.StatusNotFound, "not_found", "information country management not found", nil)
			return
		}

		c.JSON(http.StatusOK, NewResponseDTO(data))
		return
	}

	data, err := h.service.List()
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewResponseListDTO(data))
}

func (h *Handler) GetByID(c *gin.Context) {
	data, err := h.service.GetByID(c.Param("id"))
	if err != nil {
		httpx.RespondError(c, http.StatusNotFound, "not_found", "information country management not found", nil)
		return
	}

	c.JSON(http.StatusOK, NewResponseDTO(data))
}

func (h *Handler) GetBySlug(c *gin.Context) {
	data, err := h.service.GetBySlug(c.Param("slug"))
	if err != nil {
		httpx.RespondError(c, http.StatusNotFound, "not_found", "information country management not found", nil)
		return
	}

	c.JSON(http.StatusOK, NewResponseDTO(data))
}

func (h *Handler) ListByCountryID(c *gin.Context) {
	data, err := h.service.ListByCountryID(c.Param("country_id"))
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewResponseListDTO(data))
}

func (h *Handler) Update(c *gin.Context) {
	var input UpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	data, err := h.service.Update(
		c.Param("id"),
		input.Slug,
		input.Title,
		input.Description,
		input.Priority,
		input.CountryID,
	)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "update_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewResponseDTO(data))
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "delete_failed", err.Error(), nil)
		return
	}

	c.Status(http.StatusNoContent)
}
