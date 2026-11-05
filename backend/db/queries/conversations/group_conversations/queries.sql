-- name: GetGroupConversation :one
SELECT *
FROM conversations.group_conversations
WHERE conversation_id = $1;