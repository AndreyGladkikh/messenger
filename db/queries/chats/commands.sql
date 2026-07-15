-- name: CreateChat :exec
INSERT INTO chats (
    id,
    type,
    name,
    created_at
) VALUES (
  $1, $2, $3, $4
);

-- name: DeleteChat :exec
UPDATE chats
  set deleted_at = $2
WHERE id = $1;