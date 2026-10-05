-- name: GetMessage :one
SELECT * FROM conversations.messages
WHERE id = $1 LIMIT 1;

-- name: ListMessagesForConversation :many
SELECT * FROM conversations.messages
WHERE conversation_id = $1
ORDER BY created_at DESC
LIMIT $2 
OFFSET $3;