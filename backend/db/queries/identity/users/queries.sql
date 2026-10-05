-- name: GetUserByLogin :one
SELECT * FROM identity.users
WHERE login = $1;

-- name: UserExistsByLogin :one
SELECT EXISTS(
    SELECT 1 
    FROM identity.users
    WHERE login = $1
);