package middleware

import (
	"pm4devs.strawhats/internal/app"
	"pm4devs.strawhats/internal/models/permissions"
	"pm4devs.strawhats/internal/models/users"
	"pm4devs.strawhats/internal/rest"
	"pm4devs.strawhats/internal/xlogger"
)

type Middleware struct {
	logger      xlogger.Logger
	permissions permissions.PermissionsRepository
	rest        *rest.Rest
	users       users.UsersRepository
}

func New(app *app.App) *Middleware {
	return &Middleware{
		logger:      app.Logger,
		permissions: app.Models.Permissions,
		rest:        app.Rest,
		users:       app.Models.Users,
	}
}
