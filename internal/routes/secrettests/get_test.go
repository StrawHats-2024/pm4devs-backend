package secret

import (
	"net/http"
	"testing"

	gofakeit "github.com/brianvoe/gofakeit/v7"
	"pm4devs.strawhats/internal/assert"
	"pm4devs.strawhats/internal/mocks"
	"pm4devs.strawhats/internal/models/secrets"
	"pm4devs.strawhats/internal/routes/group"
	"pm4devs.strawhats/internal/routes/secret"
	"pm4devs.strawhats/internal/routes/utils"
)

func TestGetUsrSecrets(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := secretsHandler(app)
	authHandler := utils.AuthHandler(app)

	credentials := `{"email": "test@example.com", "password": "password"}`

	assert.Check(t, utils.RegisterUser(authHandler, credentials))
	token := utils.LoginUser(authHandler, credentials)
	assert.Check(t, len(token) > 0)

	type responseMessage struct {
		Message string `json:"message"`
		Data    []any  `json:"data"`
	}

	tests := []assert.HandlerTestCase[responseMessage]{
		{
			Name:   "AuthRequired",
			Body:   ``,
			Status: http.StatusUnauthorized,
		},
		{
			Name:   "Success",
			Auth:   token,
			Status: http.StatusOK,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Success!")
			},
		},
	}

	for _, tc := range tests {
		assert.RunHandlerTestCase(t, handler, http.MethodGet, secret.GetUserSecretsRoute, tc)
	}

	user, err := app.Models.Users.GetByEmail("test@example.com")
	if err != nil {
		t.Error(err)
	}
	_, err = app.Models.Secrets.NewRecord(gofakeit.Name(),
		gofakeit.Sentence(5), user.ID)
	if err != nil {
		t.Error(err)
	}
	assert.RunHandlerTestCase(t, handler, http.MethodGet, secret.GetUserSecretsRoute, assert.HandlerTestCase[responseMessage]{
		Name:   "DummyData/Success",
		Auth:   token,
		Status: http.StatusOK,
		FN: func(t *testing.T, result responseMessage) {
			assert.Equal(t, result.Message, "Success!")
			assert.NotEqual(t, len(result.Data), 0)
		},
	})
}

func TestGetGroupSecrets(t *testing.T) {
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

	credentialsThree := `{"email": "test3@example.com", "password": "password"}`
	assert.Check(t, utils.RegisterUser(authHandler, credentialsThree))
	tokenThree := utils.LoginUser(authHandler, credentialsThree)
	assert.Check(t, len(tokenThree) > 0)

	// Create a group with user 1
	groupData := `{"group_name": "TestGroup"}`
	groupRes := sendAuthRequest(groupHandler, http.MethodPost, group.CRUDGroupRoute, groupData, token)
	assert.Equal(t, groupRes, http.StatusCreated)

	// Add second user to the group
	addMemberData := `{"group_id": 1, "user_id": 2}`
	addMemberRes := sendAuthRequest(groupHandler, http.MethodPost, group.AddUserToGroupRoute, addMemberData, token)
	assert.Equal(t, addMemberRes, http.StatusOK)

	// Create a secret
	secretData := `{"encrypted_data": "test@example.com", "name": "testname"}`
	secretRes := sendAuthRequest(secretsHandler, http.MethodPost, secret.SecretCRUDRoute, secretData, token)
	assert.Equal(t, secretRes, http.StatusCreated)

	// Share secret with the group
	shareSecretdata := `{"secret_id": 1, "group_id": 1, "permission": "read-only"}`
	shareRes := sendAuthRequest(secretsHandler, http.MethodPost, secret.SecretShareGroupRoute, shareSecretdata, token)
	assert.Equal(t, shareRes, http.StatusCreated)

	// Struct to hold the response
	type responseMessage struct {
		Error   map[string]string      `json:"error"`
		Message string                 `json:"message"`
		Data    []secrets.SecretRecord `json:"data"`
	}

	tests := []assert.HandlerTestCase[responseMessage]{
		// Missing or invalid group ID
		{
			Name:   "MissingGroupID",
			Body:   `{}`,
			Method: http.MethodGet,
			Status: http.StatusUnprocessableEntity,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["group_id"], "must be provided")
			},
		},
		{
			Name:   "GroupOwnerCanAccess",
			Body:   `{"group_id": 1}`,
			Method: http.MethodGet,
			Status: http.StatusOK,
			Auth:   token, // User 1
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Success!")
				assert.Check(t, len(result.Data) > 0) // Ensure secrets are returned
			},
		},
		// Successful retrieval
		{
			Name:   "Success",
			Body:   `{"group_id": 1}`,
			Method: http.MethodGet,
			Status: http.StatusOK,
			Auth:   token, // User 2 (group member) trying to access secrets
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Success!")
				assert.Check(t, len(result.Data) > 0) // Ensure secrets are returned
			},
		},
		// Unauthorized if user is not part of the group
		{
			Name:   "Unauthorized",
			Body:   `{"group_id": 1}`,
			Method: http.MethodGet,
			Status: http.StatusUnauthorized,
			Auth:   tokenThree, // User 3
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Message, "Only group members can access secrets.")
			},
		},
	}

	for _, tc := range tests {
		assert.RunHandlerTestCase(t, secretsHandler, tc.Method, secret.GetGroupSecretsRoute, tc)
	}
}
