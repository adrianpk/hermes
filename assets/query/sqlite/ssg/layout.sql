-- Res: Layout
-- Table: layout

-- Create
INSERT INTO layout (
    id, short_id, name, description, code, created_by, updated_by, created_at, updated_at
) VALUES (
    :id, :short_id, :name, :description, :code, :created_by, :updated_by, :created_at, :updated_at
);

-- GetAll
SELECT * FROM layout;

-- Get
SELECT * FROM layout WHERE id = :id;
