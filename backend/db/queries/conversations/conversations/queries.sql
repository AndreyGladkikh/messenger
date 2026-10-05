-- name: GetConversation :one
SELECT * FROM conversations.conversations
WHERE id = $1;

-- name: GetDirectConversationBetween :one
SELECT * FROM conversations.conversations
WHERE id = (
    SELECT conversation_id FROM conversations.direct_conversations
    WHERE participant_id_low = $1 AND participant_id_high = $2
);

-- name: ListConversationsForUser :many
SELECT * FROM conversations.conversations
WHERE id = ANY(
    SELECT conversation_id
    FROM conversations.memberships
    WHERE participant_id = $1
)
ORDER BY created_at DESC
LIMIT $2 
OFFSET $3;

-- name: ListConversationIDsForUser :many
SELECT id FROM conversations.conversations
WHERE id = ANY(
    SELECT conversation_id
    FROM conversations.memberships
    WHERE participant_id = $1
);

-- name: DirectConversationExists :one
SELECT EXISTS(
    SELECT 1 
    FROM conversations.conversations c
    WHERE c.type = 'direct'
    AND EXISTS(SELECT 1 FROM conversations.memberships cp WHERE cp.conversation_id = c.id AND cp.participant_id = $1)
    AND EXISTS(SELECT 1 FROM conversations.memberships cp WHERE cp.conversation_id = c.id AND cp.participant_id = $2)
);