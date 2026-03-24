package schema

// AllModels returns all GORM models for AutoMigrate.
func AllModels() []interface{} {
	return []interface{}{
		&User{},
		&NoteStudent{},
		&DocumentsManagement{},
		&StepsManagement{},
		&ChildStepsManagement{},
		&CountryManagement{},
		&StageManagement{},
		&CountryStepsManagement{},
		&ChatConversation{},
		&ChatConversationMember{},
		&ChatMessage{},
		&ChatMessageAttachment{},
		&ChatMessageStatus{},
		&ChatMessageMention{},
		&Question{},
		&QuestionOption{},
		&AnswerQuestion{},
		&AnswerSubmission{},
		&AnswerSelectedOption{},
		&AnswerDocument{},
		&GeneratedCVAIDocument{},
		&GeneratedStatementLetterAIDocument{},
		&AnswerApproval{},
		&AnswerDocumentApproval{},
		&DocumentTranslation{},
		&QuestionBase{},
		&TicketMessage{},
	}
}
