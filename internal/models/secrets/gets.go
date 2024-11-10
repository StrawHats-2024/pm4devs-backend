package secrets

import (
	"context"
	"time"

	"pm4devs.strawhats/internal/xerrors"
)

type SharedSecretUser struct {
	SecretID   int64      `db:"secret_id" json:"secret_id"`   // ID of the secret
	UserID     int64      `db:"user_id" json:"user_id"`       // ID of the user the secret is shared with
	Permission Permission `db:"permission" json:"permission"` // Permission for the shared secret
}

type SharedSecretGroup struct {
	SecretID   int64      `db:"secret_id" json:"secret_id"`   // ID of the secret
	GroupID    int64      `db:"group_id" json:"group_id"`     // ID of the group the secret is shared with
	Permission Permission `db:"permission" json:"permission"` // Permission for the shared secret
}

type FullSharedSecretUserDetail struct {
    SecretID      int64     `json:"secret_id"`
    Name          string    `json:"name"`
    EncryptedData []byte    `json:"encrypted_data"`
    IV            []byte    `json:"iv"`
    OwnerID       int64     `json:"owner_id"`
    UserID        int64     `json:"user_id"`
    Permission    string    `json:"permission"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

func (s *Secrets) GetSecretsSharedToOtherUsers(userID int64) (*[]FullSharedSecretUserDetail, *xerrors.AppError) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    // Prepare the SQL query to select full secret details shared to other users
    query := `
        SELECT s.id AS secret_id, s.name, s.encrypted_data, s.iv, s.owner_id, ssu.user_id, ssu.permission, s.created_at, s.updated_at
        FROM secrets s
        JOIN shared_secrets_user ssu ON ssu.secret_id = s.id
        WHERE s.owner_id = $1;
    `

    // Execute the query
    rows, err := s.DB.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, xerrors.DatabaseError(err, "secrets.GetSecretsSharedToOtherUsers - query execution")
    }
    defer rows.Close()

    // Initialize a slice to hold the retrieved shared secret details
    var sharedSecrets []FullSharedSecretUserDetail

    // Iterate through the rows and scan the data into FullSharedSecretUserDetail structs
    for rows.Next() {
        var sharedSecret FullSharedSecretUserDetail
        if err := rows.Scan(
            &sharedSecret.SecretID,
            &sharedSecret.Name,
            &sharedSecret.EncryptedData,
            &sharedSecret.IV,
            &sharedSecret.OwnerID,
            &sharedSecret.UserID,
            &sharedSecret.Permission,
            &sharedSecret.CreatedAt,
            &sharedSecret.UpdatedAt,
        ); err != nil {
            return nil, xerrors.DatabaseError(err, "secrets.GetSecretsSharedToOtherUsers - scan error")
        }
        sharedSecrets = append(sharedSecrets, sharedSecret)
    }

    // Check for any error that may have occurred during iteration
    if err := rows.Err(); err != nil {
        return nil, xerrors.DatabaseError(err, "secrets.GetSecretsSharedToOtherUsers - rows error")
    }

    // Return the slice of shared secret records with full details
    return &sharedSecrets, nil
}

func (s *Secrets) GetSecretsSharedToGroups(userID int64) (*[]SharedSecretGroup, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the SQL query to select secrets shared to groups
	query := `
		SELECT s.id AS secret_id, ssg.group_id, ssg.permission
		FROM secrets s
		JOIN shared_secrets_group ssg ON ssg.secret_id = s.id
		WHERE s.owner_id = $1;
	`

	// Execute the query
	rows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "secrets.GetSecretsSharedToGroups")
	}
	defer rows.Close()

	// Initialize a slice to hold the retrieved shared secret records
	var sharedSecrets []SharedSecretGroup

	// Iterate through the rows and scan the data into SharedSecretToGroup structs
	for rows.Next() {
		var sharedSecret SharedSecretGroup
		if err := rows.Scan(&sharedSecret.SecretID, &sharedSecret.GroupID, &sharedSecret.Permission); err != nil {
			return nil, xerrors.DatabaseError(err, "secrets.GetSecretsSharedToGroups - scan")
		}
		sharedSecrets = append(sharedSecrets, sharedSecret)
	}

	// Check for any error that may have occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, xerrors.DatabaseError(err, "secrets.GetSecretsSharedToGroups - rows error")
	}

	// Return the slice of shared secret records
	return &sharedSecrets, nil
}

type SharedSecretDetail struct {
	SecretID      int64  `json:"secret_id"`
	Name          string `json:"name"`
	EncryptedData []byte `json:"encrypted_data"`
	IV            []byte `json:"iv"`
	OwnerID       int64  `json:"owner_id"`
	Permission    string `json:"permission"`
}

// GetSecretsSharedWithUser returns a list of secrets, including details, that are shared with the specified user.
func (s *Secrets) GetSecretsSharedWithUser(userID int64) (*[]SharedSecretDetail, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// SQL query to select detailed information for secrets shared with the specified user
	query := `
        SELECT s.id AS secret_id, s.name, s.encrypted_data, s.iv, s.owner_id, ssu.permission
        FROM secrets s
        JOIN shared_secrets_user ssu ON ssu.secret_id = s.id
        WHERE ssu.user_id = $1;
    `

	// Execute the query
	rows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "secrets.GetSecretsSharedWithUser")
	}
	defer rows.Close()

	// Initialize a slice to hold the detailed shared secret records
	var sharedSecrets []SharedSecretDetail

	// Iterate through the rows and scan the data into SharedSecretDetail structs
	for rows.Next() {
		var sharedSecret SharedSecretDetail
		if err := rows.Scan(&sharedSecret.SecretID, &sharedSecret.Name, &sharedSecret.EncryptedData, &sharedSecret.IV, &sharedSecret.OwnerID, &sharedSecret.Permission); err != nil {
			return nil, xerrors.DatabaseError(err, "secrets.GetSecretsSharedWithUser - scan")
		}
		sharedSecrets = append(sharedSecrets, sharedSecret)
	}

	// Check for any error that may have occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, xerrors.DatabaseError(err, "secrets.GetSecretsSharedWithUser - rows error")
	}

	// Return the slice of detailed shared secret records
	return &sharedSecrets, nil
}
