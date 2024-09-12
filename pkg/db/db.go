package db

import (
	"database/sql"
	"log"
	"os"
	"pm4devs-backend/pkg/models"

	_ "github.com/lib/pq"
)

// Storage defines the interface for data storage operations related to users, secrets, and groups.
type Storage interface {
	// CreateUser creates a new user in the database and returns the user's ID.
	CreateUser(*models.User) (int, error)

	// GetUserById retrieves a user by their ID.
	GetUserById(int) (*models.User, error)

	// GetUserByEmail retrieves a user by their email address.
	GetUserByEmail(string) (*models.User, error)

	// GetAllUsers retrieves all users from the database.
	// The passwordHash for each user will be set to an empty string for security.
	GetAllUsers() ([]*models.User, error)

	// UpdateLastLogin updates the last login timestamp for a user identified by their ID.
	UpdateLastLogin(int) error

	// CreateSecret creates a new secret and returns its ID.
	CreateSecret(*models.Secret) (int, error)

	// GetAllSecret retrieves all secrets from the database.
	GetAllSecret() ([]*models.Secret, error)

	// GetSecretById retrieves a secret by its ID.
	GetSecretById(int) (*models.Secret, error)

	// GetAllSecretsByUserID retrieves all secrets associated with a user identified by their ID.
	GetAllSecretsByUserID(int) ([]*models.Secret, error)

	// DeleteSecretById deletes a secret identified by its ID from the database.
	DeleteSecretById(int) error

	// UpdateSecretById updates a secret identified by its ID using the provided request data.
	UpdateSecretById(int, UpdateSecretReq) error

	// CreateGroup creates a new group and returns the group's ID.
	CreateGroup(*models.Group) (int, error)

	// GetUserGroups retrieves a list of groups that a user belongs to, identified by their user ID.
	GetUserGroups(userId int) ([]*GetUserGroupRes, error) // Returns list of groups user belongs to

	// GetGroupById retrieves a group by its ID.
	GetGroupById(groupId int) (*GroupWithUsers, error)

	// AddUserToGroup adds a user to a specified group.
	AddUserToGroup(groupId int, newUser AddUserToGroupReq) error

	// DeleteUserFromGroup removes a user from a specified group using their email address.
	DeleteUserFromGroup(groupId int, email string) error

	// GetUsersByGroupId retrieves all users belonging to a group identified by its ID.
	GetUsersByGroupId(groupId int) ([]UserRes, error) // Returns all users of the group with groupId

	// DeleteGroup deletes a group identified by its ID from the database.
	DeleteGroup(groupId int) error // groupId to delete group

	// UpdateGroupName updates the name of a specified group identified by its ID.
	UpdateGroupName(groupId int, newName string) error

	// IsUserAdminInGroup checks if the user is an admin in the specified group.
	IsUserAdminInGroup(userId, groupId int) (bool, error)

	// IsGroupCreator checks if the user is the creator of the specified group.
	IsGroupCreator(userId, groupId int) (bool, error)

	// Get user role with user_email and groupId
	GetUserRoleInGroup(string, int) (models.Role, error)
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
