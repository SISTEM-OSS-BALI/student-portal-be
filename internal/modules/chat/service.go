package chat

import (
	"errors"
	"strings"
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateConversation(creatorID string, input CreateConversationDTO) (schema.ChatConversation, error) {
	ctype := strings.TrimSpace(strings.ToLower(input.Type))
	if ctype == "" {
		return schema.ChatConversation{}, errors.New("type is required")
	}

	conversation := schema.ChatConversation{
		Type:        ctype,
		Title:       input.Title,
		CreatedByID: creatorID,
	}
	if err := s.repo.CreateConversation(&conversation); err != nil {
		return schema.ChatConversation{}, err
	}

	uniqueMembers := map[string]struct{}{
		creatorID: {},
	}
	for _, id := range input.MemberIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		uniqueMembers[id] = struct{}{}
	}

	members := make([]schema.ChatConversationMember, 0, len(uniqueMembers))
	now := time.Now()
	for id := range uniqueMembers {
		role := "member"
		if id == creatorID {
			role = "admin"
		}
		members = append(members, schema.ChatConversationMember{
			ConversationID: conversation.ID,
			UserID:         id,
			Role:           role,
			JoinedAt:       now,
		})
	}

	if err := s.repo.AddMembers(members); err != nil {
		return schema.ChatConversation{}, err
	}

	conversation.Members = members
	return conversation, nil
}

func (s *Service) ListConversationsByUserID(userID string) ([]schema.ChatConversation, error) {
	return s.repo.ListConversationsByUserID(userID)
}

func (s *Service) IsMember(conversationID, userID string) (bool, error) {
	return s.repo.IsMember(conversationID, userID)
}

func (s *Service) ListMessages(conversationID string, limit, offset int) ([]schema.ChatMessage, error) {
	return s.repo.ListMessages(conversationID, limit, offset)
}

func (s *Service) SendMessage(conversationID, senderID string, input SendMessageDTO) (schema.ChatMessage, error) {
	ok, err := s.repo.IsMember(conversationID, senderID)
	if err != nil {
		return schema.ChatMessage{}, err
	}
	if !ok {
		return schema.ChatMessage{}, errors.New("user is not a member of this conversation")
	}

	msgType := strings.TrimSpace(strings.ToLower(input.Type))
	if msgType == "" {
		msgType = "text"
	}
	if msgType == "text" && (input.Text == nil || strings.TrimSpace(*input.Text) == "") {
		return schema.ChatMessage{}, errors.New("text is required for type text")
	}

	contextType := strings.TrimSpace(input.ContextType)
	if contextType == "" && input.ContextUserID != nil {
		contextType = "student"
	}
	message := schema.ChatMessage{
		ConversationID: conversationID,
		SenderID:       senderID,
		Type:           msgType,
		Text:           input.Text,
		ReplyToID:      input.ReplyToID,
		ContextUserID:  input.ContextUserID,
		ContextType:    contextType,
	}
	if err := s.repo.CreateMessage(&message); err != nil {
		return schema.ChatMessage{}, err
	}
	if len(input.Attachments) > 0 {
		attachments := make([]schema.ChatMessageAttachment, 0, len(input.Attachments))
		for _, attachment := range input.Attachments {
			url := strings.TrimSpace(attachment.URL)
			if url == "" {
				continue
			}
			attachments = append(attachments, schema.ChatMessageAttachment{
				MessageID: message.ID,
				FileURL:   url,
				FileName:  strings.TrimSpace(attachment.Name),
				FileType:  strings.TrimSpace(attachment.MimeType),
				FileSize:  attachment.Size,
			})
		}
		if err := s.repo.AddAttachments(attachments); err != nil {
			return schema.ChatMessage{}, err
		}
	}
	if err := s.repo.AddMentions(message.ID, input.MentionUserIDs); err != nil {
		return schema.ChatMessage{}, err
	}

	return s.repo.GetMessageByID(message.ID)
}

func (s *Service) MarkRead(conversationID, userID string, at time.Time) error {
	ok, err := s.repo.IsMember(conversationID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("user is not a member of this conversation")
	}
	return s.repo.UpdateMemberLastReadAt(conversationID, userID, at)
}

func (s *Service) ListMentions(userID string, limit, offset int) ([]MentionMessageDTO, error) {
	return s.repo.ListMentions(userID, limit, offset)
}

func (s *Service) MarkMentionRead(messageID, userID string, at time.Time) error {
	return s.repo.UpsertMessageStatus(messageID, userID, "read", at)
}
