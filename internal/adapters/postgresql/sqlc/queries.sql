-- SQLC queries for user management
-- name: CreateUser :one
INSERT INTO users (
    email,
    username,
    password_hash,
    name,
) VALUES (
    $1, $2, $3, $4, $5
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
