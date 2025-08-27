
-- Update
UPDATE section SET
    name = :name,
    description = :description,
    path = :path,
    layout_id = :layout_id,
    image = :image,
    header = :header,
    updated_by = :updated_by,
    updated_at = :updated_at
WHERE id = :id;
