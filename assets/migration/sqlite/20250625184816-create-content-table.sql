-- +migrate Up
CREATE TABLE content (
    id TEXT PRIMARY KEY,
    short_id TEXT,
    user_id TEXT NOT NULL,
    heading TEXT NOT NULL,
    body TEXT,
    status TEXT,
    created_by TEXT,
    updated_by TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- +migrate Down
DROP TABLE content;
