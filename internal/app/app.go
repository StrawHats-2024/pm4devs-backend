package app

import (
	"pm4devs.strawhats/internal/config"
	"pm4devs.strawhats/internal/mailer"
	"pm4devs.strawhats/internal/models"
	"pm4devs.strawhats/internal/rest"
	"pm4devs.strawhats/internal/xlogger"
)

// Container for app wide dependencies
type App struct {
	BG     Backgrounder
	Config config.Config
	Logger xlogger.Logger
	Mailer mailer.Mailer
	Models *models.Models
	Rest   *rest.Rest
}

// Create a new App struct
func New(
	backgrounder Backgrounder,
	config config.Config,
	logger xlogger.Logger,
	mailer mailer.Mailer,
	models *models.Models,
	rest *rest.Rest,
) *App {
	return &App{
		BG:     backgrounder,
		Config: config,
		Logger: logger,
		Mailer: mailer,
		Models: models,
		Rest:   rest,
	}
}
