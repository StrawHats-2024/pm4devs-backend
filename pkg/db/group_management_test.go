package db

import (
	"pm4devs-backend/pkg/models"
	"testing"
)

func TestCreateGroup(t *testing.T) {
	db := setupPostgresTestDB(t)
	defer teardownPostgres(t, db)

	pg := &PostgresStore{db: db}
	group := &models.Group{
		GroupName: "Test Group",
		CreatedBy: 1, // Assume user ID 1 exists
	}

	groupID, err := pg.CreateGroup(group)
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	if groupID == 0 {
		t.Errorf("Expected group ID to be greater than 0, got %d", groupID)
	}
}

func TestGetGroupsOfUser(t *testing.T) {
	db := setupPostgresTestDB(t)
	defer teardownPostgres(t, db)

	pg := &PostgresStore{db: db}

	// Assume user ID 1 exists and has groups
	groups, err := pg.GetGroupsOfUser(1)
	if err != nil {
		t.Fatalf("Failed to get groups for user: %v", err)
	}

	if len(groups) == 0 {
		t.Error("Expected to find groups for user, got none")
	}
}

func TestAddUserToGroup(t *testing.T) {
	db := setupPostgresTestDB(t)
	defer teardownPostgres(t, db)

	pg := &PostgresStore{db: db}
	newUser := AddUserToGroupReq{
		UserEmail: "user1@example.com",
		Role:      models.Member,      
	}

	err := pg.AddUserToGroup(1, newUser) 
	if err != nil {
		t.Fatalf("Failed to add user to group: %v", err)
	}
}

func TestDeleteUserFromGroup(t *testing.T) {
	db := setupPostgresTestDB(t)
	defer teardownPostgres(t, db)

	pg := &PostgresStore{db: db}

	err := pg.DeleteUserFromGroup(1, "user1@example.com") // Group ID 1
	if err != nil {
		t.Fatalf("Failed to delete user from group: %v", err)
	}
}

func TestGetAllUsersByGroupID(t *testing.T) {
	db := setupPostgresTestDB(t)
	defer teardownPostgres(t, db)

	pg := &PostgresStore{db: db}

	users, err := pg.GetAllUsersByGroupID(1) // Assume group ID 1 exists
	if err != nil {
		t.Fatalf("Failed to get users for group: %v", err)
	}

	if len(users) == 0 {
		t.Error("Expected to find users in group, got none")
	}
}
