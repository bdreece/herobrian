-- name: FindUserByUsername :one
SELECT *
FROM users
WHERE username = @username
LIMIT 1;

-- name: FindUserByEmail :one
SELECT *
FROM users
WHERE email_address = @email_address
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (username, first_name, last_name, email_address, password_hash)
VALUES (@username, @first_name, @last_name, @email_address, @password_hash)
RETURNING id;

-- name: UpdateUser :execrows
UPDATE users
SET username = @username,
    first_name = @first_name,
    last_name = @last_name,
    email_address = @email_address,
    email_verified = @email_verified,
    password_hash = @password_hash
WHERE id = @id;

-- name: RemoveUser :execrows
DELETE FROM users
WHERE id = @id;
