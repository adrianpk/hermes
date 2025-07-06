-- +migrate Up
CREATE TABLE layout (
    id TEXT PRIMARY KEY,
    short_id TEXT,
    name TEXT NOT NULL,
    description TEXT,
    code TEXT,
    created_by TEXT,
    updated_by TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- +migrate Down
DROP TABLE layout;
