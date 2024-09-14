package sharing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pm4devs-backend/services/auth"
	"pm4devs-backend/types"
	"pm4devs-backend/utils"
)

func (s *Handler) handleRevokeAccessToUser(w http.ResponseWriter, r *http.Request) error {
	var req types.RevokeSecretAccessPayload // Define this struct as per your needs
	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return nil
	}

	// Validate the request body
	err = utils.ValidateRequestBody(req, w)
	if err != nil {
		return err
	}

	// Get user ID from the request
	userId, err := auth.GetUserIdfromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}

	// Get the secret ID from the URL
	secretId, err := getSecretIdFromUrl(r)
	if err != nil {
		http.Error(w, "Invalid secret ID", http.StatusBadRequest)
		return err
	}

	// Check the user's permission for the secret
	permissionType, err := s.store.GetUserSecretPermission(userId, secretId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking permissions: %v", err.Error()),
			http.StatusBadRequest)
		return nil
	}

	// Check if the user is allowed to revoke access
	if permissionType == types.NotAllowed {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return nil
	}

	// Revoke the secret access from the user
	err = s.store.RevokeSharingFromUser(secretId, req.UserEmail) // Assuming this method exists in your store
	if err != nil {
		return err
	}

	// Respond with success
	return utils.WriteJSON(w, http.StatusOK,
		map[string]string{"message": "Access revoked successfully"})
}
