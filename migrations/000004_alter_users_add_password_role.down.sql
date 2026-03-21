ALTER TABLE users
  DROP COLUMN password,
  DROP COLUMN role,
  ALTER COLUMN spotify_id SET NOT NULL;
