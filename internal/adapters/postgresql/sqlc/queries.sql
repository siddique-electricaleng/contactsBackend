-- SQLC queries for user management : auth domain
-- name: CreateUser :one
INSERT INTO users (
    email,
    username,
    password_hash,
    name
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByEmailOrUsername :one
SELECT *
FROM users
WHERE (email = $1 OR username = $1)
  AND deleted_at IS NULL
LIMIT 1;

-- name: MarkUserEmailVerified :exec
UPDATE users
SET email_verified = TRUE,
    updated_at     = now()
WHERE id = $1;

-- SQLC queries for email verification tokens

-- name: CreateEmailVerificationToken :one
INSERT INTO email_verification_tokens (
    user_id,
    token,
    expires_at
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetEmailVerificationToken :one
SELECT *
FROM email_verification_tokens
WHERE token = $1;

-- name: MarkEmailVerificationTokenUsed :exec
UPDATE email_verification_tokens
SET used_at = now()
WHERE id = $1;

-- name: DeleteEmailVerificationToken :exec
DELETE FROM email_verification_tokens
WHERE id = $1;

-- SQLC queries for refresh tokens

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    user_id,
    token_hash,
    expires_at
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetRefreshTokenByHash :one
SELECT *
FROM refresh_tokens
WHERE token_hash = $1
LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = now()
WHERE id = $1
  AND revoked_at IS NULL;

-- name: RevokeAllUserRefreshTokens :exec
UPDATE refresh_tokens
SET revoked_at = now()
WHERE user_id = $1
  AND revoked_at IS NULL;

-- SQLC queries for contact management : contacts domain

-- POST /contacts

-- name: CreateContact :one
INSERT INTO contacts (
    user_id,
    display_name,
    first_name,
    surname,
    note,
    source
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpsertContactPhone :one
INSERT INTO contact_phone_numbers(
  user_id,
  contact_id,
  number,
  label,
  normalized_number,
  is_primary
)
VALUES( $1, $2, $3, $4, $5, $6
)
ON CONFLICT (user_id, normalized_number) DO UPDATE
SET
  contact_id = EXCLUDED.contact_id,
  number = EXCLUDED.number,
  label = EXCLUDED.label,
  is_primary = EXCLUDED.is_primary,
  updated_at = now()
RETURNING *
;

-- name: UpsertContactEmail :one
INSERT INTO contact_emails(
  user_id,
  contact_id,
  email,
  label,
  normalized_email,
  is_primary
)
VALUES( $1, $2, $3, $4, $5, $6
) ON CONFLICT (user_id, normalized_email) DO UPDATE
SET 
  contact_id = EXCLUDED.contact_id,
  email = EXCLUDED.email,
  label = EXCLUDED.label,
  is_primary = EXCLUDED.is_primary,
  updated_at = now()
RETURNING *
;

-- GET /contacts?limit=&offset=

-- name: ListContactsForUser :many
SELECT *
FROM contacts
WHERE user_id = $1
  AND deleted_at IS NULL
ORDER BY display_name ASC;

-- GET /contacts with details
-- name: ListContactsForUserWithDetails :many
SELECT
	c.id,
	c.display_name,
	c.first_name,
	c.surname,
	c.note,
	c.source,
	COALESCE(
		json_agg(
			distinct jsonb_build_object(
			'label', p.label,
			'number', p.number,
			'isPrimary', p.is_primary
			)
		) filter (where p.id is not null)
	, '[]'::json) as phone_numbers,
	COALESCE(
		json_agg(
			distinct jsonb_build_object(
			'label', e.label,
			'email', e.email,
			'isPrimary', e.is_primary
			)
		) filter (where e.id is not null)
	, '[]'::json) as emails
FROM contacts c
LEFT JOIN contact_phone_numbers p 
		ON p.contact_id = c.id 
		AND p.user_id = c.user_id 
LEFT JOIN contact_emails e 
		ON e.contact_id = c.id
		AND e.user_id = c.user_id 
WHERE c.user_id = $1
	and c.deleted_at is null
group by c.id, c.display_name, c.first_name, c.surname, c.note, c.source
order by c.first_name asc
;

-- GET /contacts/{id}

-- name: GetContactByID :one
SELECT *
FROM contacts
WHERE id = $1
  AND user_id = $2
  AND deleted_at IS NULL
;

-- name: ListPhonesForContact :many
SELECT *
FROM contact_phone_numbers
WHERE contact_id = $1;

-- name: ListEmailsForContact :many
SELECT *
FROM contact_emails
WHERE contact_id = $1;
