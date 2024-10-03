BEGIN;

-- Revert the default 'activated' field back to false
ALTER TABLE IF EXISTS users ALTER COLUMN activated SET DEFAULT false;

-- Optional: If needed, revert the users updated during the migration
-- This part is optional, as it might not be possible to distinguish users that were activated by other means
-- UPDATE users SET activated = false WHERE created_at < (SELECT created_at FROM users ORDER BY created_at LIMIT 1);

COMMIT;
