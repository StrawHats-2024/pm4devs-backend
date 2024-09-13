package user

import (
	"database/sql"
	"fmt"
	"pm4devs-backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user *types.User) (int, error) {
	query := `INSERT INTO Users (email, username, password_hash, created_at) VALUES ($1, $2, $3, $4) RETURNING user_id`

	// Use QueryRow to get the newly created user_id
	var userId int
	err := s.db.QueryRow(
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

// GetUserById retrieves a user by their ID.
func (s *Store) GetUserById(userId int) (*types.User, error) {
	query := `
  SELECT (user_id, email, username, password_hash, created_at, last_login)
  FROM Users
  WHERE user_id = $1;
  `
	user := new(types.User)
	err := s.db.QueryRow(query, userId).Scan(
		&user.UserID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.LastLogin,
	)
	if err != nil {
		fmt.Println("err: ", err)
		return nil, err
	}
	return user, nil
}

// GetUserByEmail retrieves a user by their email address.
func (s *Store) GetUserByEmail(userEmail string) (*types.User, error) {
	query := `
  SELECT *
  FROM Users
  WHERE email = $1;
  `
	// Initialize a User struct to hold the result
	user := new(types.User)

	// Use QueryRow to fetch the single row based on email
	err := s.db.QueryRow(query, userEmail).Scan(
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
			return nil, fmt.Errorf("user with email %s not found", userEmail)
		}
		return nil, err
	}

	return user, nil
}

// GetAllUsers retrieves all users from the database.
// The passwordHash for each user will be set to an empty string for security.
func (s *Store) GetAllUsers() ([]*types.User, error) {
	query := `
  SELECT user_id, email, created_at, last_login
  FROM Users;
  `
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*types.User

	for rows.Next() {
		user := new(types.User)
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

// UpdateLastLogin updates the last login timestamp for a user identified by their ID.
func (s *Store) UpdateLastLogin(userId int) error {
	query := `
    UPDATE Users
    SET last_login = CURRENT_TIMESTAMP
    WHERE user_id = $1;
    `
	_, err := s.db.Exec(query, userId)
	return err
}
