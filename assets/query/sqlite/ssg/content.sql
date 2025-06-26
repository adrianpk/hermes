-- Res: Content
-- Table: content

-- Create
INSERT INTO content (
    id, user_id, heading, body, status, short_id, created_by, updated_by, created_at, updated_at
) VALUES (
    :id, :user_id, :heading, :body, :status, :short_id, :created_by, :updated_by, :created_at, :updated_at
);

-- GetAll
SELECT * FROM content;

-- Get
SELECT * FROM content WHERE id = :id;

