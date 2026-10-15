-- name: HandleExists :one
SELECT EXISTS(
    SELECT 1 FROM directory.entries
    WHERE handle = $1
);

-- name: SearchDirectoryEntry :many
SELECT target_id, target_type, display_name, handle 
FROM directory.entries
WHERE ts_vec @@ to_tsquery('russian', $1);
