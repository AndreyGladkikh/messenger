-- name: PutToOutbox :exec
INSERT INTO outbox (
    id,
    kind,
    type,
    payload
) VALUES (
  $1, $2, $3, $4
);

-- name: ProcessEvent :exec
UPDATE outbox
  set status = $2,
  attempts = $3,
  error = $4,
  next_retry_at = $5
WHERE id = $1;