package db

import (
	"database/sql"
	"fmt"
	"pm4devs-backend/pkg/models"

	_ "github.com/lib/pq"
)

func (pg *PostgresStore) CreateUser(user *models.User) (int, error) {

	query := `INSERT INTO Users (email, username, password_hash, created_at) VALUES ($1, $2, $3, $4) RETURNING user_id`

	// Use QueryRow to get the newly created user_id
	var userId int
	err := pg.db.QueryRow(
		query,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.CreatedAt,
	).Scan(&userId)

	if err != nil {
		return 0, err
	}

	fmt.Println("New User ID: ", user.UserID)
	return userId, nil
}
func (pg *PostgresStore) GetUserById(id int) (*models.User, error) {
	query := `
  SELECT (user_id, email, username, password_hash, created_at, last_login)
  FROM Users
  WHERE user_id = $1;
  `
	user := new(models.User)
	err := pg.db.QueryRow(query, id).Scan(
		&user.UserID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.LastLogin,
	)
	if err != nil {
    fmt.Println("err: ", err);
		return nil, err
	}
	return user, nil
}

func (pg *PostgresStore) GetAllUsers() ([]*models.User, error) {
	query := `
  SELECT user_id, email, created_at, last_login
  FROM Users;
  `
	rows, err := pg.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*models.User

	for rows.Next() {
		user := new(models.User)
		err := rows.Scan(
			&user.UserID,
			&user.Email,
			&user.CreatedAt,
			&user.LastLogin,
		)
		users = append(users, user)
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}

func (pg *PostgresStore) GetUserByEmail(email string) (*models.User, error) {
	query := `
  SELECT *
  FROM Users
  WHERE email = $1;
  `
	// Initialize a User struct to hold the result
	user := new(models.User)

	// Use QueryRow to fetch the single row based on email
	err := pg.db.QueryRow(query, email).Scan(
		&user.UserID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.LastLogin,
	)

	// Handle the case where no row is found
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, err
	}

	return user, nil
}
func (pg *PostgresStore) UpdateLastLogin(userID int) error {
	query := `
    UPDATE Users
    SET last_login = CURRENT_TIMESTAMP
    WHERE user_id = $1;
    `
	_, err := pg.db.Exec(query, userID)
	return err
}
