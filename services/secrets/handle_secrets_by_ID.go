package secrets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pm4devs-backend/services/auth"
	"pm4devs-backend/types"
	"pm4devs-backend/utils"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) handleSecretsManagementById(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	secretIdStr := vars["secret_id"]
	fmt.Println("secretId: ", secretIdStr)
	secretId, err := strconv.Atoi(secretIdStr)
	if err != nil {
		return err
	}

	userId, err := auth.GetUserIdfromRequest(r)
	if err != nil {
		return err
	}

	switch r.Method {
	case http.MethodGet:
		// only owner can get
		secret, err := h.store.GetSecretById(secretId)
		if err != nil {
			return err
		}
		if secret.UserID != userId {
			return utils.WriteJSON(w, http.StatusUnauthorized,
				utils.ApiError{Error: "Unauthorized access: You do not have permission to view this secret."})
		}
		return utils.WriteJSON(w, http.StatusOK, secret)

	case http.MethodDelete:
		// only owner can delete
		secret, err := h.store.GetSecretById(secretId)
		if err != nil {
			return err
		}
		if secret.UserID != userId {
			return utils.WriteJSON(w, http.StatusUnauthorized,
				utils.ApiError{Error: "Unauthorized access: You do not have permission to delete this secret."})
		}
		err = h.store.DeleteSecretById(secretId)
		if err != nil {
			return err
		}
		return utils.WriteJSON(w, http.StatusOK,
			struct {
				Message string `json:"message"`
			}{Message: "Secret deleted successfully"})

	case http.MethodPut:
		// only owner can update
		secret, err := h.store.GetSecretById(secretId)
		if err != nil {
			return err
		}
		if secret.UserID != userId {
			return utils.WriteJSON(w, http.StatusUnauthorized,
				utils.ApiError{Error: "Unauthorized access: You do not have permission to update this secret."})
		}
		var reqObj types.UpdateSecretPayload
		err = json.NewDecoder(r.Body).Decode(&reqObj)
		if err != nil {
			return err
		}
		err = h.store.UpdateSecretById(secretId, reqObj)
		if err != nil {
			return err
		}
		return utils.WriteJSON(w, http.StatusOK,
			struct {
				Message string `json:"message"`
			}{Message: "Secret updated successfully"})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed %s", r.Method)

	}
}
