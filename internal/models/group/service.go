package group

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

type GroupRepository interface {
	GetByGroupID(id int64) (*GroupRecordWithUsers, *xerrors.AppError)
	GetGroupUsers(name string) (*GroupRecordWithUsers, *xerrors.AppError)
	GetGroupSharedSecrets(name string) (*GroupRecordWithSecrets, *xerrors.AppError)
	UpdateGroupName(newName string, groupName string) (*GroupRecord, *xerrors.AppError)
	DeleteByGroupID(groupID int64) *xerrors.AppError
	NewRecord(name string, ownerID int64) (*GroupRecord, *xerrors.AppError)
	AddUser(groupId, userId int64) *xerrors.AppError
	RemoveUser(groupId, userId int64) *xerrors.AppError
	GetGroupsByUserID(userID int64) ([]GroupRecord, *xerrors.AppError)
	IsUserInGroup(groupID, userID int64) (bool, *xerrors.AppError)
}

type Group struct {
	DB core.Queryable
}

type GroupRecordWithUsers struct {
	GroupRecord
	Users []*users.UserRecord
}

func Repository(db core.Queryable) GroupRepository {
	return &Group{DB: db}
}

func (g *Group) NewRecord(name string, ownerID int64) (*GroupRecord, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start a new transaction
	db, ok := g.DB.(*sql.DB)
	if !ok {
		return nil, xerrors.DatabaseError(fmt.Errorf("failed to cast DB to *sql.DB"), "group.NewRecord")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "group.NewRecord")
	}
	// Rollback transaction in case of error
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Insert the new group into the groups table
	query := `
		INSERT INTO groups (name, creator_id, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, name, creator_id, created_at;
	`
	var newGroup GroupRecord
	err = tx.QueryRowContext(ctx, query, name, ownerID).
		Scan(&newGroup.ID, &newGroup.Name, &newGroup.CreatorID, &newGroup.CreatedAt)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "group.NewRecord: failed to create group")
	}

	// Add the creator as a member of the group in the group_members table
	addUserQuery := `
		INSERT INTO group_members (group_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (group_id, user_id) DO NOTHING;
	`
	_, err = tx.ExecContext(ctx, addUserQuery, newGroup.ID, ownerID)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "group.NewRecord: failed to add creator as member")
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, xerrors.DatabaseError(err, "group.NewRecord: failed to commit transaction")
	}

	return &newGroup, nil
}

func (g *Group) GetByGroupID(id int64) (*GroupRecordWithUsers, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// First query to get the group details
	queryGroup := `
		SELECT id, name, creator_id, created_at
		FROM groups
		WHERE id = $1;
	`

	var group GroupRecordWithUsers
	err := g.DB.QueryRowContext(ctx, queryGroup, id).Scan(&group.ID, &group.Name, &group.CreatorID, &group.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, xerrors.DatabaseError(err, "group.GetByGroupID")
		}
		return nil, xerrors.DatabaseError(err, "group.GetByGroupID")
	}

	// Second query to get users related to the group
	queryUsers := `
		SELECT u.id, u.email
		FROM users u
		JOIN group_members gm ON gm.user_id = u.id
		WHERE gm.group_id = $1;
	`

	rows, err := g.DB.QueryContext(ctx, queryUsers, group.ID)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "group.GetByGroupID")
	}
	defer rows.Close()

	for rows.Next() {
		var user users.UserRecord
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			return nil, xerrors.DatabaseError(err, "group.GetByGroupID")
		}
		group.Users = append(group.Users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, xerrors.DatabaseError(err, "group.GetByGroupID")
	}

	return &group, nil
}

func (g *Group) GetGroupUsers(name string) (*GroupRecordWithUsers, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// First query to get the group details
	queryGroup := `
		SELECT id, name, creator_id, created_at
		FROM groups
		WHERE name = $1;
	`

	var group GroupRecordWithUsers
	err := g.DB.QueryRowContext(ctx, queryGroup, name).Scan(&group.ID, &group.Name, &group.CreatorID, &group.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, xerrors.DatabaseError(err, "group.GetByGroupID")
		}
		return nil, xerrors.DatabaseError(err, "group.GetByGroupID")
	}

	// Second query to get users related to the group
	queryUsers := `
		SELECT u.id, u.email
		FROM users u
		JOIN group_members gm ON gm.user_id = u.id
		WHERE gm.group_id = $1;
	`

	rows, err := g.DB.QueryContext(ctx, queryUsers, group.ID)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "group.GetByGroupID")
	}
	defer rows.Close()

	for rows.Next() {
		var user users.UserRecord
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			return nil, xerrors.DatabaseError(err, "group.GetByGroupID")
		}
		group.Users = append(group.Users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, xerrors.DatabaseError(err, "group.GetByGroupID")
	}

	return &group, nil
}

type SharedSecretDetailForGroup struct {
	SecretID      int64     `json:"secret_id"`
	Name          string    `json:"name"`
	EncryptedData []byte    `json:"encrypted_data"`
	IV            []byte    `json:"iv"`
	OwnerID       int64     `json:"owner_id"`
	Permission    string    `json:"permission"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GroupRecordWithSecrets struct {
	GroupID   int64                         `json:"group_id"`
	Name      string                        `json:"name"`
	CreatorID int64                         `json:"creator_id"`
	CreatedAt time.Time                     `json:"created_at"`
	Secrets   []*SharedSecretDetailForGroup `json:"secrets"`
}

// GetGroupSharedSecrets fetches all secrets shared with the specified group and returns full secret details.
func (g *Group) GetGroupSharedSecrets(name string) (*GroupRecordWithSecrets, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// First query to get the group details
	queryGroup := `
        SELECT id, name, creator_id, created_at
        FROM groups
        WHERE name = $1;
    `

	var group GroupRecordWithSecrets
	err := g.DB.QueryRowContext(ctx, queryGroup, name).Scan(&group.GroupID, &group.Name, &group.CreatorID, &group.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, xerrors.DatabaseError(err, "group.GetGroupSharedSecrets - group not found")
		}
		return nil, xerrors.DatabaseError(err, "group.GetGroupSharedSecrets")
	}

	// Second query to get secrets shared with the group
	querySecrets := `
        SELECT s.id AS secret_id, s.name, s.encrypted_data, s.iv, s.owner_id, ssg.permission, s.created_at, s.updated_at
        FROM secrets s
        JOIN shared_secrets_group ssg ON ssg.secret_id = s.id
        WHERE ssg.group_id = $1;
    `

	rows, err := g.DB.QueryContext(ctx, querySecrets, group.GroupID)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "group.GetGroupSharedSecrets - query secrets")
	}
	defer rows.Close()

	for rows.Next() {
		var secret SharedSecretDetailForGroup
		if err := rows.Scan(&secret.SecretID, &secret.Name, &secret.EncryptedData, &secret.IV, &secret.OwnerID, &secret.Permission, &secret.CreatedAt, &secret.UpdatedAt); err != nil {
			return nil, xerrors.DatabaseError(err, "group.GetGroupSharedSecrets - scan secret")
		}
		group.Secrets = append(group.Secrets, &secret)
	}

	if err := rows.Err(); err != nil {
		return nil, xerrors.DatabaseError(err, "group.GetGroupSharedSecrets - rows error")
	}

	return &group, nil
}

func (g *Group) UpdateGroupName(newName string, groupName string) (*GroupRecord, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE groups
		SET name = $1
		WHERE name = $2
		RETURNING id, name, creator_id, created_at;
	`

	var updatedGroup GroupRecord
	err := g.DB.QueryRowContext(ctx, query, newName, groupName).
		Scan(&updatedGroup.ID, &updatedGroup.Name, &updatedGroup.CreatorID, &updatedGroup.CreatedAt)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "group.UpdateByGroupID")
	}

	return &updatedGroup, nil
}

func (g *Group) DeleteByGroupID(groupID int64) *xerrors.AppError {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM groups WHERE id = $1;`

	_, err := g.DB.ExecContext(ctx, query, groupID)
	if err != nil {
		return xerrors.DatabaseError(err, "group.DeleteByGroupID")
	}

	return nil
}

// Returns no error if user already in group
func (g *Group) AddUser(groupId, userId int64) *xerrors.AppError {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the SQL query to insert a new user into the group_members table
	query := `
		INSERT INTO group_members (group_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (group_id, user_id) DO NOTHING;  -- Prevents duplicate entries
	`

	_, err := g.DB.ExecContext(ctx, query, groupId, userId)
	if err != nil {
		return xerrors.DatabaseError(err, "group.AddUser")
	}

	return nil
}

func (g *Group) RemoveUser(groupId, userId int64) *xerrors.AppError {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the SQL query to remove a user from the group_members table
	query := `
		DELETE FROM group_members
		WHERE group_id = $1 AND user_id = $2;
	`

	result, err := g.DB.ExecContext(ctx, query, groupId, userId)
	if err != nil {
		return xerrors.DatabaseError(err, "group.RemoveUserIfExists")
	}

	// Check if any rows were deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return xerrors.DatabaseError(err, "group.RemoveUserIfExists")
	}

	if rowsAffected == 0 {
		// return xerrors.DatabaseError(fmt.Errorf("No user in group with id: %d", userId),
		// 	"group.RemoveUserIfExists")
		xerrors.ClientError(http.StatusNotFound,
			fmt.Sprintf("Invalid user_id=%d or group_id=%d", groupId, userId),
			"group.RemoveUserIfExists",
			fmt.Errorf("Invalid user or group"))
	}

	return nil
}

func (g *Group) IsUserInGroup(groupID, userID int64) (bool, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the SQL query to check if the user is the creator or a member of the group
	query := `
		SELECT 1
		FROM groups
		WHERE id = $1 AND (creator_id = $2 OR EXISTS (
			SELECT 1 
			FROM group_members 
			WHERE group_id = $1 AND user_id = $2
		));
	`

	var exists int
	err := g.DB.QueryRowContext(ctx, query, groupID, userID).Scan(&exists)

	// If no rows were found, the user is neither a member nor the creator of the group
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		// Handle database error
		return false, xerrors.DatabaseError(err, "group.IsUserInGroup")
	}

	// User is either the creator or a member of the group
	return true, nil
}
