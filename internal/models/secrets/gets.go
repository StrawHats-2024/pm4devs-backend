package secrets

import (
	"context"
	"time"

	"pm4devs.strawhats/internal/xerrors"
)

type SharedSecretUser struct {
	SecretID   int64      `db:"secret_id"`  // ID of the secret
	UserID     int64      `db:"user_id"`    // ID of the user the secret is shared with
	Permission Permission `db:"permission"` // Permission for the shared secret
}

type SharedSecretGroup struct {
	SecretID   int64      `db:"secret_id"`  // ID of the secret
	GroupID    int64      `db:"group_id"`   // ID of the group the secret is shared with
	Permission Permission `db:"permission"` // Permission for the shared secret
}

// returns list of secrets shared to other users by authenticated users
func (s *Secrets) GetSecretsSharedToOtherUsers(userID int64) (*[]SharedSecretUser, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the SQL query to select secrets shared to other users
	query := `
		SELECT s.id AS secret_id, ssu.user_id, ssu.permission
		FROM secrets s
		JOIN shared_secrets_user ssu ON ssu.secret_id = s.id
		WHERE s.owner_id = $1;
	`

	// Execute the query
	rows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "secrets.GetSecretsSharedToUsers")
	}
	defer rows.Close()

	// Initialize a slice to hold the retrieved shared secret records
	var sharedSecrets []SharedSecretUser

	// Iterate through the rows and scan the data into SharedSecretToUser structs
	for rows.Next() {
		var sharedSecret SharedSecretUser
		if err := rows.Scan(&sharedSecret.SecretID, &sharedSecret.UserID, &sharedSecret.Permission); err != nil {
			return nil, xerrors.DatabaseError(err, "secrets.GetSecretsSharedToUsers - scan")
		}
		sharedSecrets = append(sharedSecrets, sharedSecret)
	}

	// Check for any error that may have occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, xerrors.DatabaseError(err, "secrets.GetSecretsSharedToUsers - rows error")
	}

	// Return the slice of shared secret records
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
