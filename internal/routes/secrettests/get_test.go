package secret

import (
	"net/http"
	"testing"

	"pm4devs.strawhats/internal/assert"
	"pm4devs.strawhats/internal/mocks"
	"pm4devs.strawhats/internal/routes/secret"
	"pm4devs.strawhats/internal/routes/utils"
)

func TestGetUsrSecrets(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := secretsHandler(app)
	authHandler := utils.AuthHandler(app)

	credentials := `{"email": "test@example.com", "password": "password"}`

	assert.RunHandlerTestCase(t, handler, "GET", secret.GetUserSecretsRoute, assert.HandlerTestCase[failure]{
		Name:   "GetUserSecretsRoute/AuthRequire",
		Body:   ``,
		Status: http.StatusUnauthorized,
	})

	assert.Check(t, utils.RegisterUser(authHandler, credentials))
	token := utils.LoginUser(authHandler, credentials)
	assert.Check(t, len(token) > 0)

	assert.RunHandlerTestCase(t, handler, "GET", secret.GetUserSecretsRoute, assert.HandlerTestCase[failure]{
		Name:   "GetUserSecretsRoute/Success",
		Auth:   token,
		Body:   ``,
		Status: http.StatusOK,
		FN: func(t *testing.T, result failure) {
				assert.Equal(t, result.Error, "testing error")
		},
	})

}
