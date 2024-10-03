BEGIN;

-- Alter the users table to set default 'activated' to true for new users
ALTER TABLE users ALTER COLUMN activated SET DEFAULT true;

-- Update all existing users' activated field to true
UPDATE users SET activated = true WHERE activated = false;

COMMIT;
