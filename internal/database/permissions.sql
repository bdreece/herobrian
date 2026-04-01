-- name: ListPermissions :many
SELECT p.*
  FROM permissions AS p
 WHERE p.id > @after
 ORDER BY p.id ASC
 LIMIT @first;

-- name: UpsertPermission :execrows
INSERT INTO permissions (resource, method)
VALUES (@id, @resource, @method)
ON CONFLICT DO NOTHING;
