package schema

import (
	"errors"
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func tableName(tx *gorm.DB, model interface{}) (string, error) {
	stmt := &gorm.Statement{DB: tx}
	if err := stmt.Parse(model); err != nil {
		return "", err
	}
	return stmt.Schema.Table, nil
}

// Migrate applies versioned migrations in order.
// Add a new migration entry for every schema change.
func Migrate(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "20260304090000_init",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(AllModels()...)
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260304100000_question_base_country",
			Migrate: func(tx *gorm.DB) error {
				migrator := tx.Migrator()
				if !migrator.HasTable(&QuestionBase{}) {
					return nil
				}

				if !migrator.HasColumn(&QuestionBase{}, "TypeCountry") {
					if err := migrator.AddColumn(&QuestionBase{}, "TypeCountry"); err != nil {
						return err
					}
				}
				if !migrator.HasColumn(&QuestionBase{}, "CountryID") {
					if err := migrator.AddColumn(&QuestionBase{}, "CountryID"); err != nil {
						return err
					}
				}
				if !migrator.HasIndex(&QuestionBase{}, "CountryID") {
					if err := migrator.CreateIndex(&QuestionBase{}, "CountryID"); err != nil {
						return err
					}
				}
				if !migrator.HasConstraint(&QuestionBase{}, "Country") {
					if err := migrator.CreateConstraint(&QuestionBase{}, "Country"); err != nil {
						return err
					}
				}

				hasOldType := migrator.HasColumn(&QuestionBase{}, "type")
				hasNewType := migrator.HasColumn(&QuestionBase{}, "type_country")
				if hasOldType && hasNewType {
					if err := tx.Exec(`
						UPDATE question_bases
						SET type_country = type
						WHERE (type_country IS NULL OR type_country = '')
						  AND type IS NOT NULL
						  AND type <> ''
					`).Error; err != nil {
						return err
					}
					if err := migrator.DropColumn(&QuestionBase{}, "type"); err != nil {
						return err
					}
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260305110000_answer_submissions",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&AnswerSubmission{}, &AnswerQuestion{}, &AnswerSelectedOption{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260305130000_answer_documents",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&AnswerDocument{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260306110000_answer_approvals",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&AnswerApproval{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260306120000_answer_document_approvals",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&AnswerDocumentApproval{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260310120000_document_translations",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&DocumentTranslation{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260311100000_document_translation_page_count",
			Migrate: func(tx *gorm.DB) error {
				migrator := tx.Migrator()
				if !migrator.HasTable(&DocumentTranslation{}) {
					return nil
				}
				if !migrator.HasColumn(&DocumentTranslation{}, "PageCount") {
					if err := migrator.AddColumn(&DocumentTranslation{}, "PageCount"); err != nil {
						return err
					}
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260323170000_generated_cv_ai_documents",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&GeneratedCVAIDocument{}, &GeneratedStatementLetterAIDocument{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260324100000_generated_ai_document_relations",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&User{}, &GeneratedCVAIDocument{}, &GeneratedStatementLetterAIDocument{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260324143000_generated_statement_letter_word_backup",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&GeneratedStatementLetterAIDocument{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260324190000_generated_cv_word_backup",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&GeneratedCVAIDocument{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260330105000_statement_letter_ai_approvals",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&User{}, &GeneratedStatementLetterAIDocument{}, &StatementLetterAIApproval{}, &StatementLetterAIApprovalLog{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260330110000_generated_sponsor_letter_ai_documents",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&User{}, &GeneratedSponsorLetterAIDocument{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260330130000_sponsor_letter_ai_approvals",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&User{}, &GeneratedSponsorLetterAIDocument{}, &SponsorLetterAIApproval{}, &SponsorLetterAIApprovalLog{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260323183000_ticket_messages",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&TicketMessage{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260323191000_ticket_messages_remove_ticket_id",
			Migrate: func(tx *gorm.DB) error {
				migrator := tx.Migrator()
				if !migrator.HasTable(&TicketMessage{}) {
					return nil
				}
				if migrator.HasColumn(&TicketMessage{}, "TicketID") {
					if err := migrator.DropColumn(&TicketMessage{}, "TicketID"); err != nil {
						return err
					}
				}
				return tx.AutoMigrate(&TicketMessage{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260417103000_ticket_messages_conversation_nullable",
			Migrate: func(tx *gorm.DB) error {
				migrator := tx.Migrator()
				if !migrator.HasTable(&TicketMessage{}) {
					return nil
				}

				if err := tx.Exec(`
					UPDATE ticket_messages
					SET conversation_id = NULL
					WHERE conversation_id = ''
				`).Error; err != nil {
					return err
				}

				if err := tx.Exec(`
					ALTER TABLE ticket_messages
					MODIFY conversation_id VARCHAR(25) NULL
				`).Error; err != nil {
					return err
				}

				return tx.AutoMigrate(&TicketMessage{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260418110000_ticket_messages_add_status",
			Migrate: func(tx *gorm.DB) error {
				migrator := tx.Migrator()
				if !migrator.HasTable(&TicketMessage{}) {
					return nil
				}

				if !migrator.HasColumn(&TicketMessage{}, "Status") {
					if err := migrator.AddColumn(&TicketMessage{}, "Status"); err != nil {
						return err
					}
				}

				if err := tx.Exec(`
					UPDATE ticket_messages
					SET status = 'open'
					WHERE status IS NULL OR TRIM(status) = ''
				`).Error; err != nil {
					return err
				}

				if err := tx.Exec(`
					ALTER TABLE ticket_messages
					MODIFY status VARCHAR(20) NOT NULL DEFAULT 'open'
				`).Error; err != nil {
					return err
				}

				return tx.AutoMigrate(&TicketMessage{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260419120000_users_add_current_step_id",
			Migrate: func(tx *gorm.DB) error {
				migrator := tx.Migrator()
				if !migrator.HasTable(&User{}) {
					return nil
				}

				if !migrator.HasColumn(&User{}, "CurrentStepID") {
					if err := migrator.AddColumn(&User{}, "CurrentStepID"); err != nil {
						return err
					}
				}

				return tx.AutoMigrate(&User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260420101000_users_add_visa_status",
			Migrate: func(tx *gorm.DB) error {
				migrator := tx.Migrator()
				if !migrator.HasTable(&User{}) {
					return nil
				}

				if !migrator.HasColumn(&User{}, "VisaStatus") {
					if err := migrator.AddColumn(&User{}, "VisaStatus"); err != nil {
						return err
					}
				}

				return tx.AutoMigrate(&User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260420113000_generated_letters_add_source",
			Migrate: func(tx *gorm.DB) error {
				migrator := tx.Migrator()
				statementTable, err := tableName(tx, &GeneratedStatementLetterAIDocument{})
				if err != nil {
					return err
				}
				sponsorTable, err := tableName(tx, &GeneratedSponsorLetterAIDocument{})
				if err != nil {
					return err
				}

				if migrator.HasTable(&GeneratedStatementLetterAIDocument{}) {
					if !migrator.HasColumn(&GeneratedStatementLetterAIDocument{}, "Source") {
						if err := migrator.AddColumn(&GeneratedStatementLetterAIDocument{}, "Source"); err != nil {
							return err
						}
					}
					if err := tx.Exec(fmt.Sprintf(`
						UPDATE %s
						SET source = 'AI'
						WHERE source IS NULL OR TRIM(source) = ''
					`, statementTable)).Error; err != nil {
						return err
					}
				}

				if migrator.HasTable(&GeneratedSponsorLetterAIDocument{}) {
					if !migrator.HasColumn(&GeneratedSponsorLetterAIDocument{}, "Source") {
						if err := migrator.AddColumn(&GeneratedSponsorLetterAIDocument{}, "Source"); err != nil {
							return err
						}
					}
					if err := tx.Exec(fmt.Sprintf(`
						UPDATE %s
						SET source = 'AI'
						WHERE source IS NULL OR TRIM(source) = ''
					`, sponsorTable)).Error; err != nil {
						return err
					}
				}

				if err := tx.AutoMigrate(&GeneratedStatementLetterAIDocument{}, &GeneratedSponsorLetterAIDocument{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260420120000_users_add_visa_granted_at",
			Migrate: func(tx *gorm.DB) error {
				migrator := tx.Migrator()
				if !migrator.HasTable(&User{}) {
					return nil
				}

				if !migrator.HasColumn(&User{}, "VisaGrantedAt") {
					if err := migrator.AddColumn(&User{}, "VisaGrantedAt"); err != nil {
						return err
					}
				}

				userTable, err := tableName(tx, &User{})
				if err != nil {
					return err
				}

				if err := tx.Exec(fmt.Sprintf(`
					UPDATE %s
					SET visa_granted_at = updated_at
					WHERE visa_granted_at IS NULL
					  AND UPPER(TRIM(visa_status)) IN ('GRANT', 'GRANTED')
				`, userTable)).Error; err != nil {
					return err
				}

				return tx.AutoMigrate(&User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260420123000_users_add_student_status_audit",
			Migrate: func(tx *gorm.DB) error {
				migrator := tx.Migrator()
				if !migrator.HasTable(&User{}) {
					return nil
				}

				if !migrator.HasColumn(&User{}, "StudentStatusUpdatedByID") {
					if err := migrator.AddColumn(&User{}, "StudentStatusUpdatedByID"); err != nil {
						return err
					}
				}
				if !migrator.HasColumn(&User{}, "StudentStatusUpdatedAt") {
					if err := migrator.AddColumn(&User{}, "StudentStatusUpdatedAt"); err != nil {
						return err
					}
				}

				userTable, err := tableName(tx, &User{})
				if err != nil {
					return err
				}

				if err := tx.Exec(fmt.Sprintf(`
					UPDATE %s
					SET student_status_updated_at = updated_at
					WHERE student_status_updated_at IS NULL
					  AND UPPER(TRIM(student_status)) IN ('POSTPONE', 'CANCEL')
				`, userTable)).Error; err != nil {
					return err
				}

				return tx.AutoMigrate(&User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260421150000_information_country_managements",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&CountryManagement{}, &InformationCountryManagement{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260421153000_information_country_managements_add_priority",
			Migrate: func(tx *gorm.DB) error {
				migrator := tx.Migrator()
				if !migrator.HasTable(&InformationCountryManagement{}) {
					return nil
				}

				if !migrator.HasColumn(&InformationCountryManagement{}, "Priority") {
					if err := migrator.AddColumn(&InformationCountryManagement{}, "Priority"); err != nil {
						return err
					}
				}

				table, err := tableName(tx, &InformationCountryManagement{})
				if err != nil {
					return err
				}

				if err := tx.Exec(fmt.Sprintf(`
					UPDATE %s
					SET priority = 'normal'
					WHERE priority IS NULL OR TRIM(priority) = ''
				`, table)).Error; err != nil {
					return err
				}

				return tx.AutoMigrate(&InformationCountryManagement{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
	})

	if err := m.Migrate(); err != nil {
		return fmt.Errorf("migrate failed: %w", err)
	}

	return nil
}
