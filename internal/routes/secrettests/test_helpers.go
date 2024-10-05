package secret

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"pm4devs.strawhats/internal/app"
	"pm4devs.strawhats/internal/routes/group"
	"pm4devs.strawhats/internal/routes/middleware"
	"pm4devs.strawhats/internal/routes/secret"
)

// ============================================================================
// Helpers
// ============================================================================

// Creates a complete Secrets handler including middleware
func secretsHandler(app *app.App) http.HandlerFunc {
	handler := func() http.Handler {
		mux := http.NewServeMux()

		middleware := middleware.New(app)
		secrets := secret.New(app)
		secrets.Route(mux, middleware)

		return middleware.User(mux)
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}
func groupHandler(app *app.App) http.HandlerFunc {
	handler := func() http.Handler {
		mux := http.NewServeMux()

		middleware := middleware.New(app)
		group := group.New(app)
		group.Route(mux, middleware)

		return middleware.User(mux)
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

func sendAuthRequest(handler http.HandlerFunc, method, route, body, authToken string) int {
	req := httptest.NewRequest(method, route, bytes.NewBufferString(body))

	// If authToken is provided, set the Authorization header
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	resp := rr.Result()
	defer resp.Body.Close()

	return resp.StatusCode
}

// Helper failure type
type failure struct {
	Error string `json:"error"`
}
