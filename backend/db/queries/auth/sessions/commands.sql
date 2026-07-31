-- name: CreateSession :exec
INSERT INTO auth.sessions (
  id,
  user_id,
  refresh_token_hash,
  expires_at,
  user_agent,
  ip
) VALUES (
  $1, $2, $3, $4, $5, $6
);

-- name: SaveSession :exec
UPDATE auth.sessions
SET refresh_token_hash = $2,
  expires_at = $3,
  revoked_at = $4,
  last_used_at = $5
WHERE id = $1;