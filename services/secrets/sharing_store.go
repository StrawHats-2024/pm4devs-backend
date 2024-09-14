package secrets

import (
	"database/sql"
	"errors"
	"fmt"
	"pm4devs-backend/types"
)

// ShareSecretWithUser shares a secret with a user and returns an error if any.
func (pg *Store) ShareSecretWithUser(secretID int, userEmail string, permission types.PermissionType) error {
	// Check if the user exists
	var userID int
	err := pg.db.QueryRow("SELECT user_id FROM Users WHERE email = $1", userEmail).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	// Insert the shared secret record
	_, err = pg.db.Exec(`
		INSERT INTO SharedSecret (secret_id, shared_with_user, permissions) 
		VALUES ($1, $2, $3)`,
		secretID, userID, permission,
	)
	return err
}

// ShareSecretWithGroup shares a secret with a group and returns an error if any.
func (pg *Store) ShareSecretWithGroup(secretID int, groupID int, permission types.PermissionType) error {
	// Insert the shared secret record
	_, err := pg.db.Exec(`
		INSERT INTO SharedSecret (secret_id, shared_with_group, permissions) 
		VALUES ($1, $2, $3)`,
		secretID, groupID, permission,
	)
	return err
}

// GetAllSharedSecrets retrieves all secrets shared with a specific user by their ID.
func (pg *Store) GetAllSharedSecrets(userID int) ([]*types.Secret, error) {
	rows, err := pg.db.Query(`
		SELECT s.secret_id, s.user_id, s.secret_type, s.encrypted_data, s.description, s.created_at, s.updated_at
		FROM SharedSecret ss
		JOIN Secret s ON ss.secret_id = s.secret_id
		WHERE ss.shared_with_user = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []*types.Secret
	for rows.Next() {
		secret := &types.Secret{}
		if err := rows.Scan(&secret.SecretID, &secret.UserID, &secret.SecretType, &secret.EncryptedData, &secret.Description, &secret.CreatedAt, &secret.UpdatedAt); err != nil {
			return nil, err
		}
		secrets = append(secrets, secret)
	}

	return secrets, nil
}

// RevokeSharingFromUser revokes sharing of a secret from a user identified by their email.
func (pg *Store) RevokeSharingFromUser(secretID int, userEmail string) error {
	// Check if the user exists
	var userID int
	err := pg.db.QueryRow("SELECT user_id FROM Users WHERE email = $1", userEmail).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	// Delete the shared secret record
	result, err := pg.db.Exec(`
		DELETE FROM SharedSecret 
		WHERE secret_id = $1 AND shared_with_user = $2`,
		secretID, userID,
	)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("No user found with access")
	}
	return nil
}

// RevokeSharingFromGroup revokes sharing of a secret from a group.
func (pg *Store) RevokeSharingFromGroup(secretID int, groupID int) error {
	// Delete the shared secret record
	_, err := pg.db.Exec(`
		DELETE FROM SharedSecret 
		WHERE secret_id = $1 AND shared_with_group = $2`,
		secretID, groupID,
	)
	return err
}

// GetUserSecretPermission retrieves the permission of a user for a given secret.
func (pg *Store) GetUserSecretPermission(userID int, secretID int) (types.PermissionType, error) {
	// First, check if the user is the creator of the secret
	var creatorID int
	err := pg.db.QueryRow(`SELECT user_id FROM Secret WHERE secret_id = $1`, secretID).Scan(&creatorID)
	if err != nil {
		return "", fmt.Errorf("failed to find secret: %w", err)
	}

	// If the user is the creator, they have write access
	if userID == creatorID {
		return types.WriteRead, nil
	}

	// Check if the secret has been shared with the user
	var permission types.PermissionType
	err = pg.db.QueryRow(`SELECT permissions FROM SharedSecret WHERE secret_id = $1 AND shared_with_user = $2`, secretID, userID).Scan(&permission)
	if err != nil {
		return types.NotAllowed, nil
	}

	// If permission is found, return it
	return permission, nil
}
