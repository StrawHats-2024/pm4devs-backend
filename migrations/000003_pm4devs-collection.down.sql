BEGIN;

-- Drop the shared_secrets_user table
DROP TABLE IF EXISTS shared_secrets_user;

-- Drop the shared_secrets_group table
DROP TABLE IF EXISTS shared_secrets_group;

-- Drop the secrets table
DROP TABLE IF EXISTS secrets;

-- Drop the group_members table
DROP TABLE IF EXISTS group_members;

-- Drop the groups table
DROP TABLE IF EXISTS groups;

COMMIT;
