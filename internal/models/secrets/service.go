package secrets

import (
	"context"
	"time"

	"pm4devs.strawhats/internal/models/core"
	"pm4devs.strawhats/internal/models/users"
	"pm4devs.strawhats/internal/xerrors"
)

type SecretsRepository interface {
	GetByUserID(id int64) (*[]SecretRecord, *xerrors.AppError)
	GetByUserEmail(email string) (*[]SecretRecord, *xerrors.AppError)
	GetByGroupID(id int64) (*[]SecretRecord, *xerrors.AppError)
	// GetByGroupName(name string) (*[]SecretRecord, *xerrors.AppError)
	NewRecord(name, EncryptedData string, ownerID int64) (*SecretRecord, *xerrors.AppError)
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

func (s *Secrets) NewRecord(name, encryptedData string, ownerID int64) (*SecretRecord, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the SQL query to insert a new secret
	query := `
		INSERT INTO secrets (name, encrypted_data, owner_id, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at;
	`

	// Create a new SecretRecord instance to hold the result
	var secret SecretRecord

	// Execute the insert statement with the provided values
	err := s.DB.QueryRowContext(ctx, query, name, []byte(encryptedData),
		ownerID).Scan(&secret.ID, &secret.CreatedAt)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "secrets.New")
	}

	// Set the Name and EncryptedData fields
	secret.Name = name
	secret.EncryptedData = []byte(encryptedData)

	// Return the newly created secret record
	return &secret, nil
}

func (s *Secrets) GetByUserID(userID int64) (*[]SecretRecord, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the SQL query to select secrets for the given user ID
	query := `
		SELECT id, name, encrypted_data, created_at
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
