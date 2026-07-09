-- name: CreateMessage :exec
INSERT INTO messages (
    id,
    sender_id,
    chat_id,
    body,
    reply_to_message_id,
    created_at
) VALUES (
  $1, $2, $3, $4, $5, $6
);

-- name: UpdateMessage :exec
UPDATE messages
  set body = $2,
  updated_at = $3
WHERE id = $1;

-- name: DeleteMessage :exec
UPDATE messages
  set deleted_at = $2
WHERE id = $1;