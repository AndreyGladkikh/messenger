-- name: GetMessage :one
SELECT * FROM messages
WHERE id = $1 LIMIT 1;

-- name: ListMessagesForChat :many
SELECT * FROM messages
WHERE chat_id = $1
ORDER BY created_at DESC
LIMIT $2 
OFFSET $3;