CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email_address TEXT NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT false,
    password_hash TEXT NOT NULL
);

CREATE UNIQUE INDEX IX_users_username ON users (username ASC);
CREATE UNIQUE INDEX IX_users_email_address ON users (email_address ASC);
CREATE INDEX IX_users_name ON users (last_name ASC, first_name ASC);

CREATE TABLE roles (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NULL
);

CREATE UNIQUE INDEX IX_roles_name ON roles (name ASC);

CREATE TABLE user_roles (
    user_id INTEGER NOT NULL REFERENCES users (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    role_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, role_id)
);
