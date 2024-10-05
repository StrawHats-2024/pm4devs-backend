package secret

import (
	"net/http"
	"testing"

	"pm4devs.strawhats/internal/assert"
	"pm4devs.strawhats/internal/mocks"
	"pm4devs.strawhats/internal/routes/group"
	"pm4devs.strawhats/internal/routes/secret"
	"pm4devs.strawhats/internal/routes/utils"
)

func TestShareToGroup(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	secretsHandler := secretsHandler(app)
	groupHandler := groupHandler(app)
	authHandler := utils.AuthHandler(app)

	// Register and login users
	credentials := `{"email": "test@example.com", "password": "password"}`
	assert.Check(t, utils.RegisterUser(authHandler, credentials))
	token := utils.LoginUser(authHandler, credentials)
	assert.Check(t, len(token) > 0)

	credentialsTwo := `{"email": "test2@example.com", "password": "password"}`
	assert.Check(t, utils.RegisterUser(authHandler, credentialsTwo))
	tokenTwo := utils.LoginUser(authHandler, credentialsTwo)
	assert.Check(t, len(tokenTwo) > 0)

	// Create a secret
	secretData := `{"encrypted_data": "test@example.com", "name": "testname"}`
	res := sendAuthRequest(secretsHandler, http.MethodPost, secret.SecretCRUDRoute, secretData, token)
	assert.Equal(t, res, http.StatusCreated)

	// Create a group
	groupData := `{"group_name": "TestGroup"}`
	groupRes := sendAuthRequest(groupHandler, http.MethodPost, group.CRUDGroupRoute, groupData, token)
	assert.Equal(t, groupRes, http.StatusCreated)

	type responseMessage struct {
		Error   map[string]string `json:"error"`
		Message string            `json:"message"`
		Data    string            `json:"data"`
	}

	tests := []assert.HandlerTestCase[responseMessage]{
		// Method not allowed
		{
			Name:   "MethodNotAllowed",
			Method: http.MethodPut,
			Status: http.StatusMethodNotAllowed,
			Auth:   token,
		},
		// Missing or invalid SecretID and GroupID
		{
			Name:   "MissingSecretID",
			Body:   `{"group_id": 1, "permission": "read-only"}`,
			Method: http.MethodPost,
			Status: http.StatusUnprocessableEntity,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["secret_id"], "must be provided")
			},
		},
		{
			Name:   "MissingGroupID",
			Body:   `{"secret_id": 1, "permission": "read-only"}`,
			Method: http.MethodPost,
			Status: http.StatusUnprocessableEntity,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["group_id"], "must be provided")
			},
		},
		{
			Name:   "InvalidPermission",
			Body:   `{"secret_id": 1, "group_id": 1, "permission": "invalid"}`,
			Method: http.MethodPost,
			Status: http.StatusUnprocessableEntity,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["permission"], "must be 'read-only' or 'read-write'")
			},
		},
		// Unauthorized if the user is not the owner
		{
			Name:   "Unauthorized",
			Body:   `{"secret_id": 1, "group_id": 1, "permission": "read-only"}`,
			Method: http.MethodPost,
			Status: http.StatusUnauthorized,
			Auth:   tokenTwo,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Only secret owner can manage access")
			},
		},
		// Success case
		{
			Name:   "Success",
			Body:   `{"secret_id": 1, "group_id": 1, "permission": "read-only"}`,
			Method: http.MethodPost,
			Status: http.StatusCreated,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Secret shared successfully with the group.")
			},
		},
	}

	for _, tc := range tests {
		assert.RunHandlerTestCase(t, secretsHandler, tc.Method, secret.SecretShareGroupRoute, tc)
	}
}

func TestUpdateGroupPermission(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	secretsHandler := secretsHandler(app)
	groupHandler := groupHandler(app)
	authHandler := utils.AuthHandler(app)

	// Register and login users
	credentials := `{"email": "test@example.com", "password": "password"}`
	assert.Check(t, utils.RegisterUser(authHandler, credentials))
	token := utils.LoginUser(authHandler, credentials)
	assert.Check(t, len(token) > 0)

	credentialsTwo := `{"email": "test2@example.com", "password": "password"}`
	assert.Check(t, utils.RegisterUser(authHandler, credentialsTwo))
	tokenTwo := utils.LoginUser(authHandler, credentialsTwo)
	assert.Check(t, len(tokenTwo) > 0)

	// Create a secret
	secretData := `{"encrypted_data": "test@example.com", "name": "testname"}`
	res := sendAuthRequest(secretsHandler, http.MethodPost, secret.SecretCRUDRoute, secretData, token)
	assert.Equal(t, res, http.StatusCreated)

	// Create a group
	groupData := `{"group_name": "TestGroup"}`
	groupRes := sendAuthRequest(groupHandler, http.MethodPost, group.CRUDGroupRoute, groupData, token)
	assert.Equal(t, groupRes, http.StatusCreated)

	// share a secret to group
	shareSecretdata := `{"secret_id": 1, "group_id": 1, "permission": "read-write"}`
	shareRes := sendAuthRequest(secretsHandler, http.MethodPost, secret.SecretShareGroupRoute,
		shareSecretdata, token)
	assert.Equal(t, shareRes, http.StatusCreated)

	type responseMessage struct {
		Error   map[string]string `json:"error"`
		Message string            `json:"message"`
		Data    string            `json:"data"`
	}

	tests := []assert.HandlerTestCase[responseMessage]{
		// Missing or invalid SecretID and GroupID
		{
			Name:   "MissingSecretID",
			Body:   `{"group_id": 1, "permission": "read-only"}`,
			Method: http.MethodPatch,
			Status: http.StatusUnprocessableEntity,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["secret_id"], "must be provided")
			},
		},
		{
			Name:   "MissingGroupID",
			Body:   `{"secret_id": 1, "permission": "read-only"}`,
			Method: http.MethodPatch,
			Status: http.StatusUnprocessableEntity,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["group_id"], "must be provided")
			},
		},
		{
			Name:   "InvalidPermission",
			Body:   `{"secret_id": 1, "group_id": 1, "permission": "invalid"}`,
			Method: http.MethodPatch,
			Status: http.StatusUnprocessableEntity,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["permission"], "must be 'read-only' or 'read-write'")
			},
		},
		// Unauthorized if the user is not the owner
		{
			Name:   "Unauthorized",
			Body:   `{"secret_id": 1, "group_id": 1, "permission": "read-write"}`,
			Method: http.MethodPatch,
			Status: http.StatusUnauthorized,
			Auth:   tokenTwo, // User 2 trying to update User 1's group secret
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Only secret owner can manage access")
			},
		},
		// Success case
		{
			Name:   "Success",
			Body:   `{"secret_id": 1, "group_id": 1, "permission": "read-write"}`,
			Method: http.MethodPatch,
			Status: http.StatusOK,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Permission updated successfully for the group.")
			},
		},
	}

	for _, tc := range tests {
		assert.RunHandlerTestCase(t, secretsHandler, tc.Method, secret.SecretShareGroupRoute, tc)
	}
}

func TestRevokeGroupPermission(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	secretsHandler := secretsHandler(app)
	groupHandler := groupHandler(app)
	authHandler := utils.AuthHandler(app)

	// Register and login users
	credentials := `{"email": "test@example.com", "password": "password"}`
	assert.Check(t, utils.RegisterUser(authHandler, credentials))
	token := utils.LoginUser(authHandler, credentials)
	assert.Check(t, len(token) > 0)

	credentialsTwo := `{"email": "test2@example.com", "password": "password"}`
	assert.Check(t, utils.RegisterUser(authHandler, credentialsTwo))
	tokenTwo := utils.LoginUser(authHandler, credentialsTwo)
	assert.Check(t, len(tokenTwo) > 0)

	// Create a secret
	secretData := `{"encrypted_data": "test@example.com", "name": "testname"}`
	res := sendAuthRequest(secretsHandler, http.MethodPost, secret.SecretCRUDRoute, secretData, token)
	assert.Equal(t, res, http.StatusCreated)

	// Create a group
	groupData := `{"group_name": "TestGroup"}`
	groupRes := sendAuthRequest(groupHandler, http.MethodPost, group.CRUDGroupRoute, groupData, token)
	assert.Equal(t, groupRes, http.StatusCreated)

	// Share a secret with the group
	shareSecretData := `{"secret_id": 1, "group_id": 1, "permission": "read-write"}`
	shareRes := sendAuthRequest(secretsHandler, http.MethodPost, secret.SecretShareGroupRoute, shareSecretData, token)
	assert.Equal(t, shareRes, http.StatusCreated)

	type responseMessage struct {
		Error   map[string]string `json:"error"`
		Message string            `json:"message"`
		Data    string            `json:"data"`
	}

	tests := []assert.HandlerTestCase[responseMessage]{
		// Missing or invalid SecretID and GroupID
		{
			Name:   "MissingSecretID",
			Body:   `{"group_id": 1}`,
			Method: http.MethodDelete,
			Status: http.StatusUnprocessableEntity,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["secret_id"], "must be provided")
			},
		},
		{
			Name:   "MissingGroupID",
			Body:   `{"secret_id": 1}`,
			Method: http.MethodDelete,
			Status: http.StatusUnprocessableEntity,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["group_id"], "must be provided")
			},
		},
		// Unauthorized if the user is not the owner
		{
			Name:   "Unauthorized",
			Body:   `{"secret_id": 1, "group_id": 1}`,
			Method: http.MethodDelete,
			Status: http.StatusUnauthorized,
			Auth:   tokenTwo, // User 2 trying to revoke User 1's group secret
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Only secret owner can manage access")
			},
		},
		// Success case
		{
			Name:   "Success",
			Body:   `{"secret_id": 1, "group_id": 1}`,
			Method: http.MethodDelete,
			Status: http.StatusOK,
			Auth:   token, // User 1 revoking permission from Group 1
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Permission revoked successfully for the group.")
			},
		},
	}

	for _, tc := range tests {
		assert.RunHandlerTestCase(t, secretsHandler, tc.Method, secret.SecretShareGroupRoute, tc)
	}
}
