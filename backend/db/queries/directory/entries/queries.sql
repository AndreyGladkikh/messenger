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

-- name: ListDirectoryEntriesByTargets :many
WITH targets AS (
  SELECT 
    target_ids.val::uuid AS id,
    target_types.val::text AS type
  FROM unnest(@target_ids::uuid[]) WITH ORDINALITY AS target_ids(val, ord)
  JOIN unnest(@target_types::text[]) WITH ORDINALITY AS target_types(val, ord) ON target_ids.ord = target_types.ord
)
SELECT
    e.id,
    e.target_id,
    e.target_type,
    e.display_name,
    e.handle,
    e.avatar_url
FROM directory.entries AS e
JOIN targets ON e.target_id = targets.id AND e.target_type = targets.type;
