-- name: CreateEventHandlerExecution :one
INSERT INTO event_handler_executions (
  id,
  event_id,
  handler_type
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateEventHandlerExecution :exec
UPDATE event_handler_executions
  set attempts = $2,
    error = $3,
    next_retry_at = $4
WHERE id = $1;