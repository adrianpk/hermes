-- Res: Permission
-- Table: permission

-- GetAll
SELECT id, short_id, name, description, created_by, updated_by, created_at, updated_at FROM permission;

-- Get
SELECT id, short_id, name, description, created_by, updated_by, created_at, updated_at
FROM permission
WHERE id = ?;

-- Create
INSERT INTO permission (id, short_id, name, description, created_by, updated_by, created_at, updated_at)
VALUES (:id, :short_id, :name, :description, :created_by, :updated_by, :created_at, :updated_at);

-- Update
UPDATE permission SET short_id = :short_id, name = :name, description = :description, updated_by = :updated_by, updated_at = :updated_at WHERE id = :id;

-- Delete
DELETE FROM permission WHERE id = ?;