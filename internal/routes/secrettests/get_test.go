package secret

import (
	"net/http"
	"testing"

	gofakeit "github.com/brianvoe/gofakeit/v7"
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
