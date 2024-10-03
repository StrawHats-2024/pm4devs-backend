package secret

import (
	"net/http"

	"pm4devs.strawhats/internal/app"
	"pm4devs.strawhats/internal/mailer"
	"pm4devs.strawhats/internal/models/tokens"
	"pm4devs.strawhats/internal/models/users"
	"pm4devs.strawhats/internal/rest"
	"pm4devs.strawhats/internal/routes/middleware"
	"pm4devs.strawhats/internal/xlogger"
)

// Encapsulates the Application dependencies required by routes
type Secret struct {
	bg     app.Backgrounder
	logger xlogger.Logger
	mailer mailer.Mailer
	rest   *rest.Rest
	tokens tokens.TokensRepository
	users  users.UsersRepository
}

func New(app *app.App) *Secret {
	return &Secret{
		bg:     app.BG,
		logger: app.Logger,
		mailer: app.Mailer,
		rest:   app.Rest,
		tokens: app.Models.Tokens,
		users:  app.Models.Users,
	}
}

const GetUserSecretsRoute = "/v1/secrets/user"

func (s *Secret) Route(mux *http.ServeMux, mw *middleware.Middleware) {
	mux.HandleFunc(GetUserSecretsRoute, mw.Authenticated(s.getUserSecrets))
}
