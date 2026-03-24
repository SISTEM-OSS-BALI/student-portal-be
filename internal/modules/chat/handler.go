package chat

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/username/gin-gorm-api/internal/httpx"
	"github.com/username/gin-gorm-api/internal/modules/auth"
)

type Handler struct {
	service      *Service
	socketServer *SocketServer
}

func NewHandler(service *Service, socketServer *SocketServer) *Handler {
	return &Handler{service: service, socketServer: socketServer}
}

func (h *Handler) CreateConversation(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		return
	}

	var input CreateConversationDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	conversation, err := h.service.CreateConversation(claims.UserID, input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "create_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, NewConversationResponseDTO(conversation))
}

func (h *Handler) ListConversations(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		return
	}

	conversations, err := h.service.ListConversationsByUserID(claims.UserID)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewConversationResponseListDTO(conversations))
}

func (h *Handler) ListMessages(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		return
	}

	conversationID := c.Param("id")
	member, err := h.service.IsMember(conversationID, claims.UserID)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "membership_failed", err.Error(), nil)
		return
	}
	if !member {
		httpx.RespondError(c, http.StatusForbidden, "forbidden", "not a member of this conversation", nil)
		return
	}

	limit := parseInt(c.Query("limit"), 50)
	offset := parseInt(c.Query("offset"), 0)

	messages, err := h.service.ListMessages(conversationID, limit, offset)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewMessageResponseListDTO(messages))
}

func (h *Handler) SendMessage(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		return
	}

	var input SendMessageDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	message, err := h.service.SendMessage(c.Param("id"), claims.UserID, input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "send_failed", err.Error(), nil)
		return
	}

	dto := NewMessageResponseDTO(message)
	if h.socketServer != nil {
		h.socketServer.BroadcastMessage(c.Param("id"), dto)
	}
	c.JSON(http.StatusCreated, dto)
}

func (h *Handler) MarkRead(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		return
	}

	var input MarkReadDTO
	if err := c.ShouldBindJSON(&input); err != nil && !errors.Is(err, io.EOF) {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	at := time.Now()
	if input.At != nil {
		at = *input.At
	}

	if err := h.service.MarkRead(c.Param("id"), claims.UserID, at); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "mark_read_failed", err.Error(), nil)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ListMentions(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		return
	}

	limit := parseInt(c.Query("limit"), 50)
	offset := parseInt(c.Query("offset"), 0)

	mentions, err := h.service.ListMentions(claims.UserID, limit, offset)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, mentions)
}

func (h *Handler) MarkMentionRead(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		return
	}

	messageID := c.Param("id")
	if messageID == "" {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", "id is required", nil)
		return
	}

	if err := h.service.MarkMentionRead(messageID, claims.UserID, time.Now()); err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "mark_read_failed", err.Error(), nil)
		return
	}

	c.Status(http.StatusNoContent)
}

func getAuthClaims(c *gin.Context) (*auth.Claims, bool) {
	claimsAny, ok := c.Get("auth")
	if !ok {
		httpx.RespondError(c, http.StatusUnauthorized, "unauthorized", "missing auth claims", nil)
		return nil, false
	}

	claims, ok := claimsAny.(*auth.Claims)
	if !ok {
		httpx.RespondError(c, http.StatusUnauthorized, "unauthorized", "invalid auth claims", nil)
		return nil, false
	}

	return claims, true
}

func parseInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
