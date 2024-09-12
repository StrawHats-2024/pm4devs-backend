package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pm4devs-backend/pkg/db"
	"pm4devs-backend/pkg/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (s *APIServer) handleSecretsManagement(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		secrets, err := s.store.GetAllSecret()
		if err != nil {
			return err
		}
		err = WriteJSON(w, http.StatusOK, secrets)
		return nil

	case http.MethodPost:
		var secret *models.Secret
		err := json.NewDecoder(r.Body).Decode(&secret)
		if err != nil {
			return err
		}
		secret.CreatedAt = time.Now()
		secreId, err := s.store.CreateSecret(secret)
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
		err = WriteJSON(w, http.StatusOK, response)
		return err
	default:
		return fmt.Errorf("method not allowed %s", r.Method)

	}
}

func (s *APIServer) handleSecretsManagementById(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	secretIdStr := vars["secret_id"]
	fmt.Println("secretId: ", secretIdStr)
	secretId, err := strconv.Atoi(secretIdStr)
	if err != nil {
		return err
	}

	switch r.Method {
	case http.MethodGet:
		secret, err := s.store.GetSecretById(secretId)
		if err != nil {
			return err
		}
		err = WriteJSON(w, http.StatusOK, secret)
		return err

	case http.MethodDelete:
		err := s.store.DeleteSecretById(secretId)
		if err != nil {
			return err
		}
		err = WriteJSON(w, http.StatusOK,
			struct {
				Message string `json:"message"`
			}{Message: "Secret deleted successfully"})
		return err

	case http.MethodPut:
		var reqObj db.UpdateSecretReq
		err := json.NewDecoder(r.Body).Decode(&reqObj)
		if err != nil {
			return err
		}
		err = s.store.UpdateSecretById(secretId, reqObj)
		if err != nil {
			return err
		}
		err = WriteJSON(w, http.StatusOK,
			struct {
				Message string `json:"message"`
			}{Message: "Secret updated successfully"})

		return err

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed %s", r.Method)

	}
}
