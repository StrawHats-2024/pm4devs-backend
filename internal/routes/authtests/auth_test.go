package auth

import (
	"fmt"
	"net/http"
	"testing"

	"pm4devs.strawhats/internal/assert"
	"pm4devs.strawhats/internal/mocks"
	"pm4devs.strawhats/internal/models/users"
	"pm4devs.strawhats/internal/routes/auth"
	"pm4devs.strawhats/internal/routes/utils"
)

func TestAuthE2E(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := utils.AuthHandler(app)

	// Shared credentials
	var bearer string
	credentials := `{"email": "test@example.com", "password": "password"}`

	// Shared responses
	type token struct {
		Token string `json:"token"`
	}

	// Register
	assert.RunHandlerTestCase(t, handler, "POST", auth.RegisterRoute, assert.HandlerTestCase[user]{
		Name:   "Register",
		Body:   credentials,
		Status: http.StatusCreated,
		FN: func(t *testing.T, result user) {
			assert.Equal(t, result.User.Email, "test@example.com")

			app.BG.Wait()
			assert.NotEqual(t, mocks.Mailer(app).WelcomeActivationToken, "")
		},
	})

	// Activate
	assert.RunHandlerTestCase(t, handler, "PUT", auth.ActivateRoute, assert.HandlerTestCase[user]{
		Name:   "Activate",
		Body:   fmt.Sprintf(`{"token": "%s"}`, mocks.Mailer(app).WelcomeActivationToken),
		Status: http.StatusOK,
		FN: func(t *testing.T, result user) {
			assert.True(t, result.User.Activated)
		},
	})

	// Login
	assert.RunHandlerTestCase(t, handler, "POST", auth.LoginRoute, assert.HandlerTestCase[token]{
		Name:   "Login/1",
		Body:   credentials,
		Status: http.StatusOK,
		FN: func(t *testing.T, result token) {
			bearer = result.Token
			assert.NotEqual(t, bearer, "")
		},
	})

	// Logout
	assert.RunHandlerTestCase(t, handler, "POST", auth.LogoutRoute, assert.HandlerTestCase[struct{}]{
		Name:   "Logout",
		Auth:   bearer,
		Body:   ``,
		Status: http.StatusNoContent,
		FN:     nil,
	})

	// Request Reset
	assert.RunHandlerTestCase(t, handler, "POST", auth.ResetRoute, assert.HandlerTestCase[message]{
		Name:   "Reset/Post",
		Body:   `{"email": "test@example.com"}`,
		Status: http.StatusAccepted,
		FN: func(t *testing.T, result message) {
			assert.Equal(t, result.Message, "An email will be sent with reset instructions")

			app.BG.Wait()
			assert.NotEqual(t, mocks.Mailer(app).PasswordResetToken, "")
		},
	})

	// Reset Password
	assert.RunHandlerTestCase(t, handler, "PUT", auth.ResetRoute, assert.HandlerTestCase[message]{
		Name:   "Reset/Put",
		Body:   fmt.Sprintf(`{"password": "pa55word", "token": "%s"}`, mocks.Mailer(app).PasswordResetToken),
		Status: http.StatusOK,
		FN: func(t *testing.T, result message) {
			assert.Equal(t, result.Message, "Your password was reset successfully")
		},
	})

	// Login
	credentials = `{"email": "test@example.com", "password": "pa55word"}`
	assert.RunHandlerTestCase(t, handler, "POST", auth.LoginRoute, assert.HandlerTestCase[token]{
		Name:   "Login/2",
		Body:   credentials,
		Status: http.StatusOK,
		FN: func(t *testing.T, result token) {
			bearer = result.Token
			assert.NotEqual(t, bearer, "")
		},
	})

	// Delete
	assert.RunHandlerTestCase(t, handler, "POST", auth.DeleteRoute, assert.HandlerTestCase[message]{
		Name:   "Delete",
		Auth:   bearer,
		Body:   credentials,
		Status: http.StatusOK,
		FN: func(t *testing.T, result message) {
			assert.Equal(t, result.Message, "Your account has been deleted")
		},
	})
}

// ============================================================================
// Helpers
// ============================================================================

// Creates a complete Auth handler including middleware
// Helper user type
type user struct {
	User users.UserRecord `json:"user"`
}

// Helper success type
type message struct {
	Message string `json:"message"`
}

// Helper failure type
type failure struct {
	Error string `json:"error"`
}

// Helper failures type
type failures struct {
	Error map[string]string `json:"error"`
}

// ============================================================================
// Seeds
// ============================================================================

// Helper to activate a user
