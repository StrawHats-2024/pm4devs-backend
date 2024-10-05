package group

import (
	"net/http"
	"testing"

	"pm4devs.strawhats/internal/assert"
	"pm4devs.strawhats/internal/mocks"
	"pm4devs.strawhats/internal/routes/group"
	"pm4devs.strawhats/internal/routes/utils"
)

func TestAddRemoveUserToGroup(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := groupHandler(app)
	authHandler := utils.AuthHandler(app)

	credentials := `{"email": "test@example.com", "password": "password"}`
	credentialsTwo := `{"email": "test2@example.com", "password": "password"}`

	assert.Check(t, utils.RegisterUser(authHandler, credentials))
	token := utils.LoginUser(authHandler, credentials)

	assert.Check(t, utils.RegisterUser(authHandler, credentialsTwo))
	tokenTwo := utils.LoginUser(authHandler, credentialsTwo)

	assert.Check(t, len(token) > 0)

	type responseMessage struct {
		Error   map[string]string `json:"error"`
		Message string            `json:"message"`
		Data    map[string]any    `json:"data"`
	}

	tests := []assert.HandlerTestCase[responseMessage]{
		// add dummy data
		{
			Name:   "CreateGroup/POST",
			Auth:   token,
			Status: http.StatusCreated,
			Body:   `{"group_name": "testgroup"}`,
			Route:  group.CRUDGroupRoute,
			Method: http.MethodPost,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Success!")
				assert.NotEqual(t, result.Data["group_name"], "")
			},
		},
		// Add User
		{
			Name:   "AuthRequired/AddUser",
			Status: http.StatusUnauthorized,
			Method: http.MethodPost,
			Route:  group.AddUserToGroupRoute,
		},
		{
			Name:   "InvalidGroupID/group.AddUser",
			Auth:   token,
			Status: http.StatusUnprocessableEntity,
			Method: http.MethodPost,
			Route:  group.AddUserToGroupRoute,
			Body:   `{"group_id": 0, "user_id": 1}`,
		},
		{
			Name:   "InvalidUserID/group.AddUser",
			Auth:   token,
			Status: http.StatusUnprocessableEntity,
			Method: http.MethodPost,
			Route:  group.AddUserToGroupRoute,
			Body:   `{"group_id": 1, "user_id": 0}`,
		},
		{
			Name:   "NotGroupOwner/group.AddUser",
			Auth:   tokenTwo,
			Status: http.StatusUnauthorized,
			Method: http.MethodPost,
			Route:  group.AddUserToGroupRoute,
			Body:   `{"group_id": 1, "user_id": 2}`,
		},
		{
			Name:   "ValidRequest/group.AddUser",
			Auth:   token,
			Status: http.StatusOK,
			Method: http.MethodPost,
			Route:  group.AddUserToGroupRoute,
			Body:   `{"group_id": 1, "user_id": 2}`,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Success!")
			},
		},

		// Remove User
		{
			Name:   "AuthRequired/RemoveUser",
			Status: http.StatusUnauthorized,
			Method: http.MethodPost,
			Route:  group.RemoveUserFromGroupRoute,
		},
		{
			Name:   "InvalidGroupID/group.RemoveUser",
			Auth:   token,
			Status: http.StatusUnprocessableEntity,
			Method: http.MethodPost,
			Route:  group.RemoveUserFromGroupRoute,
			Body:   `{"group_id": 0, "user_id": 1}`,
		},
		{
			Name:   "InvalidUserID/group.RemoveUser",
			Auth:   token,
			Status: http.StatusUnprocessableEntity,
			Method: http.MethodPost,
			Route:  group.RemoveUserFromGroupRoute,
			Body:   `{"group_id": 1, "user_id": 0}`,
		},
		{
			Name:   "NotGroupOwner/group.RemoveUser",
			Auth:   tokenTwo,
			Status: http.StatusUnauthorized,
			Method: http.MethodPost,
			Route:  group.RemoveUserFromGroupRoute,
			Body:   `{"group_id": 1, "user_id": 2}`,
		},
		{
			Name:   "ValidRequest/group.RemoveUser",
			Auth:   token,
			Status: http.StatusOK,
			Method: http.MethodPost,
			Route:  group.RemoveUserFromGroupRoute,
			Body:   `{"group_id": 1, "user_id": 2}`,
		},
	}

	for _, tc := range tests {
		assert.RunHandlerTestCase(t, handler, tc.Method, tc.Route, tc)
	}
}
