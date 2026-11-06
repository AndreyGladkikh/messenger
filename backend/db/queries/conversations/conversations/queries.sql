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
SELECT cnv.id,
    cnv.type,
    gcnv.title,
    gcnv.avatar_url,
    msg.id as msg_id,
    msg.body as msg_body,
    msg.sender_id as msg_sender_id,
    msg.created_at as msg_created_at,
    other_memb.participant_id as direct_participant_id
FROM conversations.conversations cnv
LEFT JOIN conversations.group_conversations gcnv ON cnv.id = gcnv.conversation_id
JOIN conversations.memberships memb ON cnv.id = memb.conversation_id
LEFT JOIN conversations.memberships other_memb ON cnv.id = other_memb.conversation_id 
    AND other_memb.participant_id != memb.participant_id 
    AND cnv.type = 'direct'
LEFT JOIN lateral (
    SELECT msg.id,
        msg.body,
        msg.sender_id,
        msg.created_at
    FROM conversations.messages msg
    WHERE msg.conversation_id = cnv.id
        AND msg.deleted_at IS NULL
    ORDER BY msg.created_at DESC, msg.id DESC
    LIMIT 1
) msg ON true
WHERE memb.participant_id = $1
AND cnv.deleted_at IS NULL
ORDER BY msg.created_at DESC, cnv.created_at DESC
;

-- -- name: ListConversationsForUser :many
-- SELECT * FROM conversations.conversations c
-- LEFT JOIN conversations.group_conversations gc ON c.id = gc.conversation_id
-- WHERE id = ANY(
--     SELECT conversation_id
--     FROM conversations.memberships
--     WHERE participant_id = $1
-- )
-- ORDER BY c.created_at DESC
-- LIMIT $2 
-- OFFSET $3;

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