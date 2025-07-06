-- Res: Section
-- Table: section

-- Create
INSERT INTO section (
    id, short_id, name, description, path, layout_id, image, header, created_by, updated_by, created_at, updated_at
) VALUES (
    :id, :short_id, :name, :description, :path, :layout_id, :image, :header, :created_by, :updated_by, :created_at, :updated_at
);

-- GetAll
SELECT * FROM section;

-- Get
SELECT * FROM section WHERE id = :id;
