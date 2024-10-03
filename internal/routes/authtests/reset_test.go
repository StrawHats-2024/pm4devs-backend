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

func TestReset(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := utils.AuthHandler(app)
	credentials := `{"email": "test@example.com", "password": "password"}`

	// Bad Request
	assert.RunHandlerTestCase(t, handler, "POST", auth.ResetRoute, assert.HandlerTestCase[failure]{
		Name:   "Reset/BadRequest",
		Body:   ``,
		Status: http.StatusBadRequest,
	})

	// User DNE
	assert.RunHandlerTestCase(t, handler, "POST", auth.ResetRoute, assert.HandlerTestCase[failure]{
		Name:   "Reset/UserDNE",
		Body:   `{"email": "test@example.com"}`,
		Status: http.StatusNotFound,
	})

	// Seed - create user
	assert.Check(t, utils.RegisterUser(handler, credentials))

	// Seed – activate user
	assert.Check(t, utils.ActivateUser(handler, app))

	// Success
	assert.RunHandlerTestCase(t, handler, "POST", auth.ResetRoute, assert.HandlerTestCase[message]{
		Name:   "Reset/Success",
		Body:   `{"email": "test@example.com"}`,
		Status: http.StatusAccepted,
	})

	// Wait for reset token
	app.BG.Wait()

	// Invalid Token
	assert.RunHandlerTestCase(t, handler, "PUT", auth.ResetRoute, assert.HandlerTestCase[failure]{
		Name:   "Reset/BadToken",
		Body:   `{"password": "pa55word", "token": "token"}`,
		Status: http.StatusNotFound,
	})

	// Invalid password
	assert.RunHandlerTestCase(t, handler, "PUT", auth.ResetRoute, assert.HandlerTestCase[failures]{
		Name:   "Reset/InvalidPassword",
		Body:   fmt.Sprintf(`{"password": "please", "token": "%s"}`, mocks.Mailer(app).PasswordResetToken),
		Status: http.StatusUnprocessableEntity,
	})

	// Success
	assert.RunHandlerTestCase(t, handler, "PUT", auth.ResetRoute, assert.HandlerTestCase[failures]{
		Name:   "Reset/InvalidPassword",
		Body:   fmt.Sprintf(`{"password": "pa55word", "token": "%s"}`, mocks.Mailer(app).PasswordResetToken),
		Status: http.StatusOK,
	})
}
