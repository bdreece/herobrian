-- name: FindUserByID :one
SELECT u.*
FROM users AS u
WHERE u.id = @id
LIMIT 1;

-- name: FindUserByDisplayName :one
SELECT u.*
FROM users AS u
WHERE u.display_name = @display_name
LIMIT 1;
