-- name: PutToOutbox :exec
INSERT INTO outbox (
    id,
    event_type,
    event_payload
) VALUES (
  $1, $2, $3
);

-- name: ProcessEvent :exec
UPDATE outbox
  set status = $2,
  claimed_at = $3,
  attempts = $4,
  error = $5,
  next_retry_at = $6
WHERE id = $1;