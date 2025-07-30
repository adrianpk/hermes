-- Res: Content
-- Table: content

-- Create
INSERT INTO content (
    id, short_id, user_id, heading, body, status, created_by, updated_by, created_at, updated_at
) VALUES (
    :id, :short_id, :user_id, :heading, :body, :status, :created_by, :updated_by, :created_at, :updated_at
);

-- GetAll
SELECT * FROM content;

-- Get
SELECT * FROM content WHERE id = :id;

-- Update
UPDATE content SET
    heading = :heading,
    body = :body,
    status = :status,
    updated_by = :updated_by,
    updated_at = :updated_at
WHERE id = :id;