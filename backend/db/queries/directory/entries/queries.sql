-- name: HandleExists :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            directory.handles
        WHERE
            handle = $1
    )
;

-- name: SearchDirectoryEntry :many
SELECT
    target_id,
    target_type,
    display_name,
    handle,
    avatar_id
FROM
    directory.search_index
WHERE
    ts_vec @@ TO_TSQUERY('russian', $1)
;

-- name: GetUserProfile :one
SELECT
    *
FROM
    directory.user_profiles
WHERE
    user_id = $1
;

-- name: GetUserProfilesByIDs :many
SELECT
    *
FROM
    directory.user_profiles p
    LEFT JOIN directory.handles h ON p.user_id = h.target_id
    AND h.target_type = 'user'
WHERE
    user_id = ANY ($1::UUID[])
;
