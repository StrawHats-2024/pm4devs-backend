package secret

import (
	"net/http"

	"pm4devs.strawhats/internal/rest"
	"pm4devs.strawhats/internal/routes/middleware"
)

const GetUserSecretsRoute = "/v1/user/secrets"

func (app *Secret) getUserSecrets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user := middleware.ContextGetUser(r)
	userSecrets, err := app.secrets.GetByUserID(user.ID)
	if err != nil {
		app.rest.Error(w, err)
	}
	app.rest.WriteJSON(w, "secret.createNew", http.StatusOK, rest.Envelope{
		"message": "Success!",
		"data":    userSecrets,
	})
}

func (app *Secret) getGroupSecrets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user := middleware.ContextGetUser(r)
	userSecrets, err := app.secrets.GetByUserID(user.ID)
	if err != nil {
		app.rest.Error(w, err)
	}
	app.rest.WriteJSON(w, "secret.createNew", http.StatusOK, rest.Envelope{
		"message": "Success!",
		"data":    userSecrets,
	})
}
