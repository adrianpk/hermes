-- +migrate Up
CREATE TABLE user (
    id TEXT PRIMARY KEY,
    short_id TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL DEFAULT '',
    username TEXT NOT NULL DEFAULT '',
    email_enc BLOB,
    password_enc BLOB,
    created_by TEXT,
    updated_by TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    last_login_at TIMESTAMP,
    last_login_ip TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN DEFAULT 1
);

CREATE INDEX idx_user_email_enc ON user(email_enc);

-- +migrate Down
DROP TABLE user;
