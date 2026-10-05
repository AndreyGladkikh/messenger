-- name: CreateAccount :exec
INSERT INTO identity.accounts (
  id,
  type
) VALUES (
  $1, $2
);

-- name: CreateUser :exec
INSERT INTO identity.users (
  id,
  login,
  password_hash
) VALUES (
  $1, $2, $3
);
