-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX uniq_primary_phone_per_contact
    ON contact_phone_numbers (contact_id)
    WHERE is_primary = TRUE
    ;
CREATE UNIQUE INDEX uniq_primary_email_per_contact
    ON contact_emails (contact_id)
    WHERE is_primary = TRUE
    ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS uniq_primary_phone_per_contact;
DROP INDEX IF EXISTS uniq_primary_email_per_contact;
-- +goose StatementEnd
