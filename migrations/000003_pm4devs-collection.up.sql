BEGIN;

-- Create the groups table
CREATE TABLE IF NOT EXISTS groups (
    id bigserial PRIMARY KEY,
    name text UNIQUE NOT NULL,
    creator_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW()
);

-- Create the group_members table for multiple user relation in groups
CREATE TABLE IF NOT EXISTS group_members (
    group_id bigint NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, user_id)
);

-- Create the secrets table
CREATE TABLE IF NOT EXISTS secrets (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    encrypted_data bytea NOT NULL, -- Using bytea to store encrypted credentials
    owner_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW()
);

-- Create the shared_secrets_group table
CREATE TABLE IF NOT EXISTS shared_secrets_group (
    secret_id bigint NOT NULL REFERENCES secrets(id) ON DELETE CASCADE,
    group_id bigint NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    permission text NOT NULL CHECK (permission IN ('read-only', 'read-write')),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
    PRIMARY KEY (secret_id, group_id)
);

-- Create the shared_secrets_user table
CREATE TABLE IF NOT EXISTS shared_secrets_user (
    secret_id bigint NOT NULL REFERENCES secrets(id) ON DELETE CASCADE,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permission text NOT NULL CHECK (permission IN ('read-only', 'read-write')),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
    PRIMARY KEY (secret_id, user_id)
);

COMMIT;
