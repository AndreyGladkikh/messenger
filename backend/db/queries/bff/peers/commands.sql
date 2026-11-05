-- name: InsertWebBffPeer :exec
INSERT INTO bff.peers (user_id, target_id, target_type, preview_time)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, target_id, target_type) DO UPDATE SET preview_time = $4;

-- name: UpdateWebBffPeer :exec
UPDATE bff.peers
SET preview_time = $1
WHERE user_id = $2 AND target_id = $3 AND target_type = $4;
