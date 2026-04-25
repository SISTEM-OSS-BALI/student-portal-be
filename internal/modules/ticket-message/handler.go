package ticketmessage

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/username/gin-gorm-api/internal/httpx"
	"github.com/username/gin-gorm-api/internal/schema"
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

	message, err := h.service.Create(input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "ticket_message_create_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, NewResponseDTO(message))
}

func (h *Handler) List(c *gin.Context) {
	conversationID := c.Query("conversation_id")
	userID := c.Query("user_id")

	var (
		messages []schema.TicketMessage
		err      error
	)

	switch {
	case conversationID != "":
		messages, err = h.service.ListByConversationID(conversationID)
	case userID != "":
		messages, err = h.service.ListByUserID(userID)
	default:
		messages, err = h.service.List()
	}

	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "ticket_message_list_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewResponseListDTO(messages))
}

func (h *Handler) GetByID(c *gin.Context) {
	message, err := h.service.GetByID(c.Param("id"))
	if err != nil {
		httpx.RespondError(c, http.StatusNotFound, "ticket_message_not_found", "ticket message not found", nil)
		return
	}

	c.JSON(http.StatusOK, NewResponseDTO(message))
}

func (h *Handler) Update(c *gin.Context) {
	var input UpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	message, err := h.service.Update(c.Param("id"), input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "ticket_message_update_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewResponseDTO(message))
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	var input UpdateStatusDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	message, err := h.service.UpdateStatus(c.Param("id"), input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "ticket_message_status_update_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewResponseDTO(message))
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "ticket_message_delete_failed", err.Error(), nil)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteWithConversation(c *gin.Context) {
	if err := h.service.DeleteWithConversation(c.Param("id")); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "ticket_message_delete_with_conversation_failed", err.Error(), nil)
		return
	}

	c.Status(http.StatusNoContent)
}
