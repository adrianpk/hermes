-- Res: Team
-- Table: team

-- Create
INSERT INTO team (id, org_id, short_id, name, short_description, description, created_by, updated_by, created_at, updated_at)
VALUES (:id, :org_id, :short_id, :name, :short_description, :description, :created_by, :updated_by, :created_at, :updated_at);

-- GetAll
SELECT id, org_id, short_id, name, short_description, description, created_by, updated_by, created_at, updated_at FROM team WHERE org_id = ?;

-- Get
SELECT
    id,
    org_id,
    short_id,
    name,
    short_description,
    description,
    created_by,
    updated_by,
    created_at,
    updated_at
FROM team WHERE id = ?;

-- Update
UPDATE team SET short_id = :short_id, name = :name, short_description = :short_description, description = :description, updated_by = :updated_by, updated_at = :updated_at WHERE id = :id;

-- Delete
DELETE FROM team WHERE id = ?;