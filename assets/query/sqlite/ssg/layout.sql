-- Res: Layout
-- Table: layout

-- Create
INSERT INTO layout (
    id, short_id, name, description, code, created_by, updated_by, created_at, updated_at
) VALUES (
    :id, :shortID, :name, :description, :code, :created_by, :updated_by, :created_at, :updated_at
);

-- GetAll
SELECT * FROM layout;

-- Get
SELECT * FROM layout WHERE id = :id;

-- Update
UPDATE layout SET
    name = :name,
    description = :description,
    code = :code,
    updated_by = :updated_by,
    updated_at = :updated_at
WHERE id = :id;

-- Delete
DELETE FROM layout WHERE id = :id;
