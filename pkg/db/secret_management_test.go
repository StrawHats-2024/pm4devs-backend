package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
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
	db := setupPostgresTestDB(t)
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
	db := setupPostgresTestDB(t)
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
	db := setupPostgresTestDB(t)
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
	db := setupPostgresTestDB(t)
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
	db := setupPostgresTestDB(t)
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

func setupPostgresTestDB(t *testing.T) *sql.DB {
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

	sqlContent, err := os.ReadFile("../../db/seed.sql")
	if err != nil {
		t.Fatalf("Failed to read sql file: %v", err)
	}
	_, err = db.Exec(string(sqlContent))
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
