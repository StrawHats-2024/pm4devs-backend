package secret

import (
	"net/http"

	"pm4devs.strawhats/internal/app"
	"pm4devs.strawhats/internal/mailer"
	"pm4devs.strawhats/internal/models/group"
	"pm4devs.strawhats/internal/models/secrets"
	"pm4devs.strawhats/internal/models/tokens"
	"pm4devs.strawhats/internal/models/users"
	"pm4devs.strawhats/internal/rest"
	"pm4devs.strawhats/internal/routes/middleware"
	"pm4devs.strawhats/internal/xlogger"
)

// Encapsulates the Application dependencies required by routes
type Secret struct {
	bg      app.Backgrounder
	logger  xlogger.Logger
	mailer  mailer.Mailer
	rest    *rest.Rest
	tokens  tokens.TokensRepository
	users   users.UsersRepository
	secrets secrets.SecretsRepository
	group   group.GroupRepository
}

func New(app *app.App) *Secret {
	return &Secret{
		bg:      app.BG,
		logger:  app.Logger,
		mailer:  app.Mailer,
		rest:    app.Rest,
		tokens:  app.Models.Tokens,
		users:   app.Models.Users,
		secrets: app.Models.Secrets,
		group:   app.Models.Group,
	}
}

func (s *Secret) Route(mux *http.ServeMux, mw *middleware.Middleware) {
	mux.HandleFunc(GetUserSecretsRoute, mw.Authenticated(s.getUserSecrets))
	mux.HandleFunc(SecretCRUDRoute, mw.Authenticated(s.CRUDRoute))
	mux.HandleFunc(SecretShareUserRoute, mw.Authenticated(s.handleShareToUser))
	mux.HandleFunc(SecretShareGroupRoute, mw.Authenticated(s.handleShareToGroup))

	mux.HandleFunc(GetGroupSecretsRoute, mw.Authenticated(s.getGroupSecrets))
}
