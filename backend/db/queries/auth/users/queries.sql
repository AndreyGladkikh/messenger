-- name: GetUserByLogin :one
SELECT * FROM auth.users
WHERE login = $1;

-- name: UserExistsByLogin :one
SELECT EXISTS(
    SELECT 1 
    FROM auth.users
    WHERE login = $1
);