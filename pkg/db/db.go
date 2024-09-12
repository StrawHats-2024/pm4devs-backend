package db

import (
	"database/sql"
	"log"
	"os"
	"pm4devs-backend/pkg/models"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*models.User) (int, error)
	GetUserById(int) (*models.User, error)
	GetUserByEmail(string) (*models.User, error)
	GetAllUsers() ([]*models.User, error) // passwordHash will be set empty
	UpdateLastLogin(int) error

	CreateSecret(*models.Secret) (int, error)
	GetAllSecret() ([]*models.Secret, error)
	GetSecretById(int) (*models.Secret, error)
	GetAllSecretsByUserID(int) ([]*models.Secret, error)
	DeleteSecretById(int) error
  UpdateSecretById(int, UpdateSecretReq) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("Database URL not found in env")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}
