-- +migrate Up
CREATE TABLE section (
    id TEXT PRIMARY KEY,
    short_id TEXT,
    name TEXT NOT NULL,
    description TEXT,
    path TEXT,
    layout_id TEXT,
    image TEXT,
    header TEXT,
    created_by TEXT,
    updated_by TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- +migrate Down
DROP TABLE section;