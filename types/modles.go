package types

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type SecretType string

const (
	SSHKey   SecretType = "ssh_key"
	Password SecretType = "password"
	APIKey   SecretType = "api_key"
)

type PermissionType string

const (
	WriteRead PermissionType = "write"
	ReadOnly  PermissionType = "read"
  NotAllowed PermissionType = "Not Allowed"
)

type User struct {
	UserID       int       `json:"user_id" db:"user_id"`
	Email        string    `json:"email" db:"email"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	LastLogin    time.Time `json:"last_login" db:"last_login"`
}

func (u *User) ValidPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pw)) == nil
}

type Secret struct {
	SecretID      int        `json:"secret_id" db:"secret_id"`
	UserID        int        `json:"user_id" db:"user_id"`
	SecretType    SecretType `json:"secret_type" db:"secret_type"`
	EncryptedData string     `json:"encrypted_data" db:"encrypted_data"`
	Description   string     `json:"description" db:"description"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type Group struct {
	GroupID   int       `json:"group_id" db:"group_id"`
	GroupName string    `json:"group_name" db:"group_name"`
	CreatedBy int       `json:"created_by" db:"created_by"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Role string // member or admin

const (
	Member Role = "member"
	Admin  Role = "admin"
)

type UserGroup struct {
	UserGroupID int  `json:"user_group_id" db:"user_group_id"`
	UserID      int  `json:"user_id" db:"user_id"`
	GroupID     int  `json:"group_id" db:"group_id"`
	Role        Role `json:"role" db:"role"`
}

type SharedSecret struct {
	SharedSecretID  int        `json:"shared_secret_id" db:"shared_secret_id"`
	SecretID        int        `json:"secret_id" db:"secret_id"`
	SharedWithUser  *int       `json:"shared_with_user,omitempty" db:"shared_with_user"`
	SharedWithGroup *int       `json:"shared_with_group,omitempty" db:"shared_with_group"`
	Permissions     PermissionType `json:"permissions" db:"permissions"`
	SharedAt        time.Time  `json:"shared_at" db:"shared_at"`
}

type AuditLog struct {
	LogID     int       `json:"log_id" db:"log_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Action    string    `json:"action" db:"action"`
	Details   string    `json:"details" db:"details"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}