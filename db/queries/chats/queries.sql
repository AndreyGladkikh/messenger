-- name: GetChat :one
SELECT * FROM chats
WHERE id = $1 LIMIT 1;

-- name: ListChatsForUser :many
SELECT * FROM chats
WHERE id = ANY(
    SELECT chat_id
    FROM chat_participants
    WHERE participant_id = $1
)
ORDER BY created_at DESC
LIMIT $2 
OFFSET $3;

-- -- name: PrivateChatExists :one
-- SELECT 1 FROM private_chats
-- WHERE user1_id = $1
-- AND user2_id = $2;

-- name: PrivateChatExists :one
SELECT EXISTS (
    SELECT 1 
    FROM chats c
    WHERE c.type = 'private_chat'
    AND EXISTS(SELECT 1 FROM chat_participants cp WHERE cp.chat_id = c.id AND cp.participant_id = $1)
    AND EXISTS(SELECT 1 FROM chat_participants cp WHERE cp.chat_id = c.id AND cp.participant_id = $2)
);