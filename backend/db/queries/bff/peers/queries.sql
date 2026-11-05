-- name: GetUserPeers :many
SELECT *
FROM bff.peers p
WHERE p.user_id = $1
ORDER BY p.preview_time DESC
LIMIT $2;
