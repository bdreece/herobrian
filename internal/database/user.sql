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

-- name: UpsertUser :execrows
INSERT INTO users
(
	first_name,
	last_name,
	display_name,
	password_hash,
	picture_url,
	totp_secret
)
VALUES
(
	@first_name,
	@last_name,
	@display_name,
	@password_hash,
	@picture_url,
	@totp_secret
) ON CONFLICT (display_name) DO UPDATE 
  SET updated_at = DATETIME('now'),
      first_name = ?1,
      last_name = ?2,
      password_hash = ?4,
      picture_url = ?5,
      totp_secret = ?6
WHERE display_name = ?3;
