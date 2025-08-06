-- +migrate Up
CREATE TABLE role (
                       id TEXT PRIMARY KEY,
                       short_id TEXT NOT NULL DEFAULT '',
                       name TEXT NOT NULL DEFAULT '',
                       description TEXT NOT NULL DEFAULT '',
                       contextual BOOLEAN DEFAULT 0 NOT NULL,
                       status TEXT NOT NULL DEFAULT '',
                       created_by TEXT,
                       updated_by TEXT,
                       created_at TIMESTAMP,
                       updated_at TIMESTAMP
);

-- +migrate Down
DROP TABLE role;
