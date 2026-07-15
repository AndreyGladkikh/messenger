-- name: GetChat :one
SELECT * FROM chats
WHERE id = $1 LIMIT 1;

-- name: ListChatsForUser :many
SELECT * FROM chats
WHERE chat_id = ANY(
    SELECT chat_id
    FROM chat_participants
    WHERE participant_id = $1
)
ORDER BY created_at DESC
LIMIT $2 
OFFSET $3;