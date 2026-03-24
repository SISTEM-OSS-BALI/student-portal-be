package schema

import "time"

type UserRole string
type TranslationNeeded string
type AutoRenamePattern string
type StatusStudent string
type QuestionType string

const (
	UserRoleStudent   UserRole = "STUDENT"
	UserRoleDirector  UserRole = "DIRECTOR"
	UserRoleAdmission UserRole = "ADMISSION"
)

const (
	TranslationNeededYes TranslationNeeded = "YES"
	TranslationNeededNo  TranslationNeeded = "NO"
)

const (
	AutoRenamePatternDate       AutoRenamePattern = "DATE"
	AutoRenamePatternDocumentID AutoRenamePattern = "DOCUMENT_ID"
	AutoRenamePatternNone       AutoRenamePattern = "NONE"
)

const (
	StatusStudentOnGoing  StatusStudent = "ON GOING"
	StatusStudentOnShore  StatusStudent = "ON SHORE"
	StatusStudentPostPone StatusStudent = "POSTPONE"
	StatusStudentCancel   StatusStudent = "CANCEL"
)

type User struct {
	ID                         string                              `json:"id" gorm:"primaryKey;size:25"`
	Name                       string                              `json:"name" gorm:"size:120;not null"`
	Email                      string                              `json:"email" gorm:"size:191;uniqueIndex;not null"`
	Password                   string                              `json:"-" gorm:"size:191;not null"`
	Role                       UserRole                            `json:"role" gorm:"size:50;not null"`
	CreatedAt                  time.Time                           `json:"created_at"`
	NoPhone                    *string                             `json:"no_phone,omitempty" gorm:"size:20"`
	UpdatedAt                  time.Time                           `json:"updated_at"`
	StageID                    *string                             `json:"stage_id,omitempty" gorm:"size:25;index"`
	Stage                      *StageManagement                    `json:"stage,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Status                     StatusStudent                       `json:"status" gorm:"size:20;not null;default:'ON GOING'"`
	NameCampus                 *string                             `json:"name_campus" gorm:"size:100"`
	Degree                     *string                             `json:"degree" gorm:"size:100"`
	NameDegree                 *string                             `json:"name_degree" gorm:"size:100"`
	VisaType                   *string                             `json:"visa_type" gorm:"size:50"`
	TranslationQuota           int                                 `json:"translation_quota" gorm:"not null;default:0"`
	NotesStudent               []NoteStudent                       `json:"notes,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GeneratedCVAI              *GeneratedCVAIDocument              `json:"generated_cv_ai,omitempty" gorm:"foreignKey:StudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GeneratedStatementLetterAI *GeneratedStatementLetterAIDocument `json:"generated_statement_letter_ai,omitempty" gorm:"foreignKey:StudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type NoteStudent struct {
	ID        string    `json:"id" gorm:"primaryKey;size:25"`
	UserID    string    `json:"user_id" gorm:"size:25;not null;index"`
	User      *User     `json:"user,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DocumentsManagement struct {
	ID                string            `json:"id" gorm:"primaryKey;size:25"`
	Label             string            `json:"label" gorm:"size:120;not null"`
	InternalCode      string            `json:"internal_code" gorm:"size:50;not null;uniqueIndex"`
	FileType          string            `json:"file_type" gorm:"size:50;not null"`
	Category          string            `json:"category" gorm:"size:100;not null"`
	TranslationNeeded TranslationNeeded `json:"translation_needed" gorm:"size:10;not null"`
	Required          bool              `json:"required" gorm:"not null"`
	AutoRenamePattern AutoRenamePattern `json:"auto_rename_pattern" gorm:"size:191"`
	Notes             string            `json:"notes" gorm:"type:text"`
	Stages            []StageManagement `json:"stages,omitempty" gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

type StepsManagement struct {
	ID           string                   `json:"id" gorm:"primaryKey;size:25"`
	Label        string                   `json:"label" gorm:"size:120;not null"`
	Children     []ChildStepsManagement   `json:"children,omitempty" gorm:"many2many:steps_children;"`
	CountrySteps []CountryStepsManagement `json:"country_steps,omitempty" gorm:"foreignKey:StepID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
}

type ChildStepsManagement struct {
	ID        string            `json:"id" gorm:"primaryKey;size:25"`
	Label     string            `json:"label" gorm:"size:120;not null"`
	Steps     []StepsManagement `json:"steps,omitempty" gorm:"many2many:steps_children;"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type CountryManagement struct {
	ID           string                   `json:"id" gorm:"primaryKey;size:25"`
	NameCountry  string                   `json:"name" gorm:"size:120;not null"`
	Stages       []StageManagement        `json:"stages,omitempty" gorm:"foreignKey:CountryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CountrySteps []CountryStepsManagement `json:"country_steps,omitempty" gorm:"foreignKey:CountryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
}

type StageManagement struct {
	ID         string               `json:"id" gorm:"primaryKey;size:25"`
	CountryID  string               `json:"country_id" gorm:"size:25;not null;index"`
	DocumentID string               `json:"document_id" gorm:"size:25;not null;index"`
	Country    *CountryManagement   `json:"country,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Document   *DocumentsManagement `json:"document,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Users      []User               `json:"users,omitempty" gorm:"foreignKey:StageID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
}

type CountryStepsManagement struct {
	ID        string             `json:"id" gorm:"primaryKey;size:25"`
	CountryID string             `json:"country_id" gorm:"size:25;not null;index"`
	StepID    string             `json:"step_id" gorm:"size:25;not null;index"`
	Country   *CountryManagement `json:"country,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Step      *StepsManagement   `json:"step,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type QuestionBase struct {
	ID                       string             `json:"id" gorm:"primaryKey;size:25"`
	Name                     string             `json:"name" gorm:"size:120;not null"`
	Desc                     *string            `json:"desc,omitempty" gorm:"type:text"`
	TypeCountry              string             `json:"type" gorm:"column:type_country;size:50;not null"`
	CountryID                *string            `json:"country_id,omitempty" gorm:"size:25;index"`
	Country                  *CountryManagement `json:"country,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Questions                []Question         `json:"questions,omitempty" gorm:"foreignKey:BaseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AllowMultipleSubmissions bool               `json:"allow_multiple_submissions" gorm:"not null;default:false"`
	CreatedAt                time.Time          `json:"created_at"`
	UpdatedAt                time.Time          `json:"updated_at"`
	Active                   bool               `json:"active" gorm:"not null;default:true"`
	Version                  int                `json:"version" gorm:"not null;default:1"`
}

type Question struct {
	ID          string           `json:"id" gorm:"primaryKey;size:25"`
	BaseID      string           `json:"base_id" gorm:"size:25;not null;index"`
	Text        string           `json:"text" gorm:"type:text;not null"`
	InputType   QuestionType     `json:"input_type" gorm:"size:50;not null"`
	Required    bool             `json:"required" gorm:"not null;default:true"`
	Order       int              `json:"order" gorm:"not null;default:0"`
	HelpText    *string          `json:"help_text,omitempty" gorm:"type:text"`
	Placeholder *string          `json:"placeholder,omitempty" gorm:"type:text"`
	MinLength   *int             `json:"min_length,omitempty"`
	MaxLength   *int             `json:"max_length,omitempty"`
	Base        *QuestionBase    `json:"base,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Options     []QuestionOption `json:"options,omitempty" gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Answers     []AnswerQuestion `json:"answers,omitempty" gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Active      bool             `json:"active" gorm:"not null;default:true"`
}

type QuestionOption struct {
	ID         string                 `json:"id" gorm:"primaryKey;size:25"`
	QuestionID string                 `json:"question_id" gorm:"size:25;not null;index;uniqueIndex:uniq_question_option;index:idx_question_option_order"`
	Label      string                 `json:"label" gorm:"not null"`
	Value      string                 `json:"value" gorm:"not null;uniqueIndex:uniq_question_option"`
	Order      int                    `json:"order" gorm:"not null;default:0;index:idx_question_option_order"`
	Question   *Question              `json:"question,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SelectedBy []AnswerSelectedOption `json:"selected_by,omitempty" gorm:"foreignKey:OptionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
	Active     bool                   `json:"active" gorm:"not null;default:true"`
}

type ChatConversation struct {
	ID          string    `json:"id" gorm:"primaryKey;size:25"`
	Type        string    `json:"type" gorm:"size:16;not null"`
	Title       *string   `json:"title" gorm:"size:120"`
	CreatedByID string    `json:"created_by_id" gorm:"size:25;not null;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Members  []ChatConversationMember `json:"members,omitempty" gorm:"foreignKey:ConversationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Messages []ChatMessage            `json:"messages,omitempty" gorm:"foreignKey:ConversationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ChatConversationMember struct {
	ID             string     `json:"id" gorm:"primaryKey;size:25"`
	ConversationID string     `json:"conversation_id" gorm:"size:25;not null;index"`
	UserID         string     `json:"user_id" gorm:"size:25;not null;index"`
	Role           string     `json:"role" gorm:"size:16;not null"`
	JoinedAt       time.Time  `json:"joined_at" gorm:"not null"`
	LeftAt         *time.Time `json:"left_at,omitempty"`
	MutedUntil     *time.Time `json:"muted_until,omitempty"`
	LastReadAt     *time.Time `json:"last_read_at,omitempty"`
}

type ChatMessage struct {
	ID             string     `json:"id" gorm:"primaryKey;size:25"`
	ConversationID string     `json:"conversation_id" gorm:"size:25;not null;index"`
	SenderID       string     `json:"sender_id" gorm:"size:25;not null;index"`
	Type           string     `json:"type" gorm:"size:16;not null"`
	Text           *string    `json:"text,omitempty" gorm:"type:text"`
	ReplyToID      *string    `json:"reply_to_id,omitempty" gorm:"size:25;index"`
	ContextUserID  *string    `json:"context_user_id,omitempty" gorm:"size:25;index"`
	ContextType    string     `json:"context_type,omitempty" gorm:"size:32"`
	CreatedAt      time.Time  `json:"created_at"`
	EditedAt       *time.Time `json:"edited_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	Attachments []ChatMessageAttachment `json:"attachments,omitempty" gorm:"foreignKey:MessageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Statuses    []ChatMessageStatus     `json:"statuses,omitempty" gorm:"foreignKey:MessageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Mentions    []ChatMessageMention    `json:"mentions,omitempty" gorm:"foreignKey:MessageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ChatMessageAttachment struct {
	ID        string    `json:"id" gorm:"primaryKey;size:25"`
	MessageID string    `json:"message_id" gorm:"size:25;not null;index"`
	FileURL   string    `json:"file_url" gorm:"size:512;not null"`
	FileName  string    `json:"file_name" gorm:"size:255"`
	FileType  string    `json:"file_type" gorm:"size:64"`
	FileSize  int64     `json:"file_size"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatMessageStatus struct {
	ID        string    `json:"id" gorm:"primaryKey;size:25"`
	MessageID string    `json:"message_id" gorm:"size:25;not null;index"`
	UserID    string    `json:"user_id" gorm:"size:25;not null;index"`
	Status    string    `json:"status" gorm:"size:16;not null"`
	At        time.Time `json:"at" gorm:"not null;index"`
}

type ChatMessageMention struct {
	ID        string    `json:"id" gorm:"primaryKey;size:25"`
	MessageID string    `json:"message_id" gorm:"size:25;not null;index"`
	UserID    string    `json:"user_id" gorm:"size:25;not null;index"`
	CreatedAt time.Time `json:"created_at"`
}

// Conversation: 1-1 atau group
type Conversation struct {
	ID          string  `gorm:"primaryKey;size:26"`
	Type        string  `gorm:"size:16;not null"` // "direct" | "group"
	Title       *string `gorm:"size:120"`         // untuk group
	CreatedByID string  `gorm:"size:26;not null;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Members  []ConversationMember `gorm:"foreignKey:ConversationID"`
	Messages []Message            `gorm:"foreignKey:ConversationID"`
}

type ConversationMember struct {
	ID             string    `gorm:"primaryKey;size:26"`
	ConversationID string    `gorm:"size:26;not null;index"`
	UserID         string    `gorm:"size:26;not null;index"`
	Role           string    `gorm:"size:16;not null"` // "member" | "admin"
	JoinedAt       time.Time `gorm:"not null"`
	LeftAt         *time.Time
	MutedUntil     *time.Time
	LastReadAt     *time.Time

	// constraint unik agar user tidak join dua kali
	// gorm can't declare composite unique via tag; use migration or `gorm:"uniqueIndex:uniq_member"`
}

type Message struct {
	ID             string    `gorm:"primaryKey;size:26"`
	ConversationID string    `gorm:"size:26;not null;index"`
	SenderID       string    `gorm:"size:26;not null;index"`
	Type           string    `gorm:"size:16;not null"` // "text" | "image" | "file" | "system"
	Text           *string   `gorm:"type:text"`
	ReplyToID      *string   `gorm:"size:26;index"`
	CreatedAt      time.Time `gorm:"index"`
	EditedAt       *time.Time
	DeletedAt      *time.Time `gorm:"index"`

	Attachments []MessageAttachment `gorm:"foreignKey:MessageID"`
	Statuses    []MessageStatus     `gorm:"foreignKey:MessageID"`
}

type MessageAttachment struct {
	ID        string `gorm:"primaryKey;size:26"`
	MessageID string `gorm:"size:26;not null;index"`
	FileURL   string `gorm:"size:512;not null"`
	FileName  string `gorm:"size:255"`
	FileType  string `gorm:"size:64"` // mime
	FileSize  int64
	CreatedAt time.Time
}

type MessageStatus struct {
	ID        string    `gorm:"primaryKey;size:26"`
	MessageID string    `gorm:"size:26;not null;index"`
	UserID    string    `gorm:"size:26;not null;index"`
	Status    string    `gorm:"size:16;not null"` // "sent" | "delivered" | "read"
	At        time.Time `gorm:"not null;index"`
}

type AnswerSubmission struct {
	ID          string     `json:"id" gorm:"primaryKey;size:25"`
	BaseID      string     `json:"base_id" gorm:"size:25;not null;index"`
	StudentID   string     `json:"student_id" gorm:"size:25;not null;index"`
	Status      string     `json:"status" gorm:"size:20;not null;default:'draft'"`
	Version     int        `json:"version" gorm:"not null;default:1"`
	SubmittedAt *time.Time `json:"submitted_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type AnswerQuestion struct {
	ID           string    `json:"id" gorm:"primaryKey;size:25"`
	SubmissionID *string   `json:"submission_id" gorm:"size:25;index"` // nullable for backward compatibility
	QuestionID   string    `json:"question_id" gorm:"size:25;not null;index"`
	AnswerText   *string   `json:"answer_text" gorm:"type:text"`
	StudentID    *string   `json:"student_id" gorm:"size:25;index"` // kalau masih dipakai
	CreatedAt    time.Time `json:"created_at"`
}

type AnswerSelectedOption struct {
	AnswerID string `json:"answer_id" gorm:"size:25;not null;index"`
	OptionID string `json:"option_id" gorm:"size:25;not null;index"`
}

type AnswerDocument struct {
	ID           string               `json:"id" gorm:"primaryKey;size:25"`
	SubmissionID *string              `json:"submission_id,omitempty" gorm:"size:25;index"`
	StudentID    string               `json:"student_id" gorm:"size:25;not null;index"`
	DocumentID   string               `json:"document_id" gorm:"size:25;not null;index"`
	Document     *DocumentsManagement `json:"document,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	FileURL      string               `json:"file_url" gorm:"type:text;not null"`
	FilePath     *string              `json:"file_path,omitempty" gorm:"type:text"`
	FileName     *string              `json:"file_name,omitempty" gorm:"size:191"`
	FileType     *string              `json:"file_type,omitempty" gorm:"size:50"`
	Status       *string              `json:"status,omitempty" gorm:"size:20"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}

type GeneratedCVAIDocument struct {
	ID           string    `json:"id" gorm:"primaryKey;size:25"`
	StudentID    string    `json:"student_id" gorm:"size:25;not null;uniqueIndex"`
	Student      *User     `json:"student,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:StudentID;references:ID"`
	FileURL      string    `json:"file_url" gorm:"type:text;not null"`
	FilePath     *string   `json:"file_path,omitempty" gorm:"type:text"`
	FileName     *string   `json:"file_name,omitempty" gorm:"size:191"`
	FileType     *string   `json:"file_type,omitempty" gorm:"size:100"`
	WordFileURL  *string   `json:"word_file_url,omitempty" gorm:"type:text"`
	WordFilePath *string   `json:"word_file_path,omitempty" gorm:"type:text"`
	WordFileName *string   `json:"word_file_name,omitempty" gorm:"size:191"`
	WordFileType *string   `json:"word_file_type,omitempty" gorm:"size:100"`
	Status       *string   `json:"status,omitempty" gorm:"size:20"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type GeneratedStatementLetterAIDocument struct {
	ID           string    `json:"id" gorm:"primaryKey;size:25"`
	StudentID    string    `json:"student_id" gorm:"size:25;not null;uniqueIndex"`
	Student      *User     `json:"student,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:StudentID;references:ID"`
	FileURL      string    `json:"file_url" gorm:"type:text;not null"`
	FilePath     *string   `json:"file_path,omitempty" gorm:"type:text"`
	FileName     *string   `json:"file_name,omitempty" gorm:"size:191"`
	FileType     *string   `json:"file_type,omitempty" gorm:"size:100"`
	WordFileURL  *string   `json:"word_file_url,omitempty" gorm:"type:text"`
	WordFilePath *string   `json:"word_file_path,omitempty" gorm:"type:text"`
	WordFileName *string   `json:"word_file_name,omitempty" gorm:"size:191"`
	WordFileType *string   `json:"word_file_type,omitempty" gorm:"size:100"`
	Status       *string   `json:"status,omitempty" gorm:"size:20"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AnswerApproval struct {
	ID         string          `json:"id" gorm:"primaryKey;size:25"`
	AnswerID   string          `json:"answer_id" gorm:"size:25;not null;index"`
	Answer     *AnswerQuestion `json:"answer,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	StudentID  string          `json:"student_id" gorm:"size:25;not null;index"`
	ReviewerID string          `json:"reviewer_id" gorm:"size:25;not null;index"`
	Status     string          `json:"status" gorm:"size:20;not null;default:'pending'"`
	Note       *string         `json:"note,omitempty" gorm:"type:text"`
	ReviewedAt *time.Time      `json:"reviewed_at,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type AnswerDocumentApproval struct {
	ID               string          `json:"id" gorm:"primaryKey;size:25"`
	AnswerDocumentID string          `json:"answer_document_id" gorm:"size:25;not null;index"`
	AnswerDocument   *AnswerDocument `json:"answer_document,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	StudentID        string          `json:"student_id" gorm:"size:25;not null;index"`
	ReviewerID       string          `json:"reviewer_id" gorm:"size:25;not null;index"`
	Status           string          `json:"status" gorm:"size:20;not null;default:'pending'"`
	Note             *string         `json:"note,omitempty" gorm:"type:text"`
	ReviewedAt       *time.Time      `json:"reviewed_at,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

type DocumentTranslation struct {
	ID               string               `json:"id" gorm:"primaryKey;size:25"`
	StudentID        string               `json:"student_id" gorm:"size:25;not null;index"`
	DocumentID       string               `json:"document_id" gorm:"size:25;not null;index"`
	Document         *DocumentsManagement `json:"document,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	AnswerDocumentID *string              `json:"answer_document_id,omitempty" gorm:"size:25;index"`
	AnswerDocument   *AnswerDocument      `json:"answer_document,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UploaderID       string               `json:"uploader_id" gorm:"size:25;not null;index"`
	Uploader         *User                `json:"uploader,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	FileURL          string               `json:"file_url" gorm:"type:text;not null"`
	FilePath         *string              `json:"file_path,omitempty" gorm:"type:text"`
	FileName         *string              `json:"file_name,omitempty" gorm:"size:191"`
	FileType         *string              `json:"file_type,omitempty" gorm:"size:50"`
	PageCount        int                  `json:"page_count" gorm:"not null;default:0"`
	Status           *string              `json:"status,omitempty" gorm:"size:20"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
}

type TicketMessage struct {
	ID             string            `json:"id" gorm:"primaryKey;size:25"`
	Name           string            `json:"name" gorm:"size:120;not null"`
	UserID         string            `json:"user_id" gorm:"size:25;not null;index"`
	ConversationID string            `json:"conversation_id" gorm:"size:25;not null;index"`
	User           *User             `json:"user,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:UserID;references:ID"`
	Conversation   *ChatConversation `json:"conversation,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ConversationID;references:ID"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}
