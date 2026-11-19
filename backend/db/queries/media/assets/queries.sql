-- name: GetMediaAsset :one
SELECT
    *
FROM
    media.assets
WHERE
    id = $1
;

-- name: GetMediaAssets :many
SELECT
    *
FROM
    media.assets
WHERE
    id = ANY ($1::UUID[])
;

-- name: MediaAssetExists :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            media.assets
        WHERE
            id = $1
    )
;
