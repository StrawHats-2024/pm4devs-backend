package secrets

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"pm4devs.strawhats/internal/models/core"
	"pm4devs.strawhats/internal/models/users"
	"pm4devs.strawhats/internal/xerrors"
)

type Permission string

const (
	ReadOnly   Permission = "read-only"
	ReadWrite  Permission = "read-write"
	NOTALLOWED Permission = ""
)

type SecretsRepository interface {
	GetByUserID(id int64) (*[]SecretRecord, *xerrors.AppError)
	GetByUserEmail(email string) (*[]SecretRecord, *xerrors.AppError)
	GetByGroupID(id int64) (*[]SecretRecord, *xerrors.AppError)
	// GetByGroupName(name string) (*[]SecretRecord, *xerrors.AppError)
	NewRecord(name, EncryptedData, IV string, ownerID int64) (*SecretRecord, *xerrors.AppError)
	Delete(secretID int64) *xerrors.AppError
	Update(secretID int64, newName, newEncryptedData string) *xerrors.AppError
	GetSecretByID(secretID int64) (*SecretRecord, *xerrors.AppError)
	ShareToGroup(secretID, groupID int64, permission Permission) *xerrors.AppError
	ShareToUser(secretID, userID int64, permission Permission) *xerrors.AppError
	UpdateGroupPermission(secretID, groupID int64, permission Permission) *xerrors.AppError
	UpdateUserPermission(secretID, userID int64, permission Permission) *xerrors.AppError
	RevokeFromGroup(secretID, groupID int64) *xerrors.AppError
	RevokeFromUser(secretID, userID int64) *xerrors.AppError
	GetUserSecretPermission(userID int64, secretID int64) (Permission, *xerrors.AppError)
	GetSecretsSharedToOtherUsers(userID int64) (*[]SharedSecretUser, *xerrors.AppError)
	GetSecretsSharedToGroups(userID int64) (*[]SharedSecretGroup, *xerrors.AppError)
}

type Secrets struct {
	DB core.Queryable
}

func Repository(db core.Queryable) SecretsRepository {
	return &Secrets{DB: db}
}

func (s *Secrets) GetByGroupID(id int64) (*[]SecretRecord, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// SQL query to get secrets shared with the group
	query := `
		SELECT secrets.id, secrets.name, secrets.encrypted_data, secrets.created_at
		FROM secrets
		INNER JOIN shared_secrets_group ON shared_secrets_group.secret_id = secrets.id
		WHERE shared_secrets_group.group_id = $1;
	`

	// Slice to hold the results
	var secrets []SecretRecord

	// Execute the query and iterate over the rows
	rows, err := s.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "secrets.GetByGroupID")
	}
	defer rows.Close()

	// Loop through the rows and scan the data into the SecretRecord slice
	for rows.Next() {
		var secret SecretRecord
		if err := rows.Scan(&secret.ID, &secret.Name, &secret.EncryptedData, &secret.CreatedAt); err != nil {
			return nil, xerrors.DatabaseError(err, "secrets.GetByGroupID.Scan")
		}
		secrets = append(secrets, secret)
	}

	// Check for any error that might have occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, xerrors.DatabaseError(err, "secrets.GetByGroupID.Rows")
	}

	return &secrets, nil
}

// func (s *Secrets) GetByGroupName(name string) (*[]SecretRecord, *xerrors.AppError) {
//
// }

func (s *Secrets) GetByUserEmail(email string) (*[]SecretRecord, *xerrors.AppError) {
	userRepo := users.Users{DB: s.DB}
	user, err := userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return s.GetByUserID(user.ID)
}

func (s *Secrets) NewRecord(name, encryptedData, iv string, ownerID int64) (*SecretRecord, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the SQL query to insert a new secret
	query := `
		INSERT INTO secrets (name, encrypted_data, iv, owner_id, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, created_at;
	`

	// Create a new SecretRecord instance to hold the result
	var secret SecretRecord

	// Execute the insert statement with the provided values
	err := s.DB.QueryRowContext(ctx, query, name, []byte(encryptedData), []byte(iv),
		ownerID).Scan(&secret.ID, &secret.CreatedAt)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "secrets.New")
	}

	// Set the Name and EncryptedData fields
	secret.Name = name
	secret.EncryptedData = []byte(encryptedData)
	secret.IV = []byte(encryptedData)

	// Return the newly created secret record
	return &secret, nil
}

func (s *Secrets) GetByUserID(userID int64) (*[]SecretRecord, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the SQL query to select secrets for the given user ID
	query := `
		SELECT id, name, encrypted_data, iv, created_at
		FROM secrets
		WHERE owner_id = $1;
	`

	// Execute the query
	rows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "secrets.GetByUserID")
	}
	defer rows.Close()

	// Initialize a slice to hold the retrieved secret records
	var secrets []SecretRecord

	// Iterate through the rows and scan the data into SecretRecord structs
	for rows.Next() {
		var secret SecretRecord
		if err := rows.Scan(&secret.ID, &secret.Name, &secret.EncryptedData, &secret.CreatedAt); err != nil {
			return nil, xerrors.DatabaseError(err, "secrets.GetByUserID - scan")
		}
		secrets = append(secrets, secret)
	}

	// Check for any error that may have occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, xerrors.DatabaseError(err, "secrets.GetByUserID - rows error")
	}

	// Return the slice of secret records
	return &secrets, nil
}

// Delete a secret by secret ID and owner ID
func (s *Secrets) Delete(secretID int64) *xerrors.AppError {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		DELETE FROM secrets
		WHERE id = $1;
	`

	result, err := s.DB.ExecContext(ctx, query, secretID)
	if err != nil {
		return xerrors.DatabaseError(err, "secrets.Delete")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return xerrors.DatabaseError(err, "secrets.Delete.RowsAffected")
	}

	if rowsAffected == 0 {
		return xerrors.DatabaseError(fmt.Errorf("No secret found with id: %d", secretID),
			"secrets.Delete")
	}

	return nil
}

// Update a secret by secret ID
func (s *Secrets) Update(secretID int64, newName, newEncryptedData string) *xerrors.AppError {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE secrets
		SET name = $1, encrypted_data = $2, updated_at = NOW()
		WHERE id = $3;
	`

	_, err := s.DB.ExecContext(ctx, query, newName, []byte(newEncryptedData), secretID)
	if err != nil {
		return xerrors.DatabaseError(err, "secrets.Update")
	}

	return nil
}

func (s *Secrets) GetSecretByID(secretID int64) (*SecretRecord, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the SQL query to get the secret by its ID
	query := `
		SELECT id, name, encrypted_data, owner_id, created_at
		FROM secrets
		WHERE id = $1;
	`

	// Create a SecretRecord instance to hold the result
	var secret SecretRecord

	// Execute the query and scan the result into the secret struct
	err := s.DB.QueryRowContext(ctx, query, secretID).Scan(
		&secret.ID, &secret.Name, &secret.EncryptedData, &secret.OwnerID, &secret.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, xerrors.ClientError(http.StatusNotFound, fmt.Sprintf("No secret found with id: %d", secretID), "secrets.GetSecretByID", fmt.Errorf("Secret not found with id: %d", secretID))
		}
		return nil, xerrors.DatabaseError(err, "secrets.GetSecretByID")
	}

	return &secret, nil
}
