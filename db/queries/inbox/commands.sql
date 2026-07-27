-- name: CreateInbox :one
INSERT INTO inbox (
    event_id,
    handler,
    executed_at
) VALUES (
  $1, $2, now()
)
ON CONFLICT DO NOTHING
RETURNING *;