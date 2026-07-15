-- name: ListUnprocessedEvents :many
SELECT * FROM events
WHERE processed_at IS NULL
ORDER BY occurred_at
LIMIT 100;