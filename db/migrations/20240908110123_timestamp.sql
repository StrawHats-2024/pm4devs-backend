
-- +goose Up
-- +goose StatementBegin
CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Secret (
    secret_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    secret_type VARCHAR(50) CHECK (secret_type IN ('password', 'ssh_key', 'api_key')),
    encrypted_data TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);

CREATE TABLE Groups (
    group_id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    created_by INT REFERENCES Users(user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE UserGroup (
    user_group_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    group_id INT REFERENCES Groups(group_id) ON DELETE CASCADE,
    role VARCHAR(50) CHECK (role IN ('admin', 'member'))
);

CREATE TABLE SharedSecret (
    shared_secret_id SERIAL PRIMARY KEY,
    secret_id INT REFERENCES Secret(secret_id) ON DELETE CASCADE,
    shared_with_user INT REFERENCES Users(user_id) ON DELETE SET NULL,
    shared_with_group INT REFERENCES Groups(group_id) ON DELETE SET NULL,
    permissions VARCHAR(50) CHECK (permissions IN ('read', 'write')),
    shared_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK (shared_with_user IS NOT NULL OR shared_with_group IS NOT NULL) -- Ensure at least one is not null
);

CREATE TABLE AuditLog (
    log_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(user_id),
    action VARCHAR(255) NOT NULL,
    details TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS SharedSecret;
DROP TABLE IF EXISTS UserGroup;
DROP TABLE IF EXISTS AuditLog;
DROP TABLE IF EXISTS Secret;
DROP TABLE IF EXISTS Groups;
DROP TABLE IF EXISTS Users;
-- +goose StatementEnd
