-- name: UpsertMediaAsset :exec
INSERT INTO media.assets (
  id,
  status,
  name
) VALUES (
  $1, $2, $3
)
ON CONFLICT (id) DO UPDATE
SET status = EXCLUDED.status;

-- name: UpdateMediaAssetLocation :exec
UPDATE media.assets
SET bucket = $2,
  key = $3
WHERE id = $1;
