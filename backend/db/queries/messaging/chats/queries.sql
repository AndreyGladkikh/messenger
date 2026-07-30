-- name: GetChat :one
SELECT * FROM messaging.chats
WHERE id = $1;

-- name: ListChatsForUser :many
SELECT * FROM messaging.chats
WHERE id = ANY(
    SELECT chat_id
    FROM messaging.chat_participants
    WHERE participant_id = $1
)
ORDER BY created_at DESC
LIMIT $2 
OFFSET $3;

-- name: PrivateChatExists :one
SELECT EXISTS(
    SELECT 1 
    FROM messaging.chats c
    WHERE c.kind = 'private_chat'
    AND EXISTS(SELECT 1 FROM messaging.chat_participants cp WHERE cp.chat_id = c.id AND cp.participant_id = $1)
    AND EXISTS(SELECT 1 FROM messaging.chat_participants cp WHERE cp.chat_id = c.id AND cp.participant_id = $2)
);