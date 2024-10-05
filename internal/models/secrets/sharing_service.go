package secrets

import (
	"context"
	"time"

	"pm4devs.strawhats/internal/xerrors"
)

func (s *Secrets) ShareToUser(secretID, userID int64, permission Permission) *xerrors.AppError {
	// Context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// SQL query to insert a shared secret for a user
	query := `
		INSERT INTO shared_secrets_user (secret_id, user_id, permission, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (secret_id, user_id) DO NOTHING;
	`

	// Execute the query
	_, err := s.DB.ExecContext(ctx, query, secretID, userID, permission)
	if err != nil {
		return xerrors.DatabaseError(err, "secrets.ShareToUser")
	}

	return nil
}

func (s *Secrets) ShareToGroup(secretID, groupID int64, permission Permission) *xerrors.AppError {
	// Context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// SQL query to insert a shared secret for a group
	query := `
		INSERT INTO shared_secrets_group (secret_id, group_id, permission, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (secret_id, group_id) DO NOTHING;
	`

	// Execute the query
	_, err := s.DB.ExecContext(ctx, query, secretID, groupID, permission)
	if err != nil {
		return xerrors.DatabaseError(err, "secrets.ShareToGroup")
	}

	return nil
}

func (s *Secrets) UpdateGroupPermission(secretID, groupID int64, permission Permission) *xerrors.AppError {
	// Context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// SQL query to update the permission for a shared secret in a group
	query := `
		UPDATE shared_secrets_group
		SET permission = $1, updated_at = NOW()
		WHERE secret_id = $2 AND group_id = $3;
	`

	// Execute the update query
	_, err := s.DB.ExecContext(ctx, query, permission, secretID, groupID)
	if err != nil {
		return xerrors.DatabaseError(err, "secrets.UpdateGroupPermission")
	}

	return nil
}

func (s *Secrets) UpdateUserPermission(secretID, userID int64, permission Permission) *xerrors.AppError {
	// Context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// SQL query to update the permission for a shared secret with a user
	query := `
		UPDATE shared_secrets_user
		SET permission = $1, updated_at = NOW()
		WHERE secret_id = $2 AND user_id = $3;
	`

	// Execute the update query
	_, err := s.DB.ExecContext(ctx, query, permission, secretID, userID)
	if err != nil {
		return xerrors.DatabaseError(err, "secrets.UpdateUserPermission")
	}

	return nil
}

func (s *Secrets) RevokeFromGroup(secretID, groupID int64) *xerrors.AppError {
	// Context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// SQL query to delete a shared secret from a group
	query := `
		DELETE FROM shared_secrets_group
		WHERE secret_id = $1 AND group_id = $2;
	`

	// Execute the delete query
	_, err := s.DB.ExecContext(ctx, query, secretID, groupID)
	if err != nil {
		return xerrors.DatabaseError(err, "secrets.RevokeFromGroup")
	}

	return nil
}

func (s *Secrets) RevokeFromUser(secretID, userID int64) *xerrors.AppError {
	// Context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// SQL query to delete a shared secret from a user
	query := `
		DELETE FROM shared_secrets_user
		WHERE secret_id = $1 AND user_id = $2;
	`

	// Execute the delete query
	_, err := s.DB.ExecContext(ctx, query, secretID, userID)
	if err != nil {
		return xerrors.DatabaseError(err, "secrets.RevokeFromUser")
	}

	return nil
}
