-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS uniq_contact_normalized_phone;

CREATE UNIQUE INDEX IF NOT EXISTS contact_phone_numbers_user_norm_uniq
ON contact_phone_numbers (user_id, normalized_number);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS contact_phone_numbers_user_norm_uniq;

CREATE UNIQUE INDEX IF NOT EXISTS uniq_contact_normalized_phone
ON contact_phone_numbers (contact_id, normalized_number);
-- +goose StatementEnd
