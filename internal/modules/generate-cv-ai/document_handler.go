package generatecvai

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/httpx"
)

type GeneratedDocumentHandler struct {
	service *GeneratedDocumentService
}

func NewGeneratedDocumentHandler(service *GeneratedDocumentService) *GeneratedDocumentHandler {
	return &GeneratedDocumentHandler{service: service}
}

func (h *GeneratedDocumentHandler) Upsert(c *gin.Context) {
	var input GeneratedDocumentUpsertDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	doc, err := h.service.Upsert(input)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "generated_cv_ai_document_upsert_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewGeneratedDocumentResponseDTO(doc))
}

func (h *GeneratedDocumentHandler) List(c *gin.Context) {
	studentID := c.Query("student_id")
	docs, err := h.service.ListByStudentID(studentID)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "generated_cv_ai_document_list_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewGeneratedDocumentResponseListDTO(docs))
}

func (h *GeneratedDocumentHandler) GetByStudentID(c *gin.Context) {
	doc, err := h.service.GetByStudentID(c.Param("student_id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpx.RespondError(c, http.StatusNotFound, "generated_cv_ai_document_not_found", "generated cv ai document not found", nil)
			return
		}
		httpx.RespondError(c, http.StatusInternalServerError, "generated_cv_ai_document_get_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewGeneratedDocumentResponseDTO(doc))
}
