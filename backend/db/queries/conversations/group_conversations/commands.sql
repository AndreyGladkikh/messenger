-- name: UpsertGroupConversation :exec
INSERT INTO conversations.group_conversations(
    conversation_id,
    title,
    avatar_url
) VALUES (
    $1, $2, $3
) ON CONFLICT (conversation_id) DO UPDATE
SET title = EXCLUDED.title,
    avatar_url = EXCLUDED.avatar_url
;