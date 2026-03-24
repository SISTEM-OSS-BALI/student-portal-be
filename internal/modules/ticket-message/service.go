package ticketmessage

import (
	"errors"
	"strings"

	"github.com/username/gin-gorm-api/internal/schema"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(input CreateDTO) (schema.TicketMessage, error) {
	input = input.Normalize()
	if input.Name == "" || input.UserID == "" || input.ConversationID == "" {
		return schema.TicketMessage{}, errors.New("name, user_id, and conversation_id are required")
	}

	message := schema.TicketMessage{
		Name:           input.Name,
		UserID:         input.UserID,
		ConversationID: input.ConversationID,
	}
	if err := s.repo.Create(&message); err != nil {
		return schema.TicketMessage{}, err
	}
	return s.repo.GetByID(message.ID)
}

func (s *Service) List() ([]schema.TicketMessage, error) {
	return s.repo.List()
}

func (s *Service) ListByConversationID(conversationID string) ([]schema.TicketMessage, error) {
	return s.repo.ListByConversationID(strings.TrimSpace(conversationID))
}

func (s *Service) ListByUserID(userID string) ([]schema.TicketMessage, error) {
	return s.repo.ListByUserID(strings.TrimSpace(userID))
}

func (s *Service) GetByID(id string) (schema.TicketMessage, error) {
	return s.repo.GetByID(strings.TrimSpace(id))
}

func (s *Service) Update(id string, input UpdateDTO) (schema.TicketMessage, error) {
	message, err := s.repo.GetByID(strings.TrimSpace(id))
	if err != nil {
		return schema.TicketMessage{}, err
	}

	if input.Name != nil {
		message.Name = strings.TrimSpace(*input.Name)
	}
	if input.UserID != nil {
		message.UserID = strings.TrimSpace(*input.UserID)
	}
	if input.ConversationID != nil {
		message.ConversationID = strings.TrimSpace(*input.ConversationID)
	}

	if message.Name == "" || message.UserID == "" || message.ConversationID == "" {
		return schema.TicketMessage{}, errors.New("name, user_id, and conversation_id are required")
	}

	if err := s.repo.Update(&message); err != nil {
		return schema.TicketMessage{}, err
	}
	return s.repo.GetByID(message.ID)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(strings.TrimSpace(id))
}
