-- name: AddParticipantsToConversation :copyfrom
INSERT INTO conversations.memberships (
    id,
    conversation_id,
    participant_id,
    role
) VALUES (
  $1, $2, $3, $4
);

-- name: CreateConversationMembership :exec
INSERT INTO conversations.memberships (
    id,
    conversation_id,
    participant_id,
    role
) VALUES (
  $1, $2, $3, $4
);
