package db

import (
	"context"
	"database/sql"
	"fmt"
	"pm4devs-backend/pkg/models"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *sql.DB

func TestCreateSecret(t *testing.T) {
	db := setupPostgres(t)
	defer teardownPostgres(t, db)

	pg := &PostgresStore{db: db}
	userid := 1
	secret := &models.Secret{
		UserID:        userid,
		SecretType:    "password",
		EncryptedData: "encrypted_password",
		Description:   "Test password",
	}

	secretID, err := pg.CreateSecret(secret)
	if err != nil {
		t.Fatalf("Failed to create secret: %v", err)
	}

	if secretID == 0 {
		t.Errorf("Expected secret ID to be greater than 0, got %d", secretID)
	}
}

func TestGetAllSecretsByUserID(t *testing.T) {
	db := setupPostgres(t)
	defer teardownPostgres(t, db)

	pg := &PostgresStore{db: db}

	userID := 1 // Assuming user ID 1 exists in test DB
	secrets, err := pg.GetAllSecretsByUserID(userID)
	if err != nil {
		t.Fatalf("Failed to get secrets by user ID: %v", err)
	}

	if len(secrets) == 0 {
		t.Errorf("Expected to get some secrets for user ID %d, got 0", userID)
	}
}

func TestGetAllSecret(t *testing.T) {
	db := setupPostgres(t)
	defer teardownPostgres(t, db)

	pg := &PostgresStore{db: db}
	secrets, err := pg.GetAllSecret()
	if err != nil {
		t.Fatalf("Failed to get all secrets: %v", err)
	}

	if len(secrets) == 0 {
		t.Errorf("Expected to get some secrets, got 0")
	}
}

func TestGetSecretById(t *testing.T) {
	db := setupPostgres(t)
	defer teardownPostgres(t, db)

	pg := &PostgresStore{db: db}

	secretID := 1 // Assuming a secret with this ID exists in test DB
	secret, err := pg.GetSecretById(secretID)
	if err != nil {
		t.Fatalf("Failed to get secret by ID: %v", err)
	}

	if secret.SecretID != secretID {
		t.Errorf("Expected secret ID to be %d, got %d", secretID, secret.SecretID)
	}
}

func TestDeleteSecretById(t *testing.T) {
	db := setupPostgres(t)
	defer teardownPostgres(t, db)

	pg := &PostgresStore{db: db}
	secretID := 1 // Assuming secret with this ID exists in test DB
	err := pg.DeleteSecretById(secretID)
	if err != nil {
		t.Fatalf("Failed to delete secret by ID: %v", err)
	}

	// Try to fetch the deleted secret
	secret, err := pg.GetSecretById(secretID)
	if secret != nil || err == nil {
		t.Errorf("Expected secret to be deleted, but found it")
	}
}

func setupPostgres(t *testing.T) *sql.DB {
	username := "testuser"
	password := "testpassword"
	dbname := "testdb"
	ctx := context.Background()
	t.Log("Starting setup")
	// Use the new Run method for Postgres container
	pgContainer, err := postgres.Run(ctx,
		"postgres:latest",
		postgres.WithDatabase(dbname),
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	// Get the host and port for connecting to the container
	host, err := pgContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	// Build the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port.Port(), username, password, dbname)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	_, err = db.Exec(`
DROP TABLE IF EXISTS SharedSecret;
DROP TABLE IF EXISTS UserGroup;
DROP TABLE IF EXISTS AuditLog;
DROP TABLE IF EXISTS Secret;
DROP TABLE IF EXISTS Groups;
DROP TABLE IF EXISTS Users;

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

-- Insert dummy data
INSERT INTO Users (email, username, password_hash, created_at, last_login) VALUES
('user1@example.com', 'user1', 'hash1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('user2@example.com', 'user2', 'hash2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('user3@example.com', 'user3', 'hash3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO Secret (user_id, secret_type, encrypted_data, description, created_at, updated_at) VALUES
(1, 'password', 'encrypted_pass1', 'Password for service X', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, 'ssh_key', 'encrypted_ssh_key1', 'SSH key for server Y', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(3, 'api_key', 'encrypted_api_key1', 'API key for service Z', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO Groups (group_name, created_by, created_at) VALUES
('Admins', 1, CURRENT_TIMESTAMP),
('Developers', 2, CURRENT_TIMESTAMP),
('QA', 3, CURRENT_TIMESTAMP);

INSERT INTO UserGroup (user_id, group_id, role) VALUES
(1, 1, 'admin'),
(2, 2, 'member'),
(3, 3, 'member');

INSERT INTO SharedSecret (secret_id, shared_with_user, permissions, shared_at) VALUES
(1, 2, 'read', CURRENT_TIMESTAMP),
(2, 3, 'write', CURRENT_TIMESTAMP);

INSERT INTO AuditLog (user_id, action, details, timestamp) VALUES
(1, 'Login', 'User logged in', CURRENT_TIMESTAMP),
(2, 'Create Secret', 'User created a secret', CURRENT_TIMESTAMP),
(3, 'Update Secret', 'User updated a secret', CURRENT_TIMESTAMP);

	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	return db
}

func teardownPostgres(_ *testing.T, db *sql.DB) {
	// Close the database connection
	if db != nil {
		_ = db.Close()
	}
}
