-- +migrate Up
CREATE TABLE permission (
                             id TEXT PRIMARY KEY,
                             short_id TEXT NOT NULL DEFAULT '',
                             name TEXT NOT NULL DEFAULT '',
                             description TEXT NOT NULL DEFAULT '',
                             created_by TEXT,
                             updated_by TEXT,
                             created_at TIMESTAMP,
                             updated_at TIMESTAMP
);

-- +migrate Down
DROP TABLE permission;
