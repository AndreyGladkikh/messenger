-- name: ConversationMembershipExists :one
SELECT EXISTS(
    SELECT 1 
    FROM conversations.memberships m
    WHERE m.conversation_id = $1
    AND m.participant_id = $2
);