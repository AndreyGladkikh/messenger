-- name: GetSessionByRefreshTokenHash :one
SELECT * FROM identity.sessions
WHERE refresh_token_hash = $1;