package ticketmessage

import (
	"strings"
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type CreateDTO struct {
	Name           string `json:"name" binding:"required,max=120"`
	UserID         string `json:"user_id" binding:"required"`
	ConversationID string `json:"conversation_id" binding:"required"`
}

type UpdateDTO struct {
	Name           *string `json:"name"`
	UserID         *string `json:"user_id"`
	ConversationID *string `json:"conversation_id"`
}

type UserResponseDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ConversationResponseDTO struct {
	ID    string  `json:"id"`
	Type  string  `json:"type"`
	Title *string `json:"title,omitempty"`
}

type ResponseDTO struct {
	ID             string                   `json:"id"`
	Name           string                   `json:"name"`
	UserID         string                   `json:"user_id"`
	ConversationID string                   `json:"conversation_id"`
	User           *UserResponseDTO         `json:"user,omitempty"`
	Conversation   *ConversationResponseDTO `json:"conversation,omitempty"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
}

func NewResponseDTO(message schema.TicketMessage) ResponseDTO {
	response := ResponseDTO{
		ID:             message.ID,
		Name:           message.Name,
		UserID:         message.UserID,
		ConversationID: message.ConversationID,
		CreatedAt:      message.CreatedAt,
		UpdatedAt:      message.UpdatedAt,
	}

	if message.User != nil {
		response.User = &UserResponseDTO{
			ID:   message.User.ID,
			Name: message.User.Name,
		}
	}

	if message.Conversation != nil {
		response.Conversation = &ConversationResponseDTO{
			ID:    message.Conversation.ID,
			Type:  message.Conversation.Type,
			Title: message.Conversation.Title,
		}
	}

	return response
}

func NewResponseListDTO(messages []schema.TicketMessage) []ResponseDTO {
	resp := make([]ResponseDTO, 0, len(messages))
	for _, message := range messages {
		resp = append(resp, NewResponseDTO(message))
	}
	return resp
}

func (d CreateDTO) Normalize() CreateDTO {
	d.Name = strings.TrimSpace(d.Name)
	d.UserID = strings.TrimSpace(d.UserID)
	d.ConversationID = strings.TrimSpace(d.ConversationID)
	return d
}
