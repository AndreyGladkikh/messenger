-- name: CreateChat :exec
INSERT INTO messaging.chats (
    id,
    kind,
    name
) VALUES (
  $1, $2, $3
);

-- name: CreatePrivateChat :exec
INSERT INTO messaging.private_chats (
    chat_id,
    first_participant_id,
    second_participant_id
) VALUES (
  $1, $2, $3
);

-- name: DeleteChat :exec
UPDATE messaging.chats
  set deleted_at = $2
WHERE id = $1;