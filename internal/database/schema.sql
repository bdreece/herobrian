CREATE TABLE IF NOT EXISTS permissions (
	id INTEGER NOT NULL,
	resource TEXT NOT NULL,
	method TEXT NOT NULL,

	PRIMARY KEY ( id )
);

CREATE UNIQUE INDEX IF NOT EXISTS permissions_by_resource_method ON permissions ( resource ASC, method ASC );

CREATE TABLE IF NOT EXISTS roles (
	id INTEGER NOT NULL,
	name TEXT NOT NULL,

	PRIMARY KEY ( id )
);

CREATE UNIQUE INDEX IF NOT EXISTS roles_by_name ON roles ( name ASC );

CREATE TABLE IF NOT EXISTS role_permissions (
	role_id INTEGER NOT NULL,
	permission_id INTEGER NOT NULL,

	PRIMARY KEY ( role_id, permission_id ),

	FOREIGN KEY ( role_id ) REFERENCES roles ( id )
		ON DELETE CASCADE
		ON UPDATE RESTRICT,
	
	FOREIGN KEY ( permission_id ) REFERENCES permissions ( id )
		ON DELETE CASCADE
		ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS users (
	id            INTEGER  NOT NULL,
	created_at    DATETIME NOT NULL DEFAULT (DATETIME('now')),
	updated_at    DATETIME NOT NULL DEFAULT (DATETIME('now')),
	first_name    TEXT     NOT NULL,
	last_name     TEXT     NOT NULL,
	display_name  TEXT     NOT NULL,
	password_hash TEXT     NOT NULL,
	picture_url   TEXT         NULL,
	totp_secret   TEXT         NULL,

	PRIMARY KEY ( id )
);

CREATE UNIQUE INDEX IF NOT EXISTS users_by_display_name ON users ( display_name ASC );

CREATE INDEX IF NOT EXISTS users_by_name ON users ( last_name ASC, first_name ASC, created_at DESC );

CREATE TABLE IF NOT EXISTS user_roles (
	user_id INTEGER NOT NULL,
	role_id INTEGER NOT NULL,

	PRIMARY KEY ( user_id, role_id ),

	FOREIGN KEY ( user_id ) REFERENCES users ( id )
		ON DELETE CASCADE
		ON UPDATE RESTRICT,
	
	FOREIGN KEY ( role_id ) REFERENCES roles ( id )
		ON DELETE CASCADE
		ON UPDATE RESTRICT
);
