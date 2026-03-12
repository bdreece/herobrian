CREATE TABLE IF NOT EXISTS users (
	id            INTEGER  NOT NULL,
	created_at    DATETIME NOT NULL DEFAULT (DATETIME('now')),
	updated_at    DATETIME NOT NULL DEFAULT (DATETIME('now')),
	first_name    TEXT     NOT NULL,
	last_name     TEXT     NOT NULL,
	display_name  TEXT     NOT NULL,
	password_hash TEXT     NOT NULL,
	role          TEXT     NOT NULL   CHECK ( role IN ('landlubber', 'scallywag', 'freebooter', 'privateer', 'swashbuckler') ),
	picture_url   TEXT         NULL,
	totp_secret   TEXT         NULL,

	PRIMARY KEY ( id ),
	UNIQUE ( display_name )
);

CREATE INDEX IF NOT EXISTS users_by_name ON users (
	last_name ASC,
	first_name ASC,
	created_at DESC
);

CREATE INDEX IF NOT EXISTS users_by_display_name ON users ( display_name ASC );
