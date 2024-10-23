package secret

import (
	"net/http"
	"testing"

	"pm4devs.strawhats/internal/assert"
	"pm4devs.strawhats/internal/mocks"
	"pm4devs.strawhats/internal/routes/secret"
	"pm4devs.strawhats/internal/routes/utils"
)

func TestCRUDSecrets(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := secretsHandler(app)
	authHandler := utils.AuthHandler(app)

	credentials := `{"email": "test@example.com", "password": "password"}`

	assert.Check(t, utils.RegisterUser(authHandler, credentials))
	token := utils.LoginUser(authHandler, credentials)
	assert.Check(t, len(token) > 0)

	credentialsTwo := `{"email": "test@example2.com", "password": "password"}`

	assert.Check(t, utils.RegisterUser(authHandler, credentialsTwo))
	tokenTwo := utils.LoginUser(authHandler, credentialsTwo)
	assert.Check(t, len(tokenTwo) > 0)

	type responseMessage struct {
		Error   map[string]string `json:"error"`
		Message string            `json:"message"`
		Data    string            `json:"data"`
	}

	tests := []assert.HandlerTestCase[responseMessage]{
		{
			Name:   "MethodNotAllowed",
			Method: http.MethodPut,
			Status: http.StatusMethodNotAllowed,
			Auth:   token,
		},
		// create POST
		{
			Name:   "Name/Validation",
			Body:   `{"name": "password"}`,
			Status: http.StatusUnprocessableEntity,
			Method: http.MethodPost,
			Auth:   token,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["encrypted_data"], "must be provided")
			},
		},
		{
			Name:   "Data/Validation",
			Body:   `{"encrypted_data": "test@example.com"}`,
			Auth:   token,
			Method: http.MethodPost,
			Status: http.StatusUnprocessableEntity,
			FN: func(t *testing.T, result responseMessage) {
				assert.Equal(t, result.Error["name"], "must be provided")
			},
		},
		{
			Name:   "AuthRequired",
			Body:   ``,
			Method: http.MethodPost,
			Status: http.StatusUnauthorized,
		},
		{
			Name:   "Success",
			Body:   `{"encrypted_data": "test@example.com", "name": "testname", "iv": "testing"}`,
			Auth:   token,
			Method: http.MethodPost,
			Status: http.StatusCreated,
		},

		// get GET
		{
			Name:   "GET/InvalidBody",
			Body:   ``,
			Auth:   token,
			Method: http.MethodGet,
			Status: http.StatusBadRequest,
		},
		{
			Name:   "GET/InvalidSecretID",
			Body:   `{"secret_id": 0}`,
			Auth:   token,
			Method: http.MethodGet,
			Status: http.StatusUnprocessableEntity,
		},
		{
			Name:   "GET/StatusUnauthorized",
			Body:   `{"secret_id": 1}`,
			Auth:   tokenTwo,
			Method: http.MethodGet,
			Status: http.StatusUnauthorized,
		},
		{
			Name:   "GET/Success",
			Body:   `{"secret_id": 1}`,
			Auth:   token,
			Method: http.MethodGet,
			Status: http.StatusOK,
		},
		// UPDATE
		{
			Name:   "PATCH/InvalidBody",
			Body:   ``,
			Auth:   token,
			Method: http.MethodPatch,
			Status: http.StatusBadRequest,
		},
		{
			Name:   "PATCH/InvalidSecretID",
			Body:   `{"secret_id": 0, "name": "ekjlfa", "encrypted_data": "sfjlsk", "iv": "test"}`,
			Auth:   token,
			Method: http.MethodPatch,
			Status: http.StatusUnprocessableEntity,
		},
		{
			Name:   "PATCH/Unauthorized",
			Body:   `{"secret_id": 1, "name": "newname", "encrypted_data": "newdata", "iv": "testing"}`,
			Auth:   tokenTwo,
			Method: http.MethodPatch,
			Status: http.StatusUnauthorized,
		},
		{
			Name:   "PATCH/Success",
			Body:   `{"secret_id": 1, "name": "newname", "encrypted_data": "newdata", "iv": "testing"}`,
			Auth:   token,
			Method: http.MethodPatch,
			Status: http.StatusOK,
		},
		// DELETE
		{
			Name:   "Delete/InvalidBody",
			Body:   ``,
			Auth:   token,
			Method: http.MethodDelete,
			Status: http.StatusBadRequest,
		},
		{
			Name:   "Delete/InvalidSecretID",
			Body:   `{"secret_id": 0}`,
			Auth:   token,
			Method: http.MethodDelete,
			Status: http.StatusUnprocessableEntity,
		},
		{
			Name:   "Delete/StatusUnauthorized",
			Body:   `{"secret_id": 1}`,
			Auth:   tokenTwo,
			Method: http.MethodDelete,
			Status: http.StatusUnauthorized,
		},
		{
			Name:   "Delete/Success",
			Body:   `{"secret_id": 1}`,
			Auth:   token,
			Method: http.MethodDelete,
			Status: http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		assert.RunHandlerTestCase(t, handler, tc.Method, secret.SecretCRUDRoute, tc)
	}

}
