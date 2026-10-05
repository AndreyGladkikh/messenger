-- name: CreateConversation :exec
INSERT INTO conversations.conversations (
    id,
    type
) VALUES (
  $1, $2
);

-- name: CreateDirectConversation :exec
INSERT INTO conversations.direct_conversations (
    conversation_id,
    participant_id_low,
    participant_id_high
) VALUES (
  $1, $2, $3
);

-- name: DeleteConversation :exec
UPDATE conversations.conversations
  set deleted_at = $2
WHERE id = $1;