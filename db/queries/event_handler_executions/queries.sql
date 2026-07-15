-- name: ListEventHandlerExecutions :many
SELECT * FROM event_handler_executions
WHERE error IS NOT NULL
AND next_retry_at >= now()
LIMIT 100;

-- name: GetEventHandlerExecutionByEventId :one
SELECT * FROM event_handler_executions
WHERE event_id = $1
AND handler_type = $2
LIMIT 1;

-- name: GetEventHandlerExecutionsByEventId :many
SELECT * FROM event_handler_executions
WHERE event_id = $1;