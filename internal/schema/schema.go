package schema

import "time"
import "gorm.io/gorm"

type UserRole string
type TranslationNeeded string
type AutoRenamePattern string
type StatusStudent string
type QuestionType string
type StatementLetterDocumentStatus string
type StatementLetterApprovalStatus string
type SponsorLetterDocumentStatus string
type SponsorLetterApprovalStatus string
type GeneratedDocumentSource string

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

const (
	StatementLetterDocumentStatusDraft             StatementLetterDocumentStatus = "DRAFT"
	StatementLetterDocumentStatusSubmittedDirector StatementLetterDocumentStatus = "SUBMITTED_TO_DIRECTOR"
	StatementLetterDocumentStatusRevisionRequested StatementLetterDocumentStatus = "REVISION_REQUESTED"
	StatementLetterDocumentStatusApproved          StatementLetterDocumentStatus = "APPROVED"
)

const (
	StatementLetterApprovalStatusPending           StatementLetterApprovalStatus = "PENDING"
	StatementLetterApprovalStatusApproved          StatementLetterApprovalStatus = "APPROVED"
	StatementLetterApprovalStatusRevisionRequested StatementLetterApprovalStatus = "REVISION_REQUESTED"
	StatementLetterApprovalStatusRejected          StatementLetterApprovalStatus = "REJECTED"
	StatementLetterApprovalStatusCanceled          StatementLetterApprovalStatus = "CANCELED"
)

const (
	SponsorLetterDocumentStatusDraft             SponsorLetterDocumentStatus = "DRAFT"
	SponsorLetterDocumentStatusSubmittedDirector SponsorLetterDocumentStatus = "SUBMITTED_TO_DIRECTOR"
	SponsorLetterDocumentStatusRevisionRequested SponsorLetterDocumentStatus = "REVISION_REQUESTED"
	SponsorLetterDocumentStatusApproved          SponsorLetterDocumentStatus = "APPROVED"
)

const (
	SponsorLetterApprovalStatusPending           SponsorLetterApprovalStatus = "PENDING"
	SponsorLetterApprovalStatusApproved          SponsorLetterApprovalStatus = "APPROVED"
	SponsorLetterApprovalStatusRevisionRequested SponsorLetterApprovalStatus = "REVISION_REQUESTED"
	SponsorLetterApprovalStatusRejected          SponsorLetterApprovalStatus = "REJECTED"
	SponsorLetterApprovalStatusCanceled          SponsorLetterApprovalStatus = "CANCELED"
)

const (
	GeneratedDocumentSourceAI     GeneratedDocumentSource = "AI"
	GeneratedDocumentSourceManual GeneratedDocumentSource = "MANUAL"
)

type User struct {
	ID                           string                              `json:"id" gorm:"primaryKey;size:25"`
	Name                         string                              `json:"name" gorm:"size:120;not null"`
	Email                        string                              `json:"email" gorm:"size:191;uniqueIndex;not null"`
	Password                     string                              `json:"-" gorm:"size:191;not null"`
	Role                         UserRole                            `json:"role" gorm:"size:50;not null"`
	NoPhone                      *string                             `json:"no_phone,omitempty" gorm:"size:20"`
	StageID                      *string                             `json:"stage_id,omitempty" gorm:"size:25;index"`
	CurrentStepID                *string                             `json:"current_step_id,omitempty" gorm:"size:25;index"`
	VisaStatus                   *string                             `json:"visa_status,omitempty" gorm:"size:25;index"`
	VisaGrantedAt                *time.Time                          `json:"visa_granted_at,omitempty" gorm:"index"`
	StudentStatus                StatusStudent                       `json:"student_status" gorm:"size:20;not null;default:'ON GOING'"`
	StudentStatusUpdatedByID     *string                             `json:"student_status_updated_by_id,omitempty" gorm:"size:25;index"`
	StudentStatusUpdatedBy       *User                               `json:"student_status_updated_by,omitempty" gorm:"foreignKey:StudentStatusUpdatedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StudentStatusUpdatedAt       *time.Time                          `json:"student_status_updated_at,omitempty" gorm:"index"`
	NameConsultant               *string                             `json:"name_consultant,omitempty" gorm:"size:120;index"`
	Stage                        *StageManagement                    `json:"stage,omitempty" gorm:"foreignKey:StageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	NameCampus                   *string                             `json:"name_campus,omitempty" gorm:"size:100"`
	Degree                       *string                             `json:"degree,omitempty" gorm:"size:100"`
	NameDegree                   *string                             `json:"name_degree,omitempty" gorm:"size:100"`
	DocumentConsentSignatureURL  *string                             `json:"document_consent_signature_url,omitempty" gorm:"size:500"`
	DocumentConsentProofPhotoURL *string                             `json:"document_consent_proof_photo_url,omitempty" gorm:"size:500"`
	DocumentConsentSignedAt      *time.Time                          `json:"document_consent_signed_at,omitempty" gorm:"index"`
	DocumentConsentSigned        bool                                `json:"document_consent_signed" gorm:"not null;default:false"`
	VisaType                     *string                             `json:"visa_type,omitempty" gorm:"size:25;index"`
	VisaTypeDetail               *VisaTypeManagement                 `json:"visa_type_detail,omitempty" gorm:"foreignKey:VisaType;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Source                       *string                             `json:"source,omitempty" gorm:"size:100;index"`
	SourceCategory               *string                             `json:"source_category,omitempty" gorm:"size:50;index"`
	TranslationQuota             int                                 `json:"translation_quota" gorm:"not null;default:0"`
	HasInitialTranslations       bool                                `json:"has_initial_translations" gorm:"not null;default:false"`
	NotesStudent                 []NoteStudent                       `json:"notes,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GeneratedCVAI                *GeneratedCVAIDocument              `json:"generated_cv_ai,omitempty" gorm:"foreignKey:StudentID;references:ID"`
	GeneratedStatementLetterAI   *GeneratedStatementLetterAIDocument `json:"generated_statement_letter_ai,omitempty" gorm:"foreignKey:StudentID;references:ID"`
	GeneratedSponsorLetterAI     *GeneratedSponsorLetterAIDocument   `json:"generated_sponsor_letter_ai,omitempty" gorm:"foreignKey:StudentID;references:ID"`
	CreatedAt                    time.Time                           `json:"created_at"`
	UpdatedAt                    time.Time                           `json:"updated_at"`
}

type NoteStudent struct {
	ID        string    `json:"id" gorm:"primaryKey;size:25"`
	UserID    string    `json:"user_id" gorm:"size:25;not null;index"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PasswordResetOTP struct {
	ID        string     `json:"id" gorm:"primaryKey;size:25"`
	UserID    string     `json:"user_id" gorm:"size:25;not null;index"`
	Email     string     `json:"email" gorm:"size:191;not null;index"`
	Code      string     `json:"code" gorm:"size:191;not null"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null;index"`
	UsedAt    *time.Time `json:"used_at,omitempty" gorm:"index"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type DocumentsManagement struct {
	ID                string            `json:"id" gorm:"primaryKey;size:25"`
	Label             string            `json:"label" gorm:"size:120;not null"`
	InternalCode      string            `json:"internal_code" gorm:"size:50;not null;uniqueIndex"`
	FileType          string            `json:"file_type" gorm:"size:50;not null"`
	Category          string            `json:"category" gorm:"size:100;not null"`
	ExampleURL        *string           `json:"example_url,omitempty" gorm:"size:500"`
	TranslationNeeded TranslationNeeded `json:"translation_needed" gorm:"size:10;not null"`
	Required          bool              `json:"required" gorm:"not null"`
	AutoRenamePattern AutoRenamePattern `json:"auto_rename_pattern" gorm:"size:191"`
	Notes             string            `json:"notes" gorm:"type:text"`
	Stages            []StageManagement `json:"stages,omitempty" gorm:"foreignKey:DocumentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	DeletedAt         gorm.DeletedAt    `json:"deleted_at,omitempty" gorm:"index"`
}

type StepsManagement struct {
	ID           string                   `json:"id" gorm:"primaryKey;size:25"`
	Label        string                   `json:"label" gorm:"size:120;not null"`
	Children     []ChildStepsManagement   `json:"children,omitempty" gorm:"many2many:steps_children;"`
	CountrySteps []CountryStepsManagement `json:"country_steps,omitempty" gorm:"foreignKey:StepID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
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
	ID           string                         `json:"id" gorm:"primaryKey;size:25"`
	NameCountry  string                         `json:"name" gorm:"column:name_country;size:120;not null"`
	Stages       []StageManagement              `json:"stages,omitempty" gorm:"foreignKey:CountryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CountrySteps []CountryStepsManagement       `json:"country_steps,omitempty" gorm:"foreignKey:CountryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Informations []InformationCountryManagement `json:"informations,omitempty" gorm:"foreignKey:CountryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	VisaTypes    []VisaTypeManagement           `json:"visa_types,omitempty" gorm:"foreignKey:CountryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt    time.Time                      `json:"created_at"`
	UpdatedAt    time.Time                      `json:"updated_at"`
}

type VisaTypeManagement struct {
	ID        string             `json:"id" gorm:"primaryKey;size:25"`
	Name      string             `json:"name" gorm:"size:120;not null"`
	CountryID string             `json:"country_id" gorm:"size:25;not null;index"`
	Country   *CountryManagement `json:"country,omitempty" gorm:"foreignKey:CountryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt gorm.DeletedAt     `json:"deleted_at,omitempty" gorm:"index"`
}

type StageManagement struct {
	ID         string               `json:"id" gorm:"primaryKey;size:25"`
	CountryID  string               `json:"country_id" gorm:"size:25;not null;index"`
	DocumentID string               `json:"document_id" gorm:"size:25;not null;index"`
	Country    *CountryManagement   `json:"country,omitempty" gorm:"foreignKey:CountryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Document   *DocumentsManagement `json:"document,omitempty" gorm:"foreignKey:DocumentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Users      []User               `json:"users,omitempty" gorm:"foreignKey:StageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
}

type CountryStepsManagement struct {
	ID        string             `json:"id" gorm:"primaryKey;size:25"`
	CountryID string             `json:"country_id" gorm:"size:25;not null;index"`
	StepID    string             `json:"step_id" gorm:"size:25;not null;index"`
	Country   *CountryManagement `json:"country,omitempty" gorm:"foreignKey:CountryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Step      *StepsManagement   `json:"step,omitempty" gorm:"foreignKey:StepID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type QuestionBase struct {
	ID                       string             `json:"id" gorm:"primaryKey;size:25"`
	Name                     string             `json:"name" gorm:"size:120;not null"`
	Desc                     *string            `json:"desc,omitempty" gorm:"type:text"`
	TypeCountry              string             `json:"type" gorm:"column:type_country;size:50;not null"`
	CountryID                *string            `json:"country_id,omitempty" gorm:"size:25;index"`
	Country                  *CountryManagement `json:"country,omitempty" gorm:"foreignKey:CountryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Questions                []Question         `json:"questions,omitempty" gorm:"foreignKey:BaseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AllowMultipleSubmissions bool               `json:"allow_multiple_submissions" gorm:"not null;default:false"`
	Active                   bool               `json:"active" gorm:"not null;default:true"`
	Version                  int                `json:"version" gorm:"not null;default:1"`
	CreatedAt                time.Time          `json:"created_at"`
	UpdatedAt                time.Time          `json:"updated_at"`
	DeletedAt                gorm.DeletedAt     `json:"deleted_at,omitempty" gorm:"index"`
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
	Base        *QuestionBase    `json:"base,omitempty" gorm:"foreignKey:BaseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Options     []QuestionOption `json:"options,omitempty" gorm:"foreignKey:QuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Answers     []AnswerQuestion `json:"answers,omitempty" gorm:"foreignKey:QuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Active      bool             `json:"active" gorm:"not null;default:true"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `json:"deleted_at,omitempty" gorm:"index"`
}

type QuestionOption struct {
	ID         string                 `json:"id" gorm:"primaryKey;size:25"`
	QuestionID string                 `json:"question_id" gorm:"size:25;not null;index;uniqueIndex:uniq_question_option,priority:1;index:idx_question_option_order,priority:1"`
	Label      string                 `json:"label" gorm:"type:text;not null"`
	Value      string                 `json:"value" gorm:"type:varchar(255);not null;uniqueIndex:uniq_question_option,priority:2"`
	Order      int                    `json:"order" gorm:"not null;default:0;index:idx_question_option_order,priority:2"`
	Question   *Question              `json:"question,omitempty" gorm:"foreignKey:QuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SelectedBy []AnswerSelectedOption `json:"selected_by,omitempty" gorm:"foreignKey:OptionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Active     bool                   `json:"active" gorm:"not null;default:true"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
	DeletedAt  gorm.DeletedAt         `json:"deleted_at,omitempty" gorm:"index"`
}

type ChatConversation struct {
	ID          string                   `json:"id" gorm:"primaryKey;size:25"`
	Type        string                   `json:"type" gorm:"size:16;not null"`
	Title       *string                  `json:"title" gorm:"size:120"`
	CreatedByID string                   `json:"created_by_id" gorm:"size:25;not null;index"`
	Members     []ChatConversationMember `json:"members,omitempty" gorm:"foreignKey:ConversationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Messages    []ChatMessage            `json:"messages,omitempty" gorm:"foreignKey:ConversationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
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
	ID             string                  `json:"id" gorm:"primaryKey;size:25"`
	ConversationID string                  `json:"conversation_id" gorm:"size:25;not null;index"`
	SenderID       string                  `json:"sender_id" gorm:"size:25;not null;index"`
	Type           string                  `json:"type" gorm:"size:16;not null"`
	Text           *string                 `json:"text,omitempty" gorm:"type:text"`
	ReplyToID      *string                 `json:"reply_to_id,omitempty" gorm:"size:25;index"`
	ContextUserID  *string                 `json:"context_user_id,omitempty" gorm:"size:25;index"`
	ContextType    string                  `json:"context_type,omitempty" gorm:"size:32"`
	Sender         *User                   `json:"sender,omitempty" gorm:"foreignKey:SenderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Attachments    []ChatMessageAttachment `json:"attachments,omitempty" gorm:"foreignKey:MessageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Statuses       []ChatMessageStatus     `json:"statuses,omitempty" gorm:"foreignKey:MessageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Mentions       []ChatMessageMention    `json:"mentions,omitempty" gorm:"foreignKey:MessageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt      time.Time               `json:"created_at"`
	EditedAt       *time.Time              `json:"edited_at,omitempty"`
	DeletedAt      *time.Time              `json:"deleted_at,omitempty" gorm:"index"`
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

type Conversation struct {
	ID          string               `gorm:"primaryKey;size:26"`
	Type        string               `gorm:"size:16;not null"`
	Title       *string              `gorm:"size:120"`
	CreatedByID string               `gorm:"size:26;not null;index"`
	Members     []ConversationMember `gorm:"foreignKey:ConversationID;references:ID"`
	Messages    []Message            `gorm:"foreignKey:ConversationID;references:ID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ConversationMember struct {
	ID             string    `gorm:"primaryKey;size:26"`
	ConversationID string    `gorm:"size:26;not null;index"`
	UserID         string    `gorm:"size:26;not null;index"`
	Role           string    `gorm:"size:16;not null"`
	JoinedAt       time.Time `gorm:"not null"`
	LeftAt         *time.Time
	MutedUntil     *time.Time
	LastReadAt     *time.Time
}

type Message struct {
	ID             string              `gorm:"primaryKey;size:26"`
	ConversationID string              `gorm:"size:26;not null;index"`
	SenderID       string              `gorm:"size:26;not null;index"`
	Type           string              `gorm:"size:16;not null"`
	Text           *string             `gorm:"type:text"`
	ReplyToID      *string             `gorm:"size:26;index"`
	Attachments    []MessageAttachment `gorm:"foreignKey:MessageID;references:ID"`
	Statuses       []MessageStatus     `gorm:"foreignKey:MessageID;references:ID"`
	CreatedAt      time.Time           `gorm:"index"`
	EditedAt       *time.Time
	DeletedAt      *time.Time `gorm:"index"`
}

type MessageAttachment struct {
	ID        string `gorm:"primaryKey;size:26"`
	MessageID string `gorm:"size:26;not null;index"`
	FileURL   string `gorm:"size:512;not null"`
	FileName  string `gorm:"size:255"`
	FileType  string `gorm:"size:64"`
	FileSize  int64
	CreatedAt time.Time
}

type MessageStatus struct {
	ID        string    `gorm:"primaryKey;size:26"`
	MessageID string    `gorm:"size:26;not null;index"`
	UserID    string    `gorm:"size:26;not null;index"`
	Status    string    `gorm:"size:16;not null"`
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
	SubmissionID *string   `json:"submission_id,omitempty" gorm:"size:25;index"`
	QuestionID   string    `json:"question_id" gorm:"size:25;not null;index"`
	AnswerText   *string   `json:"answer_text,omitempty" gorm:"type:text"`
	StudentID    *string   `json:"student_id,omitempty" gorm:"size:25;index"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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
	Document     *DocumentsManagement `json:"document,omitempty" gorm:"foreignKey:DocumentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	FileURL      string               `json:"file_url" gorm:"type:text;not null"`
	FilePath     *string              `json:"file_path,omitempty" gorm:"type:text"`
	FileName     *string              `json:"file_name,omitempty" gorm:"size:191"`
	FileType     *string              `json:"file_type,omitempty" gorm:"size:50"`
	Status       *string              `json:"status,omitempty" gorm:"size:20"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}

type GeneratedCVAIDocument struct {
	ID        string `json:"id" gorm:"primaryKey;size:25"`
	StudentID string `json:"student_id" gorm:"size:25;not null;uniqueIndex"`
	Student   *User  `json:"student,omitempty" gorm:"foreignKey:StudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

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
	ID        string `json:"id" gorm:"primaryKey;size:25"`
	StudentID string `json:"student_id" gorm:"size:25;not null;uniqueIndex"`
	Student   *User  `json:"student,omitempty" gorm:"foreignKey:StudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	FileURL  string  `json:"file_url" gorm:"type:text;not null"`
	FilePath *string `json:"file_path,omitempty" gorm:"type:text"`
	FileName *string `json:"file_name,omitempty" gorm:"size:191"`
	FileType *string `json:"file_type,omitempty" gorm:"size:100"`

	WordFileURL  *string `json:"word_file_url,omitempty" gorm:"type:text"`
	WordFilePath *string `json:"word_file_path,omitempty" gorm:"type:text"`
	WordFileName *string `json:"word_file_name,omitempty" gorm:"size:191"`
	WordFileType *string `json:"word_file_type,omitempty" gorm:"size:100"`

	Status StatementLetterDocumentStatus `json:"status" gorm:"size:40;not null;default:'DRAFT'"`
	Source GeneratedDocumentSource       `json:"source" gorm:"size:20;not null;default:'AI'"`

	SubmittedToDirectorAt *time.Time `json:"submitted_to_director_at,omitempty"`
	ApprovedAt            *time.Time `json:"approved_at,omitempty"`
	RevisionRequestedAt   *time.Time `json:"revision_requested_at,omitempty"`

	CurrentApprovalID *string `json:"current_approval_id,omitempty" gorm:"size:25;index"`

	CurrentApproval *StatementLetterAIApproval  `json:"current_approval,omitempty" gorm:"foreignKey:CurrentApprovalID;references:ID"`
	Approvals       []StatementLetterAIApproval `json:"approvals,omitempty" gorm:"foreignKey:DocumentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt       time.Time                   `json:"created_at"`
	UpdatedAt       time.Time                   `json:"updated_at"`
}

type StatementLetterAIApproval struct {
	ID         string                              `json:"id" gorm:"primaryKey;size:25"`
	DocumentID string                              `json:"document_id" gorm:"size:25;not null;index"`
	Document   *GeneratedStatementLetterAIDocument `json:"document,omitempty" gorm:"foreignKey:DocumentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ReviewerID string `json:"reviewer_id" gorm:"size:25;not null;index"`
	Reviewer   *User  `json:"reviewer,omitempty" gorm:"foreignKey:ReviewerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Status     StatementLetterApprovalStatus `json:"status" gorm:"size:30;not null;default:'PENDING'"`
	Note       *string                       `json:"note,omitempty" gorm:"type:text"`
	ReviewedAt *time.Time                    `json:"reviewed_at,omitempty"`

	Logs      []StatementLetterAIApprovalLog `json:"logs,omitempty" gorm:"foreignKey:ApprovalID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time                      `json:"created_at"`
	UpdatedAt time.Time                      `json:"updated_at"`
}

type StatementLetterAIApprovalLog struct {
	ID         string                     `json:"id" gorm:"primaryKey;size:25"`
	ApprovalID string                     `json:"approval_id" gorm:"size:25;not null;index"`
	Approval   *StatementLetterAIApproval `json:"approval,omitempty" gorm:"foreignKey:ApprovalID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ActorID string `json:"actor_id" gorm:"size:25;not null;index"`
	Actor   *User  `json:"actor,omitempty" gorm:"foreignKey:ActorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	FromStatus *StatementLetterApprovalStatus `json:"from_status,omitempty" gorm:"size:30"`
	ToStatus   StatementLetterApprovalStatus  `json:"to_status" gorm:"size:30;not null"`
	Note       *string                        `json:"note,omitempty" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at"`
}

type GeneratedSponsorLetterAIDocument struct {
	ID        string `json:"id" gorm:"primaryKey;size:25"`
	StudentID string `json:"student_id" gorm:"size:25;not null;uniqueIndex"`
	Student   *User  `json:"student,omitempty" gorm:"foreignKey:StudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	FileURL  string  `json:"file_url" gorm:"type:text;not null"`
	FilePath *string `json:"file_path,omitempty" gorm:"type:text"`
	FileName *string `json:"file_name,omitempty" gorm:"size:191"`
	FileType *string `json:"file_type,omitempty" gorm:"size:100"`

	WordFileURL  *string `json:"word_file_url,omitempty" gorm:"type:text"`
	WordFilePath *string `json:"word_file_path,omitempty" gorm:"type:text"`
	WordFileName *string `json:"word_file_name,omitempty" gorm:"size:191"`
	WordFileType *string `json:"word_file_type,omitempty" gorm:"size:100"`

	Status SponsorLetterDocumentStatus `json:"status" gorm:"size:40;not null;default:'DRAFT'"`
	Source GeneratedDocumentSource     `json:"source" gorm:"size:20;not null;default:'AI'"`

	SubmittedToDirectorAt *time.Time `json:"submitted_to_director_at,omitempty"`
	ApprovedAt            *time.Time `json:"approved_at,omitempty"`
	RevisionRequestedAt   *time.Time `json:"revision_requested_at,omitempty"`

	CurrentApprovalID *string `json:"current_approval_id,omitempty" gorm:"size:25;index"`

	CurrentApproval *SponsorLetterAIApproval  `json:"current_approval,omitempty" gorm:"foreignKey:CurrentApprovalID;references:ID"`
	Approvals       []SponsorLetterAIApproval `json:"approvals,omitempty" gorm:"foreignKey:DocumentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

type SponsorLetterAIApproval struct {
	ID         string                            `json:"id" gorm:"primaryKey;size:25"`
	DocumentID string                            `json:"document_id" gorm:"size:25;not null;index"`
	Document   *GeneratedSponsorLetterAIDocument `json:"document,omitempty" gorm:"foreignKey:DocumentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ReviewerID string `json:"reviewer_id" gorm:"size:25;not null;index"`
	Reviewer   *User  `json:"reviewer,omitempty" gorm:"foreignKey:ReviewerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Status     SponsorLetterApprovalStatus `json:"status" gorm:"size:30;not null;default:'PENDING'"`
	Note       *string                     `json:"note,omitempty" gorm:"type:text"`
	ReviewedAt *time.Time                  `json:"reviewed_at,omitempty"`

	Logs      []SponsorLetterAIApprovalLog `json:"logs,omitempty" gorm:"foreignKey:ApprovalID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time                    `json:"created_at"`
	UpdatedAt time.Time                    `json:"updated_at"`
}

type SponsorLetterAIApprovalLog struct {
	ID         string                   `json:"id" gorm:"primaryKey;size:25"`
	ApprovalID string                   `json:"approval_id" gorm:"size:25;not null;index"`
	Approval   *SponsorLetterAIApproval `json:"approval,omitempty" gorm:"foreignKey:ApprovalID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ActorID string `json:"actor_id" gorm:"size:25;not null;index"`
	Actor   *User  `json:"actor,omitempty" gorm:"foreignKey:ActorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	FromStatus *SponsorLetterApprovalStatus `json:"from_status,omitempty" gorm:"size:30"`
	ToStatus   SponsorLetterApprovalStatus  `json:"to_status" gorm:"size:30;not null"`
	Note       *string                      `json:"note,omitempty" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at"`
}

type AnswerApproval struct {
	ID         string          `json:"id" gorm:"primaryKey;size:25"`
	AnswerID   string          `json:"answer_id" gorm:"size:25;not null;index"`
	Answer     *AnswerQuestion `json:"answer,omitempty" gorm:"foreignKey:AnswerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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
	AnswerDocument   *AnswerDocument `json:"answer_document,omitempty" gorm:"foreignKey:AnswerDocumentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	StudentID        string          `json:"student_id" gorm:"size:25;not null;index"`
	ReviewerID       string          `json:"reviewer_id" gorm:"size:25;not null;index"`
	Status           string          `json:"status" gorm:"size:20;not null;default:'pending'"`
	Note             *string         `json:"note,omitempty" gorm:"type:text"`
	ReviewedAt       *time.Time      `json:"reviewed_at,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

type DocumentTranslation struct {
	ID                    string               `json:"id" gorm:"primaryKey;size:25"`
	StudentID             string               `json:"student_id" gorm:"size:25;not null;index"`
	DocumentID            string               `json:"document_id" gorm:"size:25;not null;index"`
	Document              *DocumentsManagement `json:"document,omitempty" gorm:"foreignKey:DocumentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	AnswerDocumentID      *string              `json:"answer_document_id,omitempty" gorm:"size:25;index"`
	AnswerDocument        *AnswerDocument      `json:"answer_document,omitempty" gorm:"foreignKey:AnswerDocumentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UploaderID            string               `json:"uploader_id" gorm:"size:25;not null;index"`
	Uploader              *User                `json:"uploader,omitempty" gorm:"foreignKey:UploaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FileURL               string               `json:"file_url" gorm:"type:text;not null"`
	FilePath              *string              `json:"file_path,omitempty" gorm:"type:text"`
	FileName              *string              `json:"file_name,omitempty" gorm:"size:191"`
	FileType              *string              `json:"file_type,omitempty" gorm:"size:50"`
	PageCount             int                  `json:"page_count" gorm:"not null;default:0"`
	IsExistingTranslation bool                 `json:"is_existing_translation" gorm:"not null;default:false"`
	Status                *string              `json:"status,omitempty" gorm:"size:20"`
	CreatedAt             time.Time            `json:"created_at"`
	UpdatedAt             time.Time            `json:"updated_at"`
}

type TicketMessage struct {
	ID             string            `json:"id" gorm:"primaryKey;size:25"`
	Name           string            `json:"name" gorm:"size:120;not null"`
	UserID         string            `json:"user_id" gorm:"size:25;not null;index"`
	ConversationID *string           `json:"conversation_id,omitempty" gorm:"size:25;index"`
	Status         string            `json:"status" gorm:"size:20;not null;default:'open'"`
	User           *User             `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Conversation   *ChatConversation `json:"conversation,omitempty" gorm:"foreignKey:ConversationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type InformationCountryManagement struct {
	ID          string             `json:"id" gorm:"primaryKey;size:25"`
	Slug        string             `json:"slug" gorm:"size:25;not null;uniqueIndex"`
	Title       string             `json:"title" gorm:"size:120;not null"`
	Description *string            `json:"description,omitempty" gorm:"type:text"`
	Priority    string             `json:"priority" gorm:"size:30;not null;default:'normal'"`
	CountryID   string             `json:"country_id" gorm:"size:25;not null;index"`
	Country     *CountryManagement `json:"country,omitempty" gorm:"foreignKey:CountryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type Promo struct {
	ID          string    `json:"id" gorm:"primaryKey;size:25"`
	Code        string    `json:"code" gorm:"size:50;not null;uniqueIndex"`
	Description *string   `json:"description,omitempty" gorm:"type:text"`
	Discount    float64   `json:"discount" gorm:"not null;default:0"`
	ValidFrom   time.Time `json:"valid_from" gorm:"not null"`
	ValidTo     time.Time `json:"valid_to" gorm:"not null"`
	IsActive    bool      `json:"is_active" gorm:"not null;default:true;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (CountryManagement) TableName() string {
	return "country_managements"
}

func (VisaTypeManagement) TableName() string {
	return "visa_type_managements"
}

func (GeneratedCVAIDocument) TableName() string {
	return "generated_cv_ai_documents"
}

func (GeneratedStatementLetterAIDocument) TableName() string {
	return "generated_statement_letter_ai_documents"
}

func (StatementLetterAIApproval) TableName() string {
	return "statement_letter_ai_approvals"
}

func (StatementLetterAIApprovalLog) TableName() string {
	return "statement_letter_ai_approval_logs"
}

func (GeneratedSponsorLetterAIDocument) TableName() string {
	return "generated_sponsor_letter_ai_documents"
}

func (SponsorLetterAIApproval) TableName() string {
	return "sponsor_letter_ai_approvals"
}

func (SponsorLetterAIApprovalLog) TableName() string {
	return "sponsor_letter_ai_approval_logs"
}
