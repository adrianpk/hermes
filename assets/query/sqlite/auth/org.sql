-- Res: Org
-- Table: org

-- Create
INSERT INTO org (id, short_id, name, short_description, description, created_by, updated_by, created_at, updated_at)
VALUES (:id, :short_id, :name, :short_description, :description, :created_by, :updated_by, :created_at, :updated_at);

-- GetAll
SELECT id, short_id, name, short_description, description, created_by, updated_by, created_at, updated_at FROM org;

-- Get
SELECT id, short_id, name, short_description, description, created_by, updated_by, created_at, updated_at FROM org WHERE id = ?;

-- GetDefault
SELECT id, short_id, name, short_description, description, created_by, updated_by, created_at, updated_at FROM org ORDER BY created_at ASC LIMIT 1;

-- Update
UPDATE org SET name = :name, short_description = :short_description, description = :description, updated_by = :updated_by, updated_at = :updated_at WHERE id = :id;

-- Delete
DELETE FROM org WHERE id = ?;
