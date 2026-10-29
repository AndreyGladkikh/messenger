-- name: HandleExists :one
SELECT EXISTS(
    SELECT 1 FROM directory.entries
    WHERE handle = $1
);

-- name: SearchDirectoryEntry :many
SELECT target_id, target_type, display_name, handle 
FROM directory.entries
WHERE ts_vec @@ to_tsquery('russian', $1);

-- name: GetDirectoryEntry :one
SELECT display_name, handle, avatar_url 
FROM directory.entries
WHERE target_id = $1
    AND target_type = $2;
