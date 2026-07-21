-- name: CreateChat :exec
INSERT INTO chats (
    id,
    type,
    name
) VALUES (
  $1, $2, $3
);

-- name: DeleteChat :exec
UPDATE chats
  set deleted_at = $2
WHERE id = $1;