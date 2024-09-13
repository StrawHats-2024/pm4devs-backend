package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func NewPostgresStore() (*sql.DB, error) {
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

	return db, nil
}
