package groups

import (
	"database/sql"
	"fmt"
	"pm4devs-backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (pg *Store) UpdateGroupName(groupId int, newName string) error {
	// Prepare the SQL query to update the group name
	query := `
		UPDATE Groups
		SET group_name = $1
		WHERE group_id = $2;
	`

	// Execute the update query
	_, err := pg.db.Exec(query, newName, groupId)
	if err != nil {
		return err
	}

	return nil
}

func (pg *Store) DeleteGroup(groupID int) error {
	query := `DELETE FROM Groups WHERE group_id = $1;`
	result, err := pg.db.Exec(query, groupID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No group found with id %d", groupID)
	}

	return nil
}

func (pg *Store) CreateGroup(group *types.Group) (int, error) {
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

	_, err = tx.Exec(query, group.CreatedBy, groupId, types.Admin)
	if err != nil {
		return -1, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return -1, err
	}

	return groupId, nil
}

func (pg *Store) GetUserGroups(userId int) (
	[]*types.GetUserGroupRes, error,
) { // returns all groups that the user with the given userId belongs to
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

	var groups []*types.GetUserGroupRes
	for rows.Next() {
		var group types.GetUserGroupRes
		if err := rows.Scan(&group.GroupID, &group.GroupName, &group.Role); err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}

	return groups, nil
}

func (pg *Store) GetUserByEmail(email string) (*types.User, error) {
	query := `
  SELECT *
  FROM Users
  WHERE email = $1;
  `
	// Initialize a User struct to hold the result
	user := new(types.User)

	// Use QueryRow to fetch the single row based on email
	err := pg.db.QueryRow(query, email).Scan(
		&user.UserID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.LastLogin,
	)

	// Handle the case where no row is found
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, err
	}

	return user, nil
}

func (pg *Store) AddUserToGroup(groupId int, newUser types.AddUserToGroupPayload) error {
	user, err := pg.GetUserByEmail(newUser.UserEmail)
	if err != nil {
		return err
	}
	query := `
    INSERT INTO UserGroup (user_id, group_id, role) 
    VALUES ($1, $2, $3);`

	_, err = pg.db.Exec(query, user.UserID, groupId, newUser.Role)
	if err != nil {
		return err
	}

	return nil
}

func (pg *Store) DeleteUserFromGroup(groupID int, userEmail string) error {
	// Retrieve the user by their email
	fmt.Println("Entered the delet user func")
	user, err := pg.GetUserByEmail(userEmail)
	fmt.Println("user: ", user)
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

// gets list of users belonging to the group with the specified groupID
func (pg *Store) GetUsersByGroupId(groupID int) ([]types.UserInGroup, error) {
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

	var users []types.UserInGroup
	for rows.Next() {
		var user types.UserInGroup
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

func (pg *Store) GetGroupById(groupId int) (*types.GroupWithUsers, error) {
	// Query to get group details
	groupQuery := `
		SELECT group_id, group_name, created_by 
		FROM Groups 
		WHERE group_id = $1;`

	var group types.Group
	err := pg.db.QueryRow(groupQuery, groupId).Scan(&group.GroupID, &group.GroupName, &group.CreatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group with ID %d not found", groupId)
		}
		return nil, err
	}

	// Query to get users in the group
	usersQuery := `
		SELECT u.user_id, u.email, u.username, ug.role
		FROM Users u
		JOIN UserGroup ug ON u.user_id = ug.user_id
		WHERE ug.group_id = $1;`

	rows, err := pg.db.Query(usersQuery, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.GroupUserItem
	for rows.Next() {
		var user types.GroupUserItem
		if err := rows.Scan(&user.UserID, &user.Email, &user.Username, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return the group along with users
	return &types.GroupWithUsers{
		Group: group,
		Users: users,
	}, nil
}

// IsUserAdminInGroup checks if the user is an admin in the specified group.
func (pg *Store) IsUserAdminInGroup(userId, groupId int) (bool, error) {
	fmt.Println("IsUserAdminInGroup entered ")
	query := `
	SELECT COUNT(*) 
	FROM UserGroup 
	WHERE user_id = $1 AND group_id = $2 AND role = $3;`

	var count int
	err := pg.db.QueryRow(query, userId, groupId, types.Admin).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// IsGroupCreator checks if the user is the creator of the specified group.
func (pg *Store) IsGroupCreator(userId, groupId int) (bool, error) {
	fmt.Println("IsGroupCreator entered ")
	query := `
	SELECT COUNT(*) 
	FROM Groups 
	WHERE group_id = $1 AND created_by = $2;`

	var count int
	err := pg.db.QueryRow(query, groupId, userId).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetUserRoleInGroup retrieves the role of a user in a specified group based on their email.
func (pg *Store) GetUserRoleInGroup(userEmail string, groupId int) (types.Role, error) {
	fmt.Println("GetUserRoleInGroup entered ")
	var role types.Role

	// SQL query to retrieve the user's role in the specified group
	query := `
    SELECT ug.role
    FROM UserGroup ug
    JOIN Users u ON ug.user_id = u.user_id
    WHERE u.email = $1 AND ug.group_id = $2;
  `
	err := pg.db.QueryRow(query, userEmail, groupId).Scan(&role)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found in group")
		}
		return "", err // Return any other errors
	}

	return role, nil // Return the user's role
}
