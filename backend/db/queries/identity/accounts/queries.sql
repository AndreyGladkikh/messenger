-- name: GetAccount :one
SELECT * FROM identity.accounts
WHERE id = $1;
