package chat

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type CreateConversationDTO struct {
	Type      string   `json:"type" binding:"required"` // "direct" | "group"
	Title     *string  `json:"title"`
	MemberIDs []string `json:"member_ids" binding:"required,min=1"`
}

type SendMessageDTO struct {
	Type           string                     `json:"type"` // "text" | "image" | "file" | "system"
	Text           *string                    `json:"text"`
	ReplyToID      *string                    `json:"reply_to_id"`
	MentionUserIDs []string                   `json:"mention_user_ids"`
	ContextUserID  *string                    `json:"context_user_id"`
	ContextType    string                     `json:"context_type"`
	Attachments    []SendMessageAttachmentDTO `json:"attachments"`
}

type MarkReadDTO struct {
	At *time.Time `json:"at"`
}

type ConversationResponseDTO struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Title       *string   `json:"title"`
	CreatedByID string    `json:"created_by_id"`
	MemberIDs   []string  `json:"member_ids"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MessageResponseDTO struct {
	ID             string                 `json:"id"`
	ConversationID string                 `json:"conversation_id"`
	SenderID       string                 `json:"sender_id"`
	SenderName     string                 `json:"sender_name,omitempty"`
	SenderRole     string                 `json:"sender_role,omitempty"`
	Type           string                 `json:"type"`
	Text           *string                `json:"text"`
	ReplyToID      *string                `json:"reply_to_id"`
	MentionUserIDs []string               `json:"mention_user_ids"`
	ContextUserID  *string                `json:"context_user_id,omitempty"`
	ContextType    string                 `json:"context_type,omitempty"`
	Attachments    []MessageAttachmentDTO `json:"attachments,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	EditedAt       *time.Time             `json:"edited_at"`
	DeletedAt      *time.Time             `json:"deleted_at"`
}

type MentionMessageDTO struct {
	ID              string    `json:"id"`
	ConversationID  string    `json:"conversation_id"`
	SenderID        string    `json:"sender_id"`
	SenderName      string    `json:"sender_name,omitempty"`
	Type            string    `json:"type"`
	Text            *string   `json:"text"`
	CreatedAt       time.Time `json:"created_at"`
	ContextUserID   *string   `json:"context_user_id,omitempty"`
	ContextType     string    `json:"context_type,omitempty"`
	ContextUserName *string   `json:"context_user_name,omitempty"`
	IsRead          bool      `json:"is_read"`
}

type SendMessageAttachmentDTO struct {
	URL      string `json:"url"`
	Name     string `json:"name"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
}

type MessageAttachmentDTO struct {
	URL      string `json:"url"`
	Name     string `json:"name"`
	MimeType string `json:"mimeType,omitempty"`
	Size     int64  `json:"size,omitempty"`
}

func NewConversationResponseDTO(conversation schema.ChatConversation) ConversationResponseDTO {
	memberIDs := make([]string, 0, len(conversation.Members))
	for _, member := range conversation.Members {
		memberIDs = append(memberIDs, member.UserID)
	}
	return ConversationResponseDTO{
		ID:          conversation.ID,
		Type:        conversation.Type,
		Title:       conversation.Title,
		CreatedByID: conversation.CreatedByID,
		MemberIDs:   memberIDs,
		CreatedAt:   conversation.CreatedAt,
		UpdatedAt:   conversation.UpdatedAt,
	}
}

func NewConversationResponseListDTO(conversations []schema.ChatConversation) []ConversationResponseDTO {
	out := make([]ConversationResponseDTO, 0, len(conversations))
	for _, conversation := range conversations {
		out = append(out, NewConversationResponseDTO(conversation))
	}
	return out
}

func NewMessageResponseDTO(message schema.ChatMessage) MessageResponseDTO {
	mentionIDs := make([]string, 0, len(message.Mentions))
	for _, mention := range message.Mentions {
		mentionIDs = append(mentionIDs, mention.UserID)
	}
	attachments := make([]MessageAttachmentDTO, 0, len(message.Attachments))
	for _, attachment := range message.Attachments {
		attachments = append(attachments, MessageAttachmentDTO{
			URL:      attachment.FileURL,
			Name:     attachment.FileName,
			MimeType: attachment.FileType,
			Size:     attachment.FileSize,
		})
	}
	return MessageResponseDTO{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		SenderName: func() string {
			if message.Sender == nil {
				return ""
			}
			return message.Sender.Name
		}(),
		SenderRole: func() string {
			if message.Sender == nil {
				return ""
			}
			return string(message.Sender.Role)
		}(),
		Type:           message.Type,
		Text:           message.Text,
		ReplyToID:      message.ReplyToID,
		MentionUserIDs: mentionIDs,
		ContextUserID:  message.ContextUserID,
		ContextType:    message.ContextType,
		Attachments:    attachments,
		CreatedAt:      message.CreatedAt,
		EditedAt:       message.EditedAt,
		DeletedAt:      message.DeletedAt,
	}
}

func NewMessageResponseListDTO(messages []schema.ChatMessage) []MessageResponseDTO {
	out := make([]MessageResponseDTO, 0, len(messages))
	for _, message := range messages {
		out = append(out, NewMessageResponseDTO(message))
	}
	return out
}
