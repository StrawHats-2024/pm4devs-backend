package types

type UserStore interface {
	CreateUser(*User) (int, error)

	// GetUserById retrieves a user by their ID.
	GetUserById(int) (*User, error)

	// GetUserByEmail retrieves a user by their email address.
	GetUserByEmail(string) (*User, error)

	// GetAllUsers retrieves all users from the database.
	// The passwordHash for each user will be set to an empty string for security.
	GetAllUsers() ([]*User, error)

	// UpdateLastLogin updates the last login timestamp for a user identified by their ID.
	UpdateLastLogin(int) error
}

type SecretStore interface {
	// CreateSecret creates a new secret and returns its ID.
	CreateSecret(*Secret) (int, error)

	// GetAllSecret retrieves all secrets from the database.
	GetAllSecret() ([]*Secret, error)

	// GetSecretById retrieves a secret by its ID.
	GetSecretById(int) (*Secret, error)

	// GetAllSecretsByUserID retrieves all secrets associated with a user identified by their ID.
	GetAllSecretsByUserID(int) ([]*Secret, error)

	// DeleteSecretById deletes a secret identified by its ID from the database.
	DeleteSecretById(int) error

	// UpdateSecretById updates a secret identified by its ID using the provided request data.
	UpdateSecretById(int, UpdateSecretPayload) error
}

type GroupStore interface {
	// CreateGroup creates a new group and returns the group's ID.
	CreateGroup(*Group) (int, error)

	// GetUserGroups retrieves a list of groups that a user belongs to, identified by their user ID.
	GetUserGroups(userId int) ([]*GetUserGroupRes, error) // Returns list of groups user belongs to

	// GetGroupById retrieves a group by its ID.
	GetGroupById(groupId int) (*GroupWithUsers, error)

	// AddUserToGroup adds a user to a specified group.
	AddUserToGroup(groupId int, newUser AddUserToGroupPayload) error

	// DeleteUserFromGroup removes a user from a specified group using their email address.
	DeleteUserFromGroup(groupId int, email string) error

	// GetUsersByGroupId retrieves all users belonging to a group identified by its ID.
	GetUsersByGroupId(groupId int) ([]UserInGroup, error) // Returns all users of the group with groupId

	// DeleteGroup deletes a group identified by its ID from the database.
	DeleteGroup(groupId int) error // groupId to delete group

	// UpdateGroupName updates the name of a specified group identified by its ID.
	UpdateGroupName(groupId int, newName string) error

	// IsUserAdminInGroup checks if the user is an admin in the specified group.
	IsUserAdminInGroup(userId, groupId int) (bool, error)

	// IsGroupCreator checks if the user is the creator of the specified group.
	IsGroupCreator(userId, groupId int) (bool, error)

	// Get user role with user_email and groupId
	GetUserRoleInGroup(string, int) (Role, error)
}

