package schema

import (
	"errors"
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type foreignKeyRef struct {
	ConstraintName string `gorm:"column:CONSTRAINT_NAME"`
}

type indexRef struct {
	IndexName string `gorm:"column:INDEX_NAME"`
}

func tableName(tx *gorm.DB, model interface{}) (string, error) {
	stmt := &gorm.Statement{DB: tx}
	if err := stmt.Parse(model); err != nil {
		return "", err
	}
	return stmt.Schema.Table, nil
}

func replaceFK(tx *gorm.DB, table, column, refTable, refColumn, constraintName, onDeleteAction string) error {
	if tx.Dialector.Name() != "mysql" {
		// This migration targets MySQL only.
		return nil
	}

	var existing []foreignKeyRef
	if err := tx.Raw(`
		SELECT CONSTRAINT_NAME
		FROM information_schema.KEY_COLUMN_USAGE
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = ?
		  AND COLUMN_NAME = ?
		  AND REFERENCED_TABLE_NAME = ?
		  AND REFERENCED_COLUMN_NAME = ?
		  AND CONSTRAINT_NAME IS NOT NULL
	`, table, column, refTable, refColumn).Scan(&existing).Error; err != nil {
		return err
	}

	for _, fk := range existing {
		if fk.ConstraintName == "" {
			continue
		}
		if err := tx.Exec(fmt.Sprintf(
			"ALTER TABLE `%s` DROP FOREIGN KEY `%s`",
			table, fk.ConstraintName,
		)).Error; err != nil {
			return err
		}
	}

	// Ensure there is an index on the FK column (MySQL requires it).
	if err := tx.Exec(fmt.Sprintf(
		"ALTER TABLE `%s` ADD INDEX `idx_%s_%s` (`%s`)",
		table, table, column, column,
	)).Error; err != nil {
		// ignore if index already exists or can't be created
	}

	return tx.Exec(fmt.Sprintf(`
		ALTER TABLE %s
		ADD CONSTRAINT %s
		FOREIGN KEY (%s) REFERENCES %s (%s)
		ON DELETE %s
		ON UPDATE CASCADE
	`,
		tx.Statement.Quote(table),
		tx.Statement.Quote(constraintName),
		tx.Statement.Quote(column),
		tx.Statement.Quote(refTable),
		tx.Statement.Quote(refColumn),
		onDeleteAction,
	)).Error
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
		{
			ID: "20260429120000_country_managements_fix_name_column",
			Migrate: func(tx *gorm.DB) error {
				migrator := tx.Migrator()
				if !migrator.HasTable(&CountryManagement{}) {
					return nil
				}

				hasOldName := migrator.HasColumn(&CountryManagement{}, "name")
				hasNewName := migrator.HasColumn(&CountryManagement{}, "name_country")

				// Ensure the new column exists (current schema uses name_country).
				if !hasNewName {
					if err := migrator.AddColumn(&CountryManagement{}, "NameCountry"); err != nil {
						return err
					}
					hasNewName = migrator.HasColumn(&CountryManagement{}, "name_country")
				}

				if hasOldName && hasNewName {
					table, err := tableName(tx, &CountryManagement{})
					if err != nil {
						return err
					}

					if err := tx.Exec(fmt.Sprintf(`
						UPDATE %s
						SET name_country = name
						WHERE (name_country IS NULL OR TRIM(name_country) = '')
						  AND name IS NOT NULL
						  AND TRIM(name) <> ''
					`, table)).Error; err != nil {
						return err
					}

					if err := migrator.DropColumn(&CountryManagement{}, "name"); err != nil {
						return err
					}
				}

				return tx.AutoMigrate(&CountryManagement{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260504093000_country_relations_on_delete_cascade",
			Migrate: func(tx *gorm.DB) error {
				countryTable, err := tableName(tx, &CountryManagement{})
				if err != nil {
					return err
				}
				stageTable, err := tableName(tx, &StageManagement{})
				if err != nil {
					return err
				}
				countryStepsTable, err := tableName(tx, &CountryStepsManagement{})
				if err != nil {
					return err
				}
				infoTable, err := tableName(tx, &InformationCountryManagement{})
				if err != nil {
					return err
				}

				if err := replaceFK(tx, stageTable, "country_id", countryTable, "id", "fk_stage_managements_country", "CASCADE"); err != nil {
					return err
				}
				if err := replaceFK(tx, countryStepsTable, "country_id", countryTable, "id", "fk_country_steps_managements_country", "CASCADE"); err != nil {
					return err
				}
				if err := replaceFK(tx, infoTable, "country_id", countryTable, "id", "fk_information_country_managements_country", "CASCADE"); err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				countryTable, err := tableName(tx, &CountryManagement{})
				if err != nil {
					return err
				}
				stageTable, err := tableName(tx, &StageManagement{})
				if err != nil {
					return err
				}
				countryStepsTable, err := tableName(tx, &CountryStepsManagement{})
				if err != nil {
					return err
				}
				infoTable, err := tableName(tx, &InformationCountryManagement{})
				if err != nil {
					return err
				}

				if err := replaceFK(tx, stageTable, "country_id", countryTable, "id", "fk_stage_managements_country", "RESTRICT"); err != nil {
					return err
				}
				if err := replaceFK(tx, countryStepsTable, "country_id", countryTable, "id", "fk_country_steps_managements_country", "RESTRICT"); err != nil {
					return err
				}
				if err := replaceFK(tx, infoTable, "country_id", countryTable, "id", "fk_information_country_managements_country", "RESTRICT"); err != nil {
					return err
				}

				return nil
			},
		},
		{
			ID: "20260504120000_visa_type_managements",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&VisaTypeManagement{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260505100000_users_add_source",
			Migrate: func(tx *gorm.DB) error {
				migrator := tx.Migrator()
				if !migrator.HasTable(&User{}) {
					return nil
				}

				if !migrator.HasColumn(&User{}, "Source") {
					if err := migrator.AddColumn(&User{}, "Source"); err != nil {
						return err
					}
				}

				// Create index for `source` using a safe explicit name (avoid GORM using field name as index name).
				if tx.Dialector.Name() == "mysql" {
					userTable, err := tableName(tx, &User{})
					if err != nil {
						return err
					}

					var existing []indexRef
					if err := tx.Raw(`
						SELECT INDEX_NAME
						FROM information_schema.statistics
						WHERE table_schema = DATABASE()
						  AND table_name = ?
						  AND column_name = 'source'
					`, userTable).Scan(&existing).Error; err != nil {
						return err
					}

					if len(existing) == 0 {
						if err := tx.Exec(fmt.Sprintf(
							"CREATE INDEX `idx_users_source` ON `%s` (`source`)",
							userTable,
						)).Error; err != nil {
							return err
						}
					}
				}

				return tx.AutoMigrate(&User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260505120000_users_visa_type_fk",
			Migrate: func(tx *gorm.DB) error {
				migrator := tx.Migrator()
				if !migrator.HasTable(&User{}) || !migrator.HasTable(&VisaTypeManagement{}) {
					return nil
				}

				userTable, err := tableName(tx, &User{})
				if err != nil {
					return err
				}
				visaTable, err := tableName(tx, &VisaTypeManagement{})
				if err != nil {
					return err
				}

				// Normalize old data: if users.visa_type stored as label instead of id, try map by name.
				if tx.Dialector.Name() == "mysql" {
					if err := tx.Exec(fmt.Sprintf(`
						UPDATE %s u
						JOIN %s v ON TRIM(u.visa_type) = TRIM(v.name)
						SET u.visa_type = v.id
						WHERE u.visa_type IS NOT NULL
						  AND u.visa_type <> ''
						  AND LENGTH(u.visa_type) <> 25
					`, userTable, visaTable)).Error; err != nil {
						return err
					}

					// Any remaining non-id values -> NULL
					if err := tx.Exec(fmt.Sprintf(`
						UPDATE %s
						SET visa_type = NULL
						WHERE visa_type IS NOT NULL
						  AND visa_type <> ''
						  AND LENGTH(visa_type) <> 25
					`, userTable)).Error; err != nil {
						return err
					}
				}

				// Ensure column size is compatible with visa_type_managements.id (25).
				if tx.Dialector.Name() == "mysql" {
					if err := tx.Exec(fmt.Sprintf(`
						ALTER TABLE %s
						MODIFY visa_type VARCHAR(25) NULL
					`, tx.Statement.Quote(userTable))).Error; err != nil {
						return err
					}
				}

				// Replace FK with ON DELETE SET NULL.
				return replaceFK(tx, userTable, "visa_type", visaTable, "id", "fk_users_visa_type", "SET NULL")
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260505153000_documents_managements_add_example_url",
			Migrate: func(tx *gorm.DB) error {
				migrator := tx.Migrator()
				if !migrator.HasTable(&DocumentsManagement{}) {
					return nil
				}

				if !migrator.HasColumn(&DocumentsManagement{}, "ExampleURL") {
					if err := migrator.AddColumn(&DocumentsManagement{}, "ExampleURL"); err != nil {
						return err
					}
				}

				return tx.AutoMigrate(&DocumentsManagement{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260511120000_soft_delete_master_data",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(
					&DocumentsManagement{},
					&VisaTypeManagement{},
					&QuestionBase{},
					&Question{},
					&QuestionOption{},
				); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260512110000_users_and_document_translations_existing_flags",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&User{}, &DocumentTranslation{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260512170000_password_reset_otps",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&PasswordResetOTP{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260512173000_users_add_source_category",
			Migrate: func(tx *gorm.DB) error {
				migrator := tx.Migrator()
				if !migrator.HasTable(&User{}) {
					return nil
				}

				if !migrator.HasColumn(&User{}, "SourceCategory") {
					if err := migrator.AddColumn(&User{}, "SourceCategory"); err != nil {
						return err
					}
				}

				if tx.Dialector.Name() == "mysql" {
					userTable, err := tableName(tx, &User{})
					if err != nil {
						return err
					}

					var existing []indexRef
					if err := tx.Raw(`
						SELECT INDEX_NAME
						FROM information_schema.statistics
						WHERE table_schema = DATABASE()
						  AND table_name = ?
						  AND column_name = 'source_category'
					`, userTable).Scan(&existing).Error; err != nil {
						return err
					}

					if len(existing) == 0 {
						if err := tx.Exec(fmt.Sprintf(
							"CREATE INDEX `idx_users_source_category` ON `%s` (`source_category`)",
							userTable,
						)).Error; err != nil {
							return err
						}
					}
				}

				return tx.AutoMigrate(&User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		{
			ID: "20260512190000_promos_add_is_active",
			Migrate: func(tx *gorm.DB) error {
				migrator := tx.Migrator()
				if !migrator.HasTable(&Promo{}) {
					if err := tx.AutoMigrate(&Promo{}); err != nil {
						return err
					}
				}

				if !migrator.HasColumn(&Promo{}, "IsActive") {
					if err := migrator.AddColumn(&Promo{}, "IsActive"); err != nil {
						return err
					}
				}

				if err := tx.Model(&Promo{}).Where("is_active IS NULL").Update("is_active", true).Error; err != nil {
					return err
				}

				return tx.AutoMigrate(&Promo{})
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
