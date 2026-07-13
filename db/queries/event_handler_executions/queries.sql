-- name: ListEventHandlerExecutions :many
SELECT * FROM event_handler_executions
WHERE error IS NOT NULL
AND next_retry_at >= now()
LIMIT 100;