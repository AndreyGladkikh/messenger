-- name: UpsertMediaUpload :exec
INSERT INTO media.uploads (
  id,
  status,
  initiator_id,
  asset_id
) VALUES (
  $1, $2, $3, $4
)
ON CONFLICT (id) DO UPDATE
SET status = EXCLUDED.status,
  asset_id = EXCLUDED.asset_id
;
