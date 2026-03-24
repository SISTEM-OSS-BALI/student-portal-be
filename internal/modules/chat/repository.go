package chat

import (
	"errors"
	"strings"
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type Repository interface {
	CreateConversation(conversation *schema.ChatConversation) error
	AddMembers(members []schema.ChatConversationMember) error
	ListConversationsByUserID(userID string) ([]schema.ChatConversation, error)
	GetConversationByID(id string) (schema.ChatConversation, error)
	IsMember(conversationID, userID string) (bool, error)
	CreateMessage(message *schema.ChatMessage) error
	AddAttachments(attachments []schema.ChatMessageAttachment) error
	AddMentions(messageID string, userIDs []string) error
	ListMentions(userID string, limit, offset int) ([]MentionMessageDTO, error)
	UpsertMessageStatus(messageID, userID, status string, at time.Time) error
	GetMessageByID(id string) (schema.ChatMessage, error)
	ListMessages(conversationID string, limit, offset int) ([]schema.ChatMessage, error)
	UpdateMemberLastReadAt(conversationID, userID string, at time.Time) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) CreateConversation(conversation *schema.ChatConversation) error {
	return r.db.Create(conversation).Error
}

func (r *GormRepository) AddMembers(members []schema.ChatConversationMember) error {
	if len(members) == 0 {
		return nil
	}
	return r.db.Create(&members).Error
}

func (r *GormRepository) ListConversationsByUserID(userID string) ([]schema.ChatConversation, error) {
	var conversations []schema.ChatConversation
	if err := r.db.
		Joins("JOIN chat_conversation_members m ON m.conversation_id = chat_conversations.id").
		Where("m.user_id = ? AND m.left_at IS NULL", userID).
		Preload("Members").
		Order("chat_conversations.updated_at desc").
		Find(&conversations).Error; err != nil {
		return nil, err
	}
	return conversations, nil
}

func (r *GormRepository) GetConversationByID(id string) (schema.ChatConversation, error) {
	var conversation schema.ChatConversation
	if err := r.db.Preload("Members").First(&conversation, "id = ?", id).Error; err != nil {
		return schema.ChatConversation{}, err
	}
	return conversation, nil
}

func (r *GormRepository) IsMember(conversationID, userID string) (bool, error) {
	var count int64
	if err := r.db.Model(&schema.ChatConversationMember{}).
		Where("conversation_id = ? AND user_id = ? AND left_at IS NULL", conversationID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormRepository) CreateMessage(message *schema.ChatMessage) error {
	tx := r.db.Begin()
	if err := tx.Create(message).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&schema.ChatConversation{}).
		Where("id = ?", message.ConversationID).
		Update("updated_at", time.Now()).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *GormRepository) AddAttachments(attachments []schema.ChatMessageAttachment) error {
	if len(attachments) == 0 {
		return nil
	}
	return r.db.Create(&attachments).Error
}

func (r *GormRepository) AddMentions(messageID string, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	unique := make(map[string]struct{}, len(userIDs))
	mentions := make([]schema.ChatMessageMention, 0, len(userIDs))
	for _, id := range userIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, ok := unique[id]; ok {
			continue
		}
		unique[id] = struct{}{}
		mentions = append(mentions, schema.ChatMessageMention{
			MessageID: messageID,
			UserID:    id,
		})
	}
	if len(mentions) == 0 {
		return nil
	}
	return r.db.Create(&mentions).Error
}

func (r *GormRepository) ListMentions(userID string, limit, offset int) ([]MentionMessageDTO, error) {
	if strings.TrimSpace(userID) == "" {
		return []MentionMessageDTO{}, nil
	}
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var rows []MentionMessageDTO
	query := `
		SELECT
			m.id,
			m.conversation_id,
			m.sender_id,
			m.type,
			m.text,
			m.created_at,
			m.context_user_id,
			m.context_type,
			sender.name AS sender_name,
			context_user.name AS context_user_name,
			CASE WHEN ms.id IS NULL THEN false ELSE true END AS is_read
		FROM chat_messages m
		JOIN chat_message_mentions mm ON mm.message_id = m.id
		LEFT JOIN chat_message_statuses ms
			ON ms.message_id = m.id AND ms.user_id = ? AND ms.status = 'read'
		LEFT JOIN users sender ON sender.id = m.sender_id
		LEFT JOIN users context_user ON context_user.id = m.context_user_id
		WHERE mm.user_id = ?
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?
	`

	if err := r.db.Raw(query, userID, userID, limit, offset).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *GormRepository) UpsertMessageStatus(messageID, userID, status string, at time.Time) error {
	messageID = strings.TrimSpace(messageID)
	userID = strings.TrimSpace(userID)
	status = strings.TrimSpace(status)
	if messageID == "" || userID == "" || status == "" {
		return nil
	}

	var existing schema.ChatMessageStatus
	err := r.db.Where("message_id = ? AND user_id = ?", messageID, userID).First(&existing).Error
	if err == nil {
		return r.db.Model(&existing).Updates(map[string]interface{}{
			"status": status,
			"at":     at,
		}).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		record := schema.ChatMessageStatus{
			MessageID: messageID,
			UserID:    userID,
			Status:    status,
			At:        at,
		}
		return r.db.Create(&record).Error
	}
	return err
}

func (r *GormRepository) GetMessageByID(id string) (schema.ChatMessage, error) {
	var message schema.ChatMessage
	if err := r.db.Preload("Mentions").Preload("Attachments").First(&message, "id = ?", id).Error; err != nil {
		return schema.ChatMessage{}, err
	}
	return message, nil
}

func (r *GormRepository) ListMessages(conversationID string, limit, offset int) ([]schema.ChatMessage, error) {
	var messages []schema.ChatMessage
	query := r.db.Preload("Mentions").Preload("Attachments").Where("conversation_id = ?", conversationID).Order("created_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *GormRepository) UpdateMemberLastReadAt(conversationID, userID string, at time.Time) error {
	return r.db.Model(&schema.ChatConversationMember{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Update("last_read_at", at).Error
}
