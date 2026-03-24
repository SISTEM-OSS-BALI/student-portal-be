package schema

import (
	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

func setCUIDIfEmpty(id *string) {
	if id == nil {
		return
	}
	if *id == "" {
		*id = cuid.New()
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&u.ID)
	return nil
}

func (d *DocumentsManagement) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&d.ID)
	return nil
}

func (c *CountryManagement) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&c.ID)
	return nil
}

func (s *StageManagement) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&s.ID)
	return nil
}

func (s *StepsManagement) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&s.ID)
	return nil
}

func (c *ChildStepsManagement) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&c.ID)
	return nil
}

func (c *CountryStepsManagement) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&c.ID)
	return nil
}

func (n *NoteStudent) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&n.ID)
	return nil
}

func (c *ChatConversation) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&c.ID)
	return nil
}

func (m *ChatConversationMember) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&m.ID)
	return nil
}

func (m *ChatMessage) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&m.ID)
	return nil
}

func (a *ChatMessageAttachment) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&a.ID)
	return nil
}

func (s *ChatMessageStatus) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&s.ID)
	return nil
}

func (m *ChatMessageMention) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&m.ID)
	return nil
}

func (b *QuestionBase) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&b.ID)
	return nil
}

func (q *Question) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&q.ID)
	return nil
}

func (o *QuestionOption) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&o.ID)
	return nil
}

func (a *AnswerQuestion) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&a.ID)
	return nil
}

func (s *AnswerSubmission) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&s.ID)
	return nil
}

func (d *AnswerDocument) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&d.ID)
	return nil
}

func (d *GeneratedCVAIDocument) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&d.ID)
	return nil
}

func (d *GeneratedStatementLetterAIDocument) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&d.ID)
	return nil
}

func (a *AnswerApproval) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&a.ID)
	return nil
}

func (a *AnswerDocumentApproval) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&a.ID)
	return nil
}

func (d *DocumentTranslation) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&d.ID)
	return nil
}

func (t *TicketMessage) BeforeCreate(tx *gorm.DB) error {
	setCUIDIfEmpty(&t.ID)
	return nil
}
