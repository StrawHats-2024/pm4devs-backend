package secrets

import (
	"context"
	"database/sql"
	"time"

	"pm4devs.strawhats/internal/xerrors"
)

// GetUserSecretPermission retrieves the permission of a user for a given secret
func (s *Secrets) GetUserSecretPermission(userID int64, secretID int64) (
	Permission,
	*xerrors.AppError,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if currSecret, err := s.GetSecretByID(secretID); err != nil {
		return NOTALLOWED, err
	} else if currSecret.ID == userID {
		return ReadWrite, nil
	}

	// Check for direct user permission
	directPermissionQuery := `
		SELECT permission
		FROM shared_secrets_user
		WHERE secret_id = $1 AND user_id = $2;
	`

	var permission string
	err := s.DB.QueryRowContext(ctx, directPermissionQuery, secretID, userID).Scan(&permission)

	if err == nil {
		// Direct permission found
		if permission == "read-only" {
			return ReadOnly, nil
		} else if permission == "read-write" {
			return ReadWrite, nil
		}
	} else if err != sql.ErrNoRows {
		return NOTALLOWED, xerrors.DatabaseError(err, "secrets.GetUserSecretPermission (direct permission check)")
	}

	// Check if user is part of a group that has permission to the secret
	groupPermissionQuery := `
		SELECT sg.permission
		FROM shared_secrets_group sg
		JOIN group_members gm ON gm.group_id = sg.group_id
		WHERE sg.secret_id = $1 AND gm.user_id = $2
		LIMIT 1;
	`

	err = s.DB.QueryRowContext(ctx, groupPermissionQuery, secretID, userID).Scan(&permission)

	if err == nil {
		// Group permission found
		if permission == "read-only" {
			return ReadOnly, nil
		} else if permission == "read-write" {
			return ReadWrite, nil
		}
	} else if err != sql.ErrNoRows {
		return NOTALLOWED, xerrors.DatabaseError(err, "secrets.GetUserSecretPermission (group permission check)")
	}

	// If no permissions found
	return NOTALLOWED, nil
}
