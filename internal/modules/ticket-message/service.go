package ticketmessage

import (
	"errors"
	"strings"

	"github.com/username/gin-gorm-api/internal/schema"
)

const defaultTicketStatus = "open"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func normalizeStatus(value string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return "", errors.New("status is required")
	}

	normalized = strings.Join(strings.Fields(normalized), "_")
	if len(normalized) > 20 {
		return "", errors.New("status must be at most 20 characters")
	}

	return normalized, nil
}

func (s *Service) Create(input CreateDTO) (schema.TicketMessage, error) {
	input = input.Normalize()
	if input.Name == "" || input.UserID == "" {
		return schema.TicketMessage{}, errors.New("name and user_id are required")
	}

	status := defaultTicketStatus
	if input.Status != nil {
		normalizedStatus, err := normalizeStatus(*input.Status)
		if err != nil {
			return schema.TicketMessage{}, err
		}
		status = normalizedStatus
	}

	message := schema.TicketMessage{
		Name:           input.Name,
		UserID:         input.UserID,
		ConversationID: input.ConversationID,
		Status:         status,
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
		conversationID := strings.TrimSpace(*input.ConversationID)
		if conversationID == "" {
			message.ConversationID = nil
		} else {
			message.ConversationID = &conversationID
		}
	}
	if input.Status != nil {
		status, err := normalizeStatus(*input.Status)
		if err != nil {
			return schema.TicketMessage{}, err
		}
		message.Status = status
	}

	if message.Name == "" || message.UserID == "" {
		return schema.TicketMessage{}, errors.New("name and user_id are required")
	}

	if err := s.repo.Update(&message); err != nil {
		return schema.TicketMessage{}, err
	}
	return s.repo.GetByID(message.ID)
}

func (s *Service) UpdateStatus(id string, input UpdateStatusDTO) (schema.TicketMessage, error) {
	message, err := s.repo.GetByID(strings.TrimSpace(id))
	if err != nil {
		return schema.TicketMessage{}, err
	}

	status, err := normalizeStatus(input.Status)
	if err != nil {
		return schema.TicketMessage{}, err
	}

	message.Status = status

	if err := s.repo.Update(&message); err != nil {
		return schema.TicketMessage{}, err
	}
	return s.repo.GetByID(message.ID)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(strings.TrimSpace(id))
}

func (s *Service) DeleteWithConversation(id string) error {
	return s.repo.DeleteWithConversation(strings.TrimSpace(id))
}
