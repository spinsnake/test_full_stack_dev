DELETE image_tags
FROM image_tags
INNER JOIN images
  ON images.id = image_tags.image_id
WHERE images.source = 'seed:placehold.co'
  AND images.alt_text LIKE 'Demo Image %';

DELETE FROM images
WHERE source = 'seed:placehold.co'
  AND alt_text LIKE 'Demo Image %';

DELETE FROM tags
WHERE slug LIKE 'demo-tag-%'
  AND name LIKE 'Demo Tag %'
  AND NOT EXISTS (
    SELECT 1
    FROM image_tags
    WHERE image_tags.tag_id = tags.id
  );
