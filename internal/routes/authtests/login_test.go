package auth

import (
	"net/http"
	"testing"

	"pm4devs.strawhats/internal/assert"
	"pm4devs.strawhats/internal/mocks"
	"pm4devs.strawhats/internal/routes/auth"
	"pm4devs.strawhats/internal/routes/utils"
)

// Tests error cases for login
func TestLoginValidation(t *testing.T) {
	assert.Integration(t)

	app := mocks.App(t)
	handler := utils.AuthHandler(app)

	type failure struct {
		Error map[string]string `json:"error"`
	}

	tests := []assert.HandlerTestCase[failure]{
		{
			Name:   "Email/Validation",
			Body:   `{"password": "password"}`,
			Status: http.StatusUnprocessableEntity,
			FN: func(t *testing.T, result failure) {
				assert.Equal(t, result.Error["email"], "must be provided")
			},
		},
		{
			Name:   "Password/Validation",
			Body:   `{"email": "test@example.com"}`,
			Status: http.StatusUnprocessableEntity,
			FN: func(t *testing.T, result failure) {
				assert.Equal(t, result.Error["password"], "must be provided")
			},
		},
	}

	for _, tc := range tests {
		assert.RunHandlerTestCase(t, handler, "POST", auth.LoginRoute, tc)
	}
}

// Tests error cases for login
func TestLoginUnauthorized(t *testing.T) {
	assert.Integration(t)
	app := mocks.App(t)
	handler := utils.AuthHandler(app)

	type success struct {
		Token string `json:"token"`
	}

	credentials := `{"email": "test@example.com", "password": "password"}`

	// User does not exist
	assert.RunHandlerTestCase(t, handler, "POST", auth.LoginRoute, assert.HandlerTestCase[failure]{
		Name:   "User/DoesNotExist",
		Body:   credentials,
		Status: http.StatusUnauthorized,
	})

	// Seed – create user
	assert.Check(t, utils.RegisterUser(handler, credentials))

	// Password does not match
	assert.RunHandlerTestCase(t, handler, "POST", auth.LoginRoute, assert.HandlerTestCase[failure]{
		Name:   "User/WrongPassword",
		Body:   `{"email": "test@example.com", "password": "pa55word"}`,
		Status: http.StatusUnauthorized,
		FN: func(t *testing.T, result failure) {
			assert.Equal(t, result.Error, "The provided credentials are invalid")
		},
	})

	// Seed - activate user
	assert.Check(t, utils.ActivateUser(handler, app))

	// Success
	assert.RunHandlerTestCase(t, handler, "POST", auth.LoginRoute, assert.HandlerTestCase[success]{
		Name:   "Login/Success",
		Body:   credentials,
		Status: http.StatusOK,
		FN: func(t *testing.T, result success) {
			assert.True(t, len(result.Token) > 0)
		},
	})
}
