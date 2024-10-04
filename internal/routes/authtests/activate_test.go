package auth

import (
	"fmt"
	"net/http"
	"testing"

	"pm4devs.strawhats/internal/assert"
	"pm4devs.strawhats/internal/mocks"
	"pm4devs.strawhats/internal/routes/auth"
	"pm4devs.strawhats/internal/routes/utils"
)

func TestActivate(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := utils.AuthHandler(app)

	// Seed - create user
	assert.Check(t, utils.RegisterUser(handler, `{"email": "test@example.com", "password": "password"}`))

	// Invalid Token
	assert.RunHandlerTestCase[failure](t, handler, "PUT", auth.ActivateRoute, assert.HandlerTestCase[failure]{
		Name:   "Activate/Invalid",
		Body:   `{"token": "token"}`,
		Status: http.StatusNotFound,
	})

	app.BG.Wait()
	token := mocks.Mailer(app).WelcomeActivationToken

	// Success
	assert.RunHandlerTestCase[user](t, handler, "PUT", auth.ActivateRoute, assert.HandlerTestCase[user]{
		Name:   "Activate/Success",
		Body:   fmt.Sprintf(`{"token": "%s"}`, token),
		Status: http.StatusOK,
		FN: func(t *testing.T, result user) {
			assert.True(t, result.User.Activated)
		},
	})
}
