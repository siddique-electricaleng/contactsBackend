-- +goose Up
CREATE TABLE IF NOT EXISTS contact_phone_numbers (
    id bigserial PRIMARY KEY,

    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    contact_id uuid NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,

    label text,                       -- e.g. mobile, home, work
    number text NOT NULL,             -- raw
    normalized_number text NOT NULL,  -- E.164 or cleaned form

    is_primary boolean NOT NULL DEFAULT false,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- Unique: user cannot store same normalized number twice
CREATE UNIQUE INDEX unique_user_normalized_phone
    ON contact_phone_numbers(user_id, normalized_number);

-- +goose Down
DROP TABLE IF EXISTS contact_phone_numbers;