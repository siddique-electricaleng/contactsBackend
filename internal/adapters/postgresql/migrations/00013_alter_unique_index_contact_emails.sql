-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS uniq_contact_normalized_email;

CREATE UNIQUE INDEX IF NOT EXISTS contact_emails_user_norm_uniq
ON contact_emails (user_id, normalized_email);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS contact_emails_user_norm_uniq;

CREATE UNIQUE INDEX IF NOT EXISTS uniq_contact_normalized_email
ON contact_emails (contact_id, normalized_email);

-- +goose StatementEnd
