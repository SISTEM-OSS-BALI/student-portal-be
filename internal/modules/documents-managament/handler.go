package documents

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/username/gin-gorm-api/internal/httpx"
	"github.com/username/gin-gorm-api/internal/utils"
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

	doc, err := h.service.Create(input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "create_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, NewResponseDTO(doc))
}

func (h *Handler) List(c *gin.Context) {
	docs, err := h.service.List()
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewResponseListDTO(docs))
}

func (h *Handler) GetByID(c *gin.Context) {
	doc, err := h.service.GetByID(c.Param("id"))
	if err != nil {
		httpx.RespondError(c, http.StatusNotFound, "not_found", "document not found", nil)
		return
	}
	c.JSON(http.StatusOK, NewResponseDTO(doc))
}

func (h *Handler) Update(c *gin.Context) {
	var input UpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	doc, err := h.service.Update(c.Param("id"), input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "update_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewResponseDTO(doc))
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "delete_failed", err.Error(), nil)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DocumentTranslationRequired(c *gin.Context) {
	docs, err := h.service.DocumentTranslationRequired()
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "document_translation_required_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewResponseListDTO(docs))
}

func (h *Handler) CountPDFPages(c *gin.Context) {
	var input CountPDFPagesRequestDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	count, err := utils.CountPDFPagesFromURL(input.URL, nil)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "count_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, CountPDFPagesResponseDTO{
		URL:       input.URL,
		PageCount: count,
	})
}

func (h *Handler) MergePDF(c *gin.Context) {
	var input MergePDFRequestDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	merged, err := utils.MergePDFsFromURLs([]string{input.OriginalURL, input.TranslationURL}, &utils.PDFMergeOptions{
		AddDividerPage: input.AddDividerPage,
	})
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "merge_failed", err.Error(), nil)
		return
	}

	fileName := strings.TrimSpace(input.FileName)
	if fileName == "" {
		fileName = "merged.pdf"
	}
	if !strings.HasSuffix(strings.ToLower(fileName), ".pdf") {
		fileName = fmt.Sprintf("%s.pdf", fileName)
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	c.Data(http.StatusOK, "application/pdf", merged)
}
