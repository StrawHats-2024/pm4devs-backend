package secret

import (
	"net/http"

	"pm4devs.strawhats/internal/rest"
	"pm4devs.strawhats/internal/routes/middleware"
	"pm4devs.strawhats/internal/validator"
)

const GetUserSecretsRoute = "/v1/secrets/user"
const GetGroupSecretsRoute = "/v1/secrets/group"

// TODO: impliment this
const GetSharedByUserSecretRoute = "/v1/secrets/shared"

func (app *Secret) getUserSecrets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user := middleware.ContextGetUser(r)
	userSecrets, err := app.secrets.GetByUserID(user.ID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	app.rest.WriteJSON(w, "secret.createNew", http.StatusOK, rest.Envelope{
		"message": "Success!",
		"data":    userSecrets,
	})
}

func (app *Secret) getGroupSecrets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.rest.MethodNotAllowed(w, r, "GET")
		return
	}
	var input struct {
		GroupID int64 `json:"group_id"`
	}
	// Parse request
	if err := app.rest.ReadJSON(w, r, "secrets.getGroupSecrets", &input); err != nil {
		app.rest.Error(w, err)
		return
	}
	// Validate parameters
	v := validator.New()
	v.Check(input.GroupID > 0, "group_id", "must be provided")
	if err := v.Valid("secrets.getGroupSecrets"); err != nil {
		app.rest.Error(w, err)
		return
	}
	user := middleware.ContextGetUser(r)
	// check for permission
	exits, err := app.group.IsUserInGroup(input.GroupID, user.ID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	if !exits {
		app.rest.WriteJSON(w, "secrets.getGroupSecrets", http.StatusUnauthorized, rest.Envelope{
			"message": "Only group members can access secrets.",
		})
		return
	}
	data, err := app.secrets.GetByGroupID(input.GroupID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	app.rest.WriteJSON(w, "secrets.getGroupSecrets", http.StatusOK, rest.Envelope{
		"message": "Success!",
		"data":    data,
	})
}
