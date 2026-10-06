-- name: UpsertSession :exec
INSERT INTO identity.sessions (
  id,
  user_id,
  refresh_token_hash,
  expires_at,
  user_agent,
  ip,
  revoked_at,
  last_used_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
ON CONFLICT (id) DO UPDATE
SET refresh_token_hash = EXCLUDED.refresh_token_hash,
  expires_at = EXCLUDED.expires_at,
  user_agent = EXCLUDED.user_agent,
  ip = EXCLUDED.ip,
  revoked_at = EXCLUDED.revoked_at,
  last_used_at = EXCLUDED.last_used_at;
