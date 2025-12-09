-- +goose Up
CREATE TABLE IF NOT EXISTS contact_emails (
    id bigserial PRIMARY KEY,

    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    contact_id uuid NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,

    label text,
    email text NOT NULL,               -- raw
    normalized_email text NOT NULL,    -- lowercase/cleaned

    is_primary boolean NOT NULL DEFAULT false,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- Unique: user cannot store same normalized email twice
CREATE UNIQUE INDEX unique_user_normalized_email
    ON contact_emails(user_id, normalized_email);

-- +goose Down
DROP TABLE IF EXISTS contact_emails;