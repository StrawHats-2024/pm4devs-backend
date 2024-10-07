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
	UpdateGroupName(newName string, groupID int64) (*GroupRecord, *xerrors.AppError)
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

	query := `
		INSERT INTO groups (name, creator_id, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, name, creator_id, created_at;
	`

	var newGroup GroupRecord
	err := g.DB.QueryRowContext(ctx, query, name, ownerID).
		Scan(&newGroup.ID, &newGroup.Name, &newGroup.CreatorID, &newGroup.CreatedAt)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "group.NewRecord")
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

func (g *Group) UpdateGroupName(newName string, groupID int64) (*GroupRecord, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE groups
		SET name = $1
		WHERE id = $2
		RETURNING id, name, creator_id, created_at;
	`

	var updatedGroup GroupRecord
	err := g.DB.QueryRowContext(ctx, query, newName, groupID).
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
