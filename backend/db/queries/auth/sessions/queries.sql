-- name: GetSessionByRefreshTokenHash :one
SELECT * FROM auth.sessions
WHERE refresh_token_hash = $1;