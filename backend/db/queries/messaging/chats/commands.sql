-- name: CreateChat :exec
INSERT INTO messaging.chats (
    id,
    kind,
    name
) VALUES (
  $1, $2, $3
);

-- name: DeleteChat :exec
UPDATE messaging.chats
  set deleted_at = $2
WHERE id = $1;