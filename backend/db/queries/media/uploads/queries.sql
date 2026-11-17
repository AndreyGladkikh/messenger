-- name: GetUpload :one
SELECT * FROM media.uploads
WHERE id = $1;
