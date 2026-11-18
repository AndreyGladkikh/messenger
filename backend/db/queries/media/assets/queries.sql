-- name: GetMediaAsset :one
SELECT
    *
FROM
    media.assets
WHERE
    id = $1
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
