package generatecvai

import (
	"context"
	"errors"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/username/gin-gorm-api/internal/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Generate(c *gin.Context) {
	start := time.Now()
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("[ai.generate] panic after=%s panic=%v\n%s", time.Since(start), rec, debug.Stack())
			if !c.Writer.Written() {
				httpx.RespondError(c, http.StatusInternalServerError, "generate_panic", "unexpected error while generating CV", nil)
			}
		}
	}()

	var input GenerateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[ai.generate] bind error: %v", err)
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	timeoutSec := getEnvAsInt("AI_GENERATE_TIMEOUT_SEC", 180)
	baseCtx := context.WithoutCancel(c.Request.Context())
	ctx, cancel := context.WithTimeout(baseCtx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	log.Printf(
		"[ai.generate] start student_id=%v answers=%d sections=%d template_path=%v template_url=%v timeout=%ds",
		stringPtrValue(input.StudentID),
		len(input.Answers),
		len(input.Sections),
		stringPtrValue(input.TemplatePath),
		stringPtrValue(input.TemplateURL),
		timeoutSec,
	)

	result, err := h.service.Generate(ctx, input)
	if err != nil {
		log.Printf("[ai.generate] failed after=%s err=%v", time.Since(start), err)

		switch {
		case errors.Is(err, ErrPromptRequired), errors.Is(err, ErrModelRequired):
			httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		case errors.Is(err, ErrOllamaTimeout), errors.Is(err, context.DeadlineExceeded):
			httpx.RespondError(c, http.StatusGatewayTimeout, "ollama_timeout", err.Error(), nil)
		case errors.Is(err, ErrOllamaRequestFailed), errors.Is(err, ErrInvalidOllamaPayload):
			httpx.RespondError(c, http.StatusBadGateway, "ollama_failed", err.Error(), nil)
		default:
			httpx.RespondError(c, http.StatusInternalServerError, "generate_failed", err.Error(), nil)
		}
		return
	}

	log.Printf(
		"[ai.generate] success after=%s model=%s done=%v done_reason=%s response_len=%d file_url=%s",
		time.Since(start),
		result.Model,
		result.Done,
		result.DoneReason,
		len(result.Response),
		result.FileURL,
	)

	c.JSON(http.StatusOK, result)
}

func stringPtrValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
