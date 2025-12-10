-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS unique_user_normalized_phone;

DROP INDEX IF EXISTS unique_user_normalized_email;

CREATE UNIQUE INDEX uniq_contact_normalized_phone
  ON contact_phone_numbers(contact_id, normalized_number);

CREATE UNIQUE INDEX uniq_contact_normalized_email
  ON contact_emails(contact_id, normalized_email);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS uniq_contact_normalized_phone;
DROP INDEX IF EXISTS uniq_contact_normalized_email;
-- +goose StatementEnd
