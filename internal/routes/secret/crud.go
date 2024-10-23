package secret

import (
	"net/http"

	"pm4devs.strawhats/internal/models/secrets"
	"pm4devs.strawhats/internal/rest"
	"pm4devs.strawhats/internal/routes/middleware"
	"pm4devs.strawhats/internal/validator"
)

const SecretCRUDRoute = "/v1/secrets"

func (app *Secret) CRUDRoute(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.get(w, r)

	case http.MethodPost:
		app.createNew(w, r)

	case http.MethodPatch:
		app.update(w, r)

	case http.MethodDelete:
		app.delete(w, r)

	default:
		app.rest.MethodNotAllowed(w, r, "GET, POST, PATCH, DELETE")
	}
}

func (app *Secret) get(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SecretID int64 `json:"secret_id"`
	}
	// Parse request
	if err := app.rest.ReadJSON(w, r, "secrets.get", &input); err != nil {
		app.rest.Error(w, err)
		return
	}
	// Validate parameters
	v := validator.New()
	v.Check(input.SecretID > 0, "secret_id", "must be provided")
	if err := v.Valid("secrets.get"); err != nil {
		app.rest.Error(w, err)
		return
	}
	user := middleware.ContextGetUser(r)
	currSecret, err := app.secrets.GetSecretByID(input.SecretID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	permission, err := app.secrets.GetUserSecretPermission(user.ID, input.SecretID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	if permission == secrets.NOTALLOWED {
		app.rest.WriteJSON(w, "secrets.get", http.StatusUnauthorized, rest.Envelope{
			"message": "Your not allowed",
		})
		return
	}
	app.rest.WriteJSON(w, "secrets.get", http.StatusOK, rest.Envelope{
		"message": "Success!",
		"data":    currSecret,
	})
}

func (app *Secret) update(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SecretID      int64  `json:"secret_id"`
		Name          string `json:"name"`
		EncryptedData string `json:"encrypted_data"`
		IV            string `json:"iv"`
	}
	// Parse request
	if err := app.rest.ReadJSON(w, r, "secrets.update", &input); err != nil {
		app.rest.Error(w, err)
		return
	}
	// Validate parameters
	v := validator.New()
	v.Check(input.SecretID > 0, "secret_id", "must be provided")
	v.Check(len(input.Name) > 0, "name", "must be provided")
	v.Check(len(input.EncryptedData) > 0, "encrypted_data", "must be provided")
	v.Check(len(input.IV) > 0, "iv", "must be provided")
	if err := v.Valid("secrets.update"); err != nil {
		app.rest.Error(w, err)
		return
	}

	user := middleware.ContextGetUser(r)
	permission, err := app.secrets.GetUserSecretPermission(user.ID, input.SecretID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	if permission != secrets.ReadWrite {
		app.rest.WriteJSON(w, "secrets.update", http.StatusUnauthorized, rest.Envelope{
			"message": "Only owner can update a secret",
		})
		return
	}
	err = app.secrets.Update(input.SecretID, input.Name, input.EncryptedData, input.IV)
	if err != nil {
		app.rest.Error(w, err)
		return
	}

	app.rest.WriteJSON(w, "secrets.update", http.StatusOK, rest.Envelope{
		"message": "Success!",
	})
}

func (app *Secret) delete(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SecretID int64 `json:"secret_id"`
	}
	// Parse request
	if err := app.rest.ReadJSON(w, r, "secrets.delete", &input); err != nil {
		app.rest.Error(w, err)
		return
	}
	// Validate parameters
	v := validator.New()
	v.Check(input.SecretID > 0, "secret_id", "must be provided")
	if err := v.Valid("secrets.delete"); err != nil {
		app.rest.Error(w, err)
		return
	}
	user := middleware.ContextGetUser(r)
	currSecret, err := app.secrets.GetSecretByID(input.SecretID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	if currSecret.OwnerID != user.ID {
		app.rest.WriteJSON(w, "secrets.update", http.StatusUnauthorized, rest.Envelope{
			"message": "Only owner can delete a secret",
		})
		return
	}
	err = app.secrets.Delete(input.SecretID)
	if err != nil {
		app.rest.Error(w, err)
		return
	}
	app.rest.WriteJSON(w, "secrets.update", http.StatusNoContent, rest.Envelope{
		"message": "Success!",
	})
}

func (app *Secret) createNew(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var input struct {
		Name          string `json:"name"`
		EncryptedData string `json:"encrypted_data"`
		IV            string `json:"iv"`
	}
	// Parse request
	if err := app.rest.ReadJSON(w, r, "secrets.createNew", &input); err != nil {
		app.rest.Error(w, err)
		return
	}
	// Validate parameters
	v := validator.New()
	v.Check(len(input.Name) > 0, "name", "must be provided")
	v.Check(len(input.EncryptedData) > 0, "encrypted_data", "must be provided")
	v.Check(len(input.IV) > 0, "iv", "Initialization vector must be provided")
	if err := v.Valid("secrets.createNew"); err != nil {
		app.rest.Error(w, err)
		return
	}

	user := middleware.ContextGetUser(r)
	newSecret, err := app.secrets.NewRecord(input.Name, input.EncryptedData, input.IV, user.ID)
	if err != nil {
		app.rest.Error(w, err)
	}
	app.rest.WriteJSON(w, "secret.createNew", http.StatusCreated, rest.Envelope{
		"message":   "Success! Your secret has been created.",
		"secret_id": newSecret.ID,
	})
}
