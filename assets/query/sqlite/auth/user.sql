-- Res: User
-- Table: user

-- GetAll
SELECT id, username, email_enc, password_enc, name, short_id, created_by, updated_by, created_at, updated_at, last_login_at, last_login_ip, is_active FROM user;

-- Get
SELECT id, name, username, email_enc, password_enc, short_id, created_by, updated_by, created_at, updated_at, last_login_at, last_login_ip, is_active
FROM user
WHERE id = ?;

-- GetPreload
SELECT DISTINCT
    u.id, u.name, u.username, u.email_enc, u.password_enc, u.short_id, u.created_by, u.updated_by, u.created_at, u.updated_at, u.last_login_at, u.last_login_ip, u.is_active,
    r.id AS role_id, r.name AS role_name,
    p.id AS permission_id, p.name AS permission_name
FROM user u
       LEFT JOIN user_role ur ON u.id = ur.user_id
       LEFT JOIN role r ON ur.role_id = r.id
       LEFT JOIN role_permission rp ON r.id = rp.role_id
       LEFT JOIN permission p ON rp.permission_id = p.id
WHERE u.id = ?;

-- Create
INSERT INTO user (id, username, email_enc, name, password_enc, short_id, created_by, updated_by, created_at, updated_at)
VALUES (:id, :username, :email_enc, :name, :password_enc, :short_id, :created_by, :updated_by, :created_at, :updated_at);

-- Update
UPDATE user SET username = :username, email_enc = :email_enc, name = :name, short_id = :short_id, updated_by = :updated_by, updated_at = :updated_at WHERE id = :id;

-- Delete
DELETE FROM user WHERE id = ?;

-- UpdatePassword
UPDATE user
SET password_enc = ?, updated_by = ?, updated_at = ?
WHERE id = ?;

-- GetWithPermissions
SELECT
    u.id, u.name, u.username, u.email_enc, u.password_enc, u.short_id, u.created_by, u.updated_by, u.created_at, u.updated_at, u.last_login_at, u.last_login_ip, u.is_active,
    p.id AS "permissions.id",
    p.name AS "permissions.name",
    p.description AS "permissions.description",
    p.short_id AS "permissions.short_id"
FROM
    user u
LEFT JOIN (
    -- Direct permissions
    SELECT user_id, permission_id FROM user_permission
    UNION
    -- Permissions from GLOBAL roles
    SELECT ur.user_id, rp.permission_id
    FROM user_role ur
    JOIN role_permission rp ON ur.role_id = rp.role_id
    WHERE ur.context_type = '' OR ur.context_type IS NULL
) up ON u.id = up.user_id
LEFT JOIN permission p ON up.permission_id = p.id
WHERE
    u.id = ?;
