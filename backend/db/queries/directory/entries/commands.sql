-- name: UpsertDirectoryEntry :exec
INSERT INTO directory.entries (
    id,
    target_id,
    target_type,
    display_name,
    handle
) VALUES (
  $1, $2, $3, $4, $5
) ON CONFLICT (id) DO UPDATE
SET display_name = EXCLUDED.display_name,
  handle = EXCLUDED.handle;

-- name: UpdateDirectoryEntry :exec
UPDATE directory.entries
  set display_name = $3,
  handle = $4
WHERE target_id = $1
AND target_type = $2;
