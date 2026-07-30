-- name: CreateUser :exec
INSERT INTO auth.users (
  id,
  login,
  password_hash
) VALUES (
  $1, $2, $3
);