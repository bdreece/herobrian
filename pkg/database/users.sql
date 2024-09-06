-- name: FindUser :one
SELECT *
FROM users
WHERE id = @id
LIMIT 1;

-- name: FindUserByUsername :one
SELECT *
FROM users
WHERE username = @username
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (username, password_hash, role_id)
VALUES (@username, @password_hash, @role_id)
RETURNING id;

-- name: UpsertUser :execrows
REPLACE INTO users (id, username, password_hash, role_id)
VALUES (@id, @username, @password_hash, @role_id);

-- name: UpdateUser :execrows
UPDATE users
SET username = @username,
    password_hash = @password_hash,
    role_id = @role_id
WHERE id = @id;

-- name: RemoveUser :execrows
DELETE FROM users
WHERE id = @id;
