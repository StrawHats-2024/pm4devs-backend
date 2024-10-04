package auth

import (
	"net/http"
	"testing"

	"pm4devs.strawhats/internal/assert"
	"pm4devs.strawhats/internal/mocks"
	"pm4devs.strawhats/internal/routes/auth"
	"pm4devs.strawhats/internal/routes/utils"
)

func TestLogout(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := utils.AuthHandler(app)

	credentials := `{"email": "test@example.com", "password": "password"}`

	// Require Authed User
	assert.RunHandlerTestCase(t, handler, "POST", auth.LogoutRoute, assert.HandlerTestCase[failure]{
		Name:   "Delete/AuthRequire",
		Body:   ``,
		Status: http.StatusUnauthorized,
	})

	// Seed – create user, activate user, login user
	assert.Check(t, utils.RegisterUser(handler, credentials))
	assert.Check(t, utils.ActivateUser(handler, app))
	token := utils.LoginUser(handler, credentials)
	assert.Check(t, len(token) > 0)

	// Success
	assert.RunHandlerTestCase(t, handler, "POST", auth.LogoutRoute, assert.HandlerTestCase[struct{}]{
		Name:   "Delete/Success",
		Auth:   token,
		Body:   ``,
		Status: http.StatusNoContent,
	})
}
