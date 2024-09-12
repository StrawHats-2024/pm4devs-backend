package db

import (
	"fmt"
	"pm4devs-backend/pkg/models"
)

func (pg *PostgresStore) CreateGroup(group *models.Group) (int, error) {
	// Start a transaction
	tx, err := pg.db.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback() // Rollback if anything goes wrong

	// Step 1: Insert the new group
	query := `
		INSERT INTO Groups (group_name, created_by) 
		VALUES ($1, $2) RETURNING group_id;`

	var groupId int
	err = tx.QueryRow(query, group.GroupName, group.CreatedBy).Scan(&groupId)
	if err != nil {
		return -1, err
	}

	// Step 2: Insert the creator as an admin in UserGroup
	query = `
		INSERT INTO UserGroup (user_id, group_id, role) 
		VALUES ($1, $2, $3);`

	_, err = tx.Exec(query, group.CreatedBy, groupId, models.Admin)
	if err != nil {
		return -1, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return -1, err
	}

	return groupId, nil
}

func (pg *PostgresStore) GetGroupsOfUser(userId int) ([]*GetUserGroupByIdRes, error) {
	query := `
    SELECT 
        g.group_id,
        g.group_name,
        ug.role
    FROM 
        UserGroup ug
    JOIN 
        Groups g ON ug.group_id = g.group_id
    WHERE 
        ug.user_id = $1;`

	rows, err := pg.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*GetUserGroupByIdRes
	for rows.Next() {
		var group GetUserGroupByIdRes
		if err := rows.Scan(&group.GroupID, &group.GroupName, &group.Role); err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}

	return groups, nil
}

func (pg *PostgresStore) AddUserToGroup(group_id int, newUser AddUserToGroupReq) error {
	user, err := pg.GetUserByEmail(newUser.UserEmail)
	if err != nil {
		return err
	}
	query := `
    INSERT INTO UserGroup (user_id, group_id, role) 
    VALUES ($1, $2, $3);`

	_, err = pg.db.Exec(query, user.UserID, group_id, newUser.Role)
	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresStore) DeleteUserFromGroup(groupID int, userEmail string) error {
	// Retrieve the user by their email
	user, err := pg.GetUserByEmail(userEmail)
	if err != nil {
		return err
	}

	// Prepare the SQL query to delete the user from the specified group
	query := `
    DELETE FROM UserGroup 
    WHERE user_id = $1 AND group_id = $2;`

	// Execute the query
	result, err := pg.db.Exec(query, user.UserID, groupID)
	if err != nil {
		return err
	}

	// Check if the user was actually removed
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with email %s is not a member of group %d", userEmail, groupID)
	}

	return nil
}

func (pg *PostgresStore) GetAllUsersByGroupID(groupID int) ([]UserRes, error) {
	query := `
		SELECT u.user_id, u.username
		FROM UserGroup ug
		JOIN Users u ON ug.user_id = u.user_id
		WHERE ug.group_id = $1;`

	rows, err := pg.db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserRes
	for rows.Next() {
		var user UserRes
		if err := rows.Scan(&user.UserID, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

type AddUserToGroupReq struct {
	UserEmail string      `json:"user_email"`
	Role      models.Role `json:"role"`
}

type GetUserGroupByIdRes struct {
	GroupID   int    `json:"group_id"`
	GroupName string `json:"group_name"`
	Role      string `json:"role"`
}

type UserRes struct {
	UserID   int    `json:"user_id" db:"user_id"`
	Username string `json:"username" db:"username"`
}


