-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_contacts_user_id_display_name ON contacts(user_id, display_name);
CREATE INDEX idx_contacts_phone_numbers_contact_id ON contact_phone_numbers(contact_id);
CREATE INDEX idx_contacts_emails_contact_id ON contact_emails(contact_id);                                             
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_contacts_user_id_display_name;
DROP INDEX IF EXISTS idx_contacts_phone_numbers_contact_id;
DROP INDEX IF EXISTS idx_contacts_emails_contact_id;
-- +goose StatementEnd
