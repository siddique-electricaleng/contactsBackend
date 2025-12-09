-- +goose Up
-- +goose StatementBegin
ALTER TABLE contacts
    RENAME COLUMN given_name TO first_name;

ALTER TABLE contacts
    RENAME COLUMN family_name TO surname;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE contacts
    RENAME COLUMN first_name TO given_name;

ALTER TABLE contacts
    RENAME COLUMN surname TO family_name;
-- +goose StatementEnd
