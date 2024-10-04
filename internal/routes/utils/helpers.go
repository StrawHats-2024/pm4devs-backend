package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"pm4devs.strawhats/internal/app"
	"pm4devs.strawhats/internal/mocks"
	"pm4devs.strawhats/internal/models/users"
	"pm4devs.strawhats/internal/routes/auth"
	"pm4devs.strawhats/internal/routes/middleware"
)

// ============================================================================
// Helpers
// ============================================================================

// Creates a complete Auth handler including middleware
func AuthHandler(app *app.App) http.HandlerFunc {
	handler := func() http.Handler {
		mux := http.NewServeMux()

		middleware := middleware.New(app)
		auth := auth.New(app)
		auth.Route(mux, middleware)

		return middleware.User(mux)
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

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
func ActivateUser(handler http.HandlerFunc, app *app.App) bool {
	app.BG.Wait()
	body := fmt.Sprintf(`{"token": "%s"}`, mocks.Mailer(app).WelcomeActivationToken)
	return sendRequest(handler, "PUT", auth.ActivateRoute, body) == http.StatusOK
}

// Helper to login a user
func LoginUser(handler http.HandlerFunc, credentials string) string {
	var result struct {
		Token string `json:"token"`
	}
	sendRequestGetResult(handler, "POST", auth.LoginRoute, credentials, &result)
	return result.Token
}

// Helper to create a user
func RegisterUser(handler http.HandlerFunc, credentials string) bool {
	statusCode := sendRequest(handler, "POST", auth.RegisterRoute, credentials)
	return statusCode == http.StatusCreated
}

// Sends a request and returns the HTTP status
func sendRequest(handler http.HandlerFunc, method, route, body string) int {
	req := httptest.NewRequest(method, route, bytes.NewBufferString(body))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	resp := rr.Result()
	defer resp.Body.Close()
	return resp.StatusCode
}

func sendRequestGetResult[T any](handler http.HandlerFunc, method, route, body string, dst *T) *T {
	req := httptest.NewRequest(method, route, bytes.NewBufferString(body))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	resp := rr.Result()
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&dst)
	return dst
}
