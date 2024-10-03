package secret

import (
	"net/http"
	"testing"

	"pm4devs.strawhats/internal/assert"
	"pm4devs.strawhats/internal/mocks"
	"pm4devs.strawhats/internal/routes/secret"
	"pm4devs.strawhats/internal/routes/utils"
)

func TestCreateNewSecrets(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := secretsHandler(app)
	authHandler := utils.AuthHandler(app)

	credentials := `{"email": "test@example.com", "password": "password"}`

	assert.Check(t, utils.RegisterUser(authHandler, credentials))
	token := utils.LoginUser(authHandler, credentials)
	assert.Check(t, len(token) > 0)

	type responseMessage struct {
		Error    map[string]string `json:"error"`
		Message  string            `json:"message"`
		SecretID string            `json:"secret_id"`
	}

	tests := []assert.HandlerTestCase[responseMessage]{
		{
			Name:   "Name/Validation",
			Body:   `{"name": "password"}`,
			Status: http.StatusUnprocessableEntity,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["encrypted_data"], "must be provided")
			},
		},
		{
			Name:   "Data/Validation",
			Body:   `{"encrypted_data": "test@example.com"}`,
			Auth:   token,
			Status: http.StatusUnprocessableEntity,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["name"], "must be provided")
			},
		},
		{
			Name:   "AuthRequired",
			Body:   ``,
			Status: http.StatusUnauthorized,
		},
		{
			Name:   "Success",
			Body:   `{"encrypted_data": "test@example.com", "name": "testname"}`,
			Auth:   token,
			Status: http.StatusCreated,
		},
	}

	for _, tc := range tests {
		assert.RunHandlerTestCase(t, handler, http.MethodPost, secret.CreateNewSecretRoute, tc)
	}

}
