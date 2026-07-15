-- name: ProcessEvent :exec
UPDATE events
  set processed_at = $2
WHERE id = $1;