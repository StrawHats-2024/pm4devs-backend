package secret

import (
	"net/http"

	"pm4devs.strawhats/internal/app"
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

// Helper failure type
type failure struct {
	Error string `json:"error"`
}

