package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"pm4devs-backend/pkg/models"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(models.User) error
	GetUserById(int) (*models.User, error)
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

func (pg *PostgresStore) CreateUser(user models.User) error {
  fmt.Println("user: ", user);
	query := `INSERT INTO Users (email, password_hash, last_login) VALUES ($1, $2, $3)`
	row, err := pg.db.Query(
		query,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
	)
	newUser, err := scanIntoUser(row)
	fmt.Println("newUser: ", newUser)
	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresStore) GetUserById(id int) (*models.User, error) {
	query := `
  SELECT email, created_at
  FROM Users
  WHERE user_id = $1;
  `
	row, err := pg.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	user, err := scanIntoUser(row)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func scanIntoUser(rows *sql.Rows) (*models.User, error) {
	user := new(models.User)
	err := rows.Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.LastLogin,
	)

	return user, err
}
