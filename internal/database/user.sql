-- name: FindUserByID :one
SELECT u.*
FROM users AS u
WHERE u.id = @id
LIMIT 1;
