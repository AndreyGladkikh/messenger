-- name: AddParticipantsToChat :copyfrom
INSERT INTO chat_participants (
    id,
    chat_id,
    participant_id,
    role
) VALUES (
  $1, $2, $3, $4
);
