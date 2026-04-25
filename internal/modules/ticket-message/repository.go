package ticketmessage

import (
	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type Repository interface {
	Create(message *schema.TicketMessage) error
	List() ([]schema.TicketMessage, error)
	ListByConversationID(conversationID string) ([]schema.TicketMessage, error)
	ListByUserID(userID string) ([]schema.TicketMessage, error)
	GetByID(id string) (schema.TicketMessage, error)
	Update(message *schema.TicketMessage) error
	Delete(id string) error
	DeleteWithConversation(id string) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) baseQuery() *gorm.DB {
	return r.db.Preload("User").Preload("Conversation")
}

func (r *GormRepository) Create(message *schema.TicketMessage) error {
	return r.db.Create(message).Error
}

func (r *GormRepository) List() ([]schema.TicketMessage, error) {
	var messages []schema.TicketMessage
	if err := r.baseQuery().Order("created_at desc").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *GormRepository) ListByConversationID(conversationID string) ([]schema.TicketMessage, error) {
	var messages []schema.TicketMessage
	if err := r.baseQuery().
		Where("conversation_id = ?", conversationID).
		Order("created_at asc").
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *GormRepository) ListByUserID(userID string) ([]schema.TicketMessage, error) {
	var messages []schema.TicketMessage
	if err := r.baseQuery().
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *GormRepository) GetByID(id string) (schema.TicketMessage, error) {
	var message schema.TicketMessage
	if err := r.baseQuery().Where("id = ?", id).First(&message).Error; err != nil {
		return schema.TicketMessage{}, err
	}
	return message, nil
}

func (r *GormRepository) Update(message *schema.TicketMessage) error {
	return r.db.Save(message).Error
}

func (r *GormRepository) Delete(id string) error {
	return r.db.Delete(&schema.TicketMessage{}, "id = ?", id).Error
}

func (r *GormRepository) DeleteWithConversation(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var message schema.TicketMessage
		if err := tx.First(&message, "id = ?", id).Error; err != nil {
			return err
		}

		if message.ConversationID != nil && *message.ConversationID != "" {
			return tx.Delete(&schema.ChatConversation{}, "id = ?", *message.ConversationID).Error
		}

		return tx.Delete(&schema.TicketMessage{}, "id = ?", id).Error
	})
}
