-- name: ListRoles :many
SELECT r.*
  FROM roles AS r
 WHERE r.name > @after
 ORDER BY r.name ASC
 LIMIT @first;

-- name: FindRoleByID :one
SELECT r.*
  FROM roles AS r
 WHERE r.id = @id
 LIMIT 1;

-- name: UpsertRole :execrows
INSERT INTO roles (name)
VALUES (@name)
ON CONFLICT DO NOTHING;
