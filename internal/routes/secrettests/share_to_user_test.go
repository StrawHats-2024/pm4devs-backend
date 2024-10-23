package secret

import (
	"net/http"
	"testing"

	"pm4devs.strawhats/internal/assert"
	"pm4devs.strawhats/internal/mocks"
	"pm4devs.strawhats/internal/routes/secret"
	"pm4devs.strawhats/internal/routes/utils"
)

func TestShareToUser(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := secretsHandler(app)
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

	secretData := `{"encrypted_data": "test@example.com", "name": "testname", "iv": "testing"}`
	res := sendAuthRequest(handler, http.MethodPost, secret.SecretCRUDRoute, secretData, token)
	assert.Equal(t, res, http.StatusCreated)

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
		// Missing or invalid SecretID and UserID
		{
			Name:   "MissingSecretID",
			Body:   `{"user_id": 2, "permission": "read-only"}`,
			Method: http.MethodPost,
			Status: http.StatusUnprocessableEntity,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["secret_id"], "must be provided")
			},
		},
		{
			Name:   "MissingUserID",
			Body:   `{"secret_id": 1, "permission": "read-only"}`,
			Method: http.MethodPost,
			Status: http.StatusUnprocessableEntity,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["user_id"], "must be provided")
			},
		},
		{
			Name:   "InvalidPermission",
			Body:   `{"secret_id": 1, "user_id": 2, "permission": "invalid"}`,
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
			Body:   `{"secret_id": 1, "user_id": 2, "permission": "read-only"}`,
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
			Body:   `{"secret_id": 1, "user_id": 2, "permission": "read-only"}`,
			Method: http.MethodPost,
			Status: http.StatusCreated,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Secret shared successfully with the user.")
			},
		},
	}

	for _, tc := range tests {
		assert.RunHandlerTestCase(t, handler, tc.Method, secret.SecretShareUserRoute, tc)
	}
}

func TestUpdateUserPermission(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := secretsHandler(app)
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

	secretData := `{"encrypted_data": "test@example.com", "name": "testname", "iv": "testing"}`
	res := sendAuthRequest(handler, http.MethodPost, secret.SecretCRUDRoute, secretData, token)
	assert.Equal(t, res, http.StatusCreated)

	type responseMessage struct {
		Error   map[string]string `json:"error"`
		Message string            `json:"message"`
		Data    string            `json:"data"`
	}

	tests := []assert.HandlerTestCase[responseMessage]{
		// Missing or invalid SecretID and UserID
		{
			Name:   "MissingSecretID",
			Body:   `{"user_id": 2, "permission": "read-only"}`,
			Method: http.MethodPatch,
			Status: http.StatusUnprocessableEntity,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["secret_id"], "must be provided")
			},
		},
		{
			Name:   "MissingUserID",
			Body:   `{"secret_id": 1, "permission": "read-only"}`,
			Method: http.MethodPatch,
			Status: http.StatusUnprocessableEntity,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["user_id"], "must be provided")
			},
		},
		{
			Name:   "InvalidPermission",
			Body:   `{"secret_id": 1, "user_id": 2, "permission": "invalid"}`,
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
			Body:   `{"secret_id": 1, "user_id": 2, "permission": "read-only"}`,
			Method: http.MethodPatch,
			Status: http.StatusUnauthorized,
			Auth:   tokenTwo,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Only secret owner can manage access")
			},
		},
		// Success case
		{
			Name:   "Success",
			Body:   `{"secret_id": 1, "user_id": 2, "permission": "read-write"}`,
			Method: http.MethodPatch,
			Status: http.StatusOK,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Permission updated successfully for the user.")
			},
		},
	}

	for _, tc := range tests {
		assert.RunHandlerTestCase(t, handler, tc.Method, secret.SecretShareUserRoute, tc)
	}
}

func TestRevokeUserPermission(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := secretsHandler(app)
	authHandler := utils.AuthHandler(app)

	credentials := `{"email": "test@example.com", "password": "password"}`

	// Register and log in user 1
	assert.Check(t, utils.RegisterUser(authHandler, credentials))
	token := utils.LoginUser(authHandler, credentials)
	assert.Check(t, len(token) > 0)

	credentialsTwo := `{"email": "test2@example.com", "password": "password"}`

	// Register and log in user 2
	assert.Check(t, utils.RegisterUser(authHandler, credentialsTwo))
	tokenTwo := utils.LoginUser(authHandler, credentialsTwo)
	assert.Check(t, len(tokenTwo) > 0)

	secretData := `{"encrypted_data": "test@example.com", "name": "testname", "iv": "testing"}`
	res := sendAuthRequest(handler, http.MethodPost, secret.SecretCRUDRoute, secretData, token)
	assert.Equal(t, res, http.StatusCreated)
	// Sample test cases
	type responseMessage struct {
		Error   map[string]string `json:"error"`
		Message string            `json:"message"`
	}

	tests := []assert.HandlerTestCase[responseMessage]{
		{
			Name:   "Validation/MissingSecretID",
			Body:   `{"user_id": 2}`,
			Method: http.MethodDelete,
			Auth:   token,
			Status: http.StatusUnprocessableEntity,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["secret_id"], "must be provided")
			},
		},
		{
			Name:   "Validation/MissingUserID",
			Body:   `{"secret_id": 1}`,
			Method: http.MethodDelete,
			Auth:   token,
			Status: http.StatusUnprocessableEntity,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["user_id"], "must be provided")
			},
		},
		{
			Name:   "Unauthorized/NotOwner",
			Body:   `{"secret_id": 1, "user_id": 2}`,
			Method: http.MethodDelete,
			Auth:   tokenTwo, // User 2 trying to revoke User 1's secret
			Status: http.StatusUnauthorized,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Only secret owner can manage access")
			},
		},
		{
			Name:   "Success",
			Body:   `{"secret_id": 1, "user_id": 2}`,
			Method: http.MethodDelete,
			Auth:   token, // User 1 revoking permission from User 2
			Status: http.StatusOK,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Permission revoked successfully for the user.")
			},
		},
	}

	for _, tc := range tests {
		assert.RunHandlerTestCase(t, handler, tc.Method, secret.SecretShareUserRoute, tc)
	}
}
