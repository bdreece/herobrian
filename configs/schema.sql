CREATE TABLE IF NOT EXISTS roles (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS IX_roles_name ON roles (name ASC);

REPLACE INTO roles
    (id, name)
VALUES
    (0, 'User'),
    (1, 'Moderator'),
    (2, 'Admin'),
    (3, 'Super');

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    role_id INTEGER NOT NULL REFERENCES roles (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

CREATE UNIQUE INDEX IF NOT EXISTS IX_users_username ON users (username ASC);
