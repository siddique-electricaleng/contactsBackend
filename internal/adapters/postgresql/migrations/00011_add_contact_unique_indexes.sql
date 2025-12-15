-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS contacts_no_exact_dupe_active
ON contacts(user_id, display_name, first_name, surname)
WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS contacts_no_exact_dupe_active;
-- +goose StatementEnd
