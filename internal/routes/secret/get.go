package secret

import (
	"net/http"

	"pm4devs.strawhats/internal/rest"
	"pm4devs.strawhats/internal/routes/middleware"
	"pm4devs.strawhats/internal/validator"
)

const GetUserSecretsRoute = "/v1/secrets/user"

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

const GetGroupSecretsRoute = "/v1/secrets/group"

func (app *Secret) getGroupSecrets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.rest.MethodNotAllowed(w, r, "GET")
		return
	}
	var input struct {
		GroupName string `json:"group_name"`
	}
	// Parse request
	if err := app.rest.ReadJSON(w, r, "secrets.getGroupSecrets", &input); err != nil {
		app.rest.Error(w, err)
		return
	}
	// Validate parameters
	v := validator.New()
	v.Check(len(input.GroupName) > 0, "group_id", "must be provided")
	if err := v.Valid("secrets.getGroupSecrets"); err != nil {
		app.rest.Error(w, err)
		return
	}
	user := middleware.ContextGetUser(r)
	group, err := app.group.GetGroupUsers(input.GroupName)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	// check for permission
	exits, err := app.group.IsUserInGroup(group.ID, user.ID)
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
	data, err := app.secrets.GetByGroupID(group.ID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	app.rest.WriteJSON(w, "secrets.getGroupSecrets", http.StatusOK, rest.Envelope{
		"message": "Success!",
		"data":    data,
	})
}

const GetSecretsSharedToUser = "/v1/secrets/sharedto/user"

func (app *Secret) getSharedToUserSecrets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user := middleware.ContextGetUser(r)
	userSecrets, err := app.secrets.GetSecretsSharedWithUser(user.ID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	app.rest.WriteJSON(w, "secret.createNew", http.StatusOK, rest.Envelope{
		"message": "Success!",
		"data":    userSecrets,
	})
}

const GetSecretsSharedByUser = "/v1/secrets/sharedby/user"

func (app *Secret) getSharedByUserSecrets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user := middleware.ContextGetUser(r)
	shared, err := app.secrets.GetSecretsSharedToOtherUsers(user.ID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	app.rest.WriteJSON(w, "secret.createNew", http.StatusOK, rest.Envelope{
		"message": "Success!",
		"data":    shared,
	})
}

const GetSecretsSharedToGroup = "/v1/secrets/sharedto/group"


func (app *Secret) getSharedToGroupSecrets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user := middleware.ContextGetUser(r)
	shared, err := app.secrets.GetSecretsSharedToGroups(user.ID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	app.rest.WriteJSON(w, "secret.createNew", http.StatusOK, rest.Envelope{
		"message": "Success!",
		"data":    shared,
	})
}
