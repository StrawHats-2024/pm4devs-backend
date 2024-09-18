package secrets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pm4devs-backend/services/auth"
	"pm4devs-backend/types"
	"pm4devs-backend/utils"
	"time"
)

func (h *Handler) handleSecretsManagement(w http.ResponseWriter, r *http.Request) error {
	userIdFromCookie, err := auth.GetUserIdfromRequest(r)
	fmt.Println("userIdFromCookie: ", userIdFromCookie)
	if err != nil {
		return err
	}

	switch r.Method {
	case http.MethodGet:
		//TODO: Deal with permissions
		secrets, err := h.store.GetAllSecret()
		if err != nil {
			return err
		}
		err = utils.WriteJSON(w, http.StatusOK, secrets)
		return nil

	case http.MethodPost:
		// only add with your own userid
		var secret *types.Secret
		err := json.NewDecoder(r.Body).Decode(&secret)
		if err != nil {
			return err
		}
		secret.CreatedAt = time.Now()
    secret.UserID = userIdFromCookie
		secreId, err := h.store.CreateSecret(secret)
		if err != nil {
			return err
		}
		response := struct {
			SecretID int    `json:"secret_id"` // Change to exported field
			Message  string `json:"message"`   // Change to exported field
		}{
			SecretID: secreId,
			Message:  "Secret created successfully",
		}
		fmt.Println("created secreId: ", secreId)
		err = utils.WriteJSON(w, http.StatusOK, response)
		return err

	default:
		return fmt.Errorf("method not allowed %s", r.Method)

	}
}
