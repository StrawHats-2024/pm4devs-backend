package sharing

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

func (s *Handler) handleShareToUser(w http.ResponseWriter, r *http.Request) error {
	var req types.ShareSecretToUserPayload
	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	// Validate the request body
	err = utils.ValidateRequestBody(req, w)
	if err != nil {
		return err
	}

	// Get user ID from the request
	userId, err := auth.GetUserIdfromRequest(r)
	if err != nil {
		return err
	}

	// Get the secret ID from the URL
	secretId, err := getSecretIdFromUrl(r)
	fmt.Println("secretId: ", secretId)
	if err != nil {
		return err
	}

	// Check the user's permission for the secret
	permissionType, err := s.store.GetUserSecretPermission(userId, secretId)
	if err != nil {
		return err
	}

	// Check if the user is allowed to share the secret
	if permissionType == types.NotAllowed {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return nil
	}

	// Validate the permission request
	if permissionType == types.ReadOnly && req.Permissions == types.WriteRead {
		return utils.WriteJSON(w, http.StatusForbidden, utils.ApiError{Error: "You don't have correct permissions for this action."})
	}

	// Share the secret with the user
	err = s.store.ShareSecretWithUser(secretId, req.UserEmail, req.Permissions)
	if err != nil {
		return err
	}

	// Respond with success
	return utils.WriteJSON(w, http.StatusOK,
		map[string]string{"message": "Secret shared with user successfully"})
}

func getSecretIdFromUrl(r *http.Request) (int, error) {

	vars := mux.Vars(r)
	groupIdStr := vars["secret_id"]
	secretId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		return -1, err
	}
	return secretId, nil
}
