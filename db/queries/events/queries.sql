-- -- name: ListUnprocessedEvents :many
-- SELECT * FROM events
-- WHERE processed_at IS NULL
-- ORDER BY occurred_at
-- LIMIT 100
-- FOR UPDATE;

-- -- name: GetNextUnprocessedEvent :one
-- SELECT * FROM events
-- WHERE processed_at IS NULL
-- ORDER BY occurred_at
-- LIMIT 1
-- FOR UPDATE SKIP LOCKED;

-- name: GetEventToProcess :one
WITH events_to_process as (
    SELECT * FROM events
    WHERE status = 'pending'
    OR (status = 'retry' AND now() >= next_retry_at)
    ORDER BY occurred_at
    LIMIT 1
    FOR UPDATE SKIP LOCKED
)
UPDATE events e
SET status = 'processing'
FROM events_to_process ep
WHERE e.id = ep.id
RETURNING e.*;