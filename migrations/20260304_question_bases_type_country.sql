-- Choose the statement that matches your current schema.
-- If the column is currently named `type`:
-- ALTER TABLE question_bases CHANGE COLUMN `type` `type_country` VARCHAR(50) NOT NULL;

-- If the column does not exist yet:
-- ALTER TABLE question_bases ADD COLUMN `type_country` VARCHAR(50) NOT NULL;
