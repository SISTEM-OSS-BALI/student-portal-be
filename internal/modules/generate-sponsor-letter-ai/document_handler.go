package generatesponsorletterai

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/httpx"
	"github.com/username/gin-gorm-api/internal/modules/auth"
)

type GeneratedDocumentHandler struct{ service *GeneratedDocumentService }

func NewGeneratedDocumentHandler(service *GeneratedDocumentService) *GeneratedDocumentHandler {
	return &GeneratedDocumentHandler{service: service}
}

func getGeneratedDocumentAuthClaims(c *gin.Context) (*auth.Claims, bool) {
	value, ok := c.Get("auth")
	if !ok {
		return nil, false
	}
	claims, ok := value.(*auth.Claims)
	return claims, ok
}

func (h *GeneratedDocumentHandler) Upsert(c *gin.Context) {
	var input GeneratedDocumentUpsertDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}
	doc, err := h.service.Upsert(input)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "generated_sponsor_letter_ai_document_upsert_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewGeneratedDocumentResponseDTO(doc))
}

func (h *GeneratedDocumentHandler) Template(c *gin.Context) {
	c.JSON(http.StatusOK, NewGeneratedDocumentTemplateDTO())
}

func (h *GeneratedDocumentHandler) SubmitToDirector(c *gin.Context) {
	claims, ok := getGeneratedDocumentAuthClaims(c)
	if !ok {
		httpx.RespondError(c, http.StatusUnauthorized, "unauthorized", "missing auth claims", nil)
		return
	}
	var input SubmitToDirectorDTO
	if err := c.ShouldBindJSON(&input); err != nil && !errors.Is(err, io.EOF) {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}
	doc, err := h.service.SubmitToDirector(c.Param("id"), claims.UserID, claims.Role, input.Note)
	if err != nil {
		statusCode := http.StatusBadRequest
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, ErrSponsorLetterAlreadySubmitted):
			statusCode = http.StatusConflict
		case errors.Is(err, ErrSponsorLetterDirectorNotFound):
			statusCode = http.StatusBadRequest
		}
		httpx.RespondError(c, statusCode, "generated_sponsor_letter_ai_document_submit_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewGeneratedDocumentResponseDTO(doc))
}

func (h *GeneratedDocumentHandler) CancelSubmitToDirector(c *gin.Context) {
	claims, ok := getGeneratedDocumentAuthClaims(c)
	if !ok {
		httpx.RespondError(c, http.StatusUnauthorized, "unauthorized", "missing auth claims", nil)
		return
	}
	var input SubmitToDirectorDTO
	if err := c.ShouldBindJSON(&input); err != nil && !errors.Is(err, io.EOF) {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}
	doc, err := h.service.CancelSubmitToDirector(c.Param("id"), claims.UserID, input.Note)
	if err != nil {
		statusCode := http.StatusBadRequest
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, ErrSponsorLetterSubmissionNotCancelable):
			statusCode = http.StatusConflict
		}
		httpx.RespondError(c, statusCode, "generated_sponsor_letter_ai_document_cancel_submit_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewGeneratedDocumentResponseDTO(doc))
}

func (h *GeneratedDocumentHandler) List(c *gin.Context) {
	studentID := c.Query("student_id")
	docs, err := h.service.ListByStudentID(studentID)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "generated_sponsor_letter_ai_document_list_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewGeneratedDocumentResponseDTOs(docs))
}

func (h *GeneratedDocumentHandler) GetByStudentID(c *gin.Context) {
	doc, err := h.service.GetByStudentID(c.Param("student_id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpx.RespondError(c, http.StatusNotFound, "generated_sponsor_letter_ai_document_not_found", "generated sponsor letter ai document not found", nil)
			return
		}
		httpx.RespondError(c, http.StatusInternalServerError, "generated_sponsor_letter_ai_document_get_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewGeneratedDocumentResponseDTO(doc))
}

func (h *GeneratedDocumentHandler) DownloadPDF(c *gin.Context) {
	download, err := h.service.GetApprovedDownload(c.Param("id"))
	if err != nil {
		statusCode := http.StatusBadRequest
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, ErrSponsorLetterDownloadNotApproved):
			statusCode = http.StatusConflict
		case errors.Is(err, ErrSponsorLetterDownloadUnavailable):
			statusCode = http.StatusNotFound
		}
		httpx.RespondError(c, statusCode, "generated_sponsor_letter_ai_document_download_failed", err.Error(), nil)
		return
	}

	if download.FilePath != nil {
		if _, statErr := os.Stat(*download.FilePath); statErr == nil {
			c.FileAttachment(*download.FilePath, download.FileName)
			return
		}
		if strings.TrimSpace(download.FileURL) == "" {
			httpx.RespondError(c, http.StatusNotFound, "generated_sponsor_letter_ai_document_download_failed", "sponsor letter pdf file not found", nil)
			return
		}
	}

	c.Redirect(http.StatusTemporaryRedirect, download.FileURL)
}
