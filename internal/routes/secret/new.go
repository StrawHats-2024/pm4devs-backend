package secret

import (
	"net/http"

	"pm4devs.strawhats/internal/rest"
	"pm4devs.strawhats/internal/routes/middleware"
	"pm4devs.strawhats/internal/validator"
)

const CreateNewSecretRoute = "/v1/secrets"

func (app *Secret) createNew(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var input struct {
		Name          string `json:"name"`
		EncryptedData string `json:"encrypted_data"`
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
	if err := v.Valid("secrets.createNew"); err != nil {
		app.rest.Error(w, err)
		return
	}

	user := middleware.ContextGetUser(r)
	newSecret, err := app.secrets.NewRecord(input.Name, input.EncryptedData, user.ID)
	if err != nil {
		app.rest.Error(w, err)
	}
	app.rest.WriteJSON(w, "secret.createNew", http.StatusCreated, rest.Envelope{
		"message":   "Success! Your secret has been created.",
		"secret_id": newSecret.ID,
	})
}