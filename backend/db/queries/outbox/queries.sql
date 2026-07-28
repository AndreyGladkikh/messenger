-- name: GetEventToProcess :one
WITH events_to_process as (
    SELECT * FROM outbox
    WHERE status = 'pending'
    OR (status = 'retry' AND now() >= next_retry_at)
    ORDER BY occurred_at
    LIMIT 1
    FOR UPDATE SKIP LOCKED
)
UPDATE outbox e
SET status = 'processing',
claimed_at = now()
FROM events_to_process ep
WHERE e.id = ep.id
RETURNING e.*;