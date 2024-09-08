package db

import (
	"database/sql"
	"log"
	"os"
	"pm4devs-backend/pkg/models"

	_ "github.com/lib/pq"
)

type Storage interface {
  CreateUser(*models.User) error
  GetUserById(int) error
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
