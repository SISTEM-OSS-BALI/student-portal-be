package schema

import (
	"errors"
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

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
	})

	if err := m.Migrate(); err != nil {
		return fmt.Errorf("migrate failed: %w", err)
	}

	return nil
}
