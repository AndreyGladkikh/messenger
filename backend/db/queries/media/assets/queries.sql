-- name: GetMediaAsset :one
SELECT * FROM media.assets
WHERE id = $1;
