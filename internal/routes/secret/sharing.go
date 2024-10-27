package secret

import (
	"fmt"
	"net/http"

	"pm4devs.strawhats/internal/models/secrets"
	"pm4devs.strawhats/internal/rest"
	"pm4devs.strawhats/internal/routes/middleware"
	"pm4devs.strawhats/internal/validator"
)

const SecretShareGroupRoute = "/v1/secrets/share/group"

func (app *Secret) handleShareToGroup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		app.shareToGroup(w, r)
	case http.MethodPatch:
		app.updateGroupPermission(w, r)
	case http.MethodDelete:
		app.revokeGroupPermission(w, r)
	default:
		app.rest.MethodNotAllowed(w, r, "POST, PATCH, DELETE")
	}
}

const SecretShareUserRoute = "/v1/secrets/share/user"

func (app *Secret) handleShareToUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		app.shareToUser(w, r)
	case http.MethodPatch:
		app.updateUserPermission(w, r)
	case http.MethodDelete:
		app.revokeUserPermission(w, r)
	default:
		app.rest.MethodNotAllowed(w, r, "POST, PATCH, DELETE")
	}
}

func (app *Secret) shareToUser(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Set the response type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Define input structure
	var input struct {
		SecretID   int64              `json:"secret_id"`
		UserID     int64              `json:"user_id"`
		Permission secrets.Permission `json:"permission"`
	}

	// Parse the request
	if err := app.rest.ReadJSON(w, r, "secrets.shareToUser", &input); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Validate input
	v := validator.New()
	v.Check(input.SecretID > 0, "secret_id", "must be provided")
	v.Check(input.UserID > 0, "user_id", "must be provided")
	v.Check(input.Permission == "read-only" || input.Permission == "read-write", "permission", "must be 'read-only' or 'read-write'")
	if err := v.Valid("secrets.shareToUser"); err != nil {
		app.rest.Error(w, err)
		return
	}

	err := app.validateSecretOwnership(w, r, input.SecretID)
	if err != nil {
		return
	}

	// Call the method to share the secret with the user
	if err := app.secrets.ShareToUser(input.SecretID, input.UserID, input.Permission); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Respond with success
	app.rest.WriteJSON(w, "secret.shareToUser", http.StatusCreated, rest.Envelope{
		"message": "Secret shared successfully with the user.",
	})
}

func (app *Secret) shareToGroup(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Set the response type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Define input structure
	var input struct {
		SecretID   int64              `json:"secret_id"`
		GroupName  string             `json:"group_name"`
		Permission secrets.Permission `json:"permission"`
	}

	// Parse the request
	if err := app.rest.ReadJSON(w, r, "secrets.shareToGroup", &input); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Validate input
	v := validator.New()
	v.Check(input.SecretID > 0, "secret_id", "must be provided")
	v.Check(len(input.GroupName) > 0, "group_name", "must be provided")
	v.Check(input.Permission == "read-only" || input.Permission == "read-write",
		"permission", "must be 'read-only' or 'read-write'")
	if err := v.Valid("secrets.shareToGroup"); err != nil {
		app.rest.Error(w, err)
		return
	}

	err := app.validateSecretOwnership(w, r, input.SecretID)
	if err != nil {
		app.rest.WriteJSON(w, "secret.shareToGroup", http.StatusUnauthorized, rest.Envelope{
			"message": "Valied to validate ownership",
		})
		return
	}
	group, err2 := app.group.GetByGroupName(input.GroupName)
	if err2 != nil {
		app.rest.Error(w, err2)
	}

	// Call the method to share the secret with the group
	if err := app.secrets.ShareToGroup(input.SecretID, group.ID, input.Permission); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Respond with success
	app.rest.WriteJSON(w, "secret.shareToGroup", http.StatusCreated, rest.Envelope{
		"message": "Secret shared successfully with the group.",
	})
}

func (app *Secret) updateGroupPermission(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is PATCH
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Set the response type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Define input structure
	var input struct {
		SecretID   int64              `json:"secret_id"`
		GroupName  string             `json:"group_name"`
		Permission secrets.Permission `json:"permission"`
	}

	// Parse the request
	if err := app.rest.ReadJSON(w, r, "secrets.updateGroupPermission", &input); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Validate input
	v := validator.New()
	v.Check(input.SecretID > 0, "secret_id", "must be provided")
	v.Check(len(input.GroupName) > 0, "group_name", "must be provided")
	v.Check(input.Permission == "read-only" || input.Permission == "read-write", "permission", "must be 'read-only' or 'read-write'")
	if err := v.Valid("secrets.updateGroupPermission"); err != nil {
		app.rest.Error(w, err)
		return
	}
	err := app.validateSecretOwnership(w, r, input.SecretID)
	if err != nil {
		return
	}

	group, err2 := app.group.GetByGroupName(input.GroupName)
	if err2 != nil {
		app.rest.Error(w, err2)
	}

	// Call the method to update the permission
	if err := app.secrets.UpdateGroupPermission(input.SecretID, group.ID, input.Permission); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Respond with success
	app.rest.WriteJSON(w, "secret.updateGroupPermission", http.StatusOK, rest.Envelope{
		"message": "Permission updated successfully for the group.",
	})
}

func (app *Secret) updateUserPermission(w http.ResponseWriter, r *http.Request) {
	// Set the response type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Define input structure
	var input struct {
		SecretID   int64              `json:"secret_id"`
		UserEmail  string             `json:"user_email"`
		Permission secrets.Permission `json:"permission"`
	}

	// Parse the request
	if err := app.rest.ReadJSON(w, r, "secrets.updateUserPermission", &input); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Validate input
	v := validator.New()
	v.Check(input.SecretID > 0, "secret_id", "must be provided")
	v.Check(len(input.UserEmail) > 0, "user_email", "must be provided")
	v.Check(input.Permission == "read-only" || input.Permission == "read-write",
		"permission", "must be 'read-only' or 'read-write'")
	if err := v.Valid("secrets.updateUserPermission"); err != nil {
		app.rest.Error(w, err)
		return
	}
	err := app.validateSecretOwnership(w, r, input.SecretID)
	if err != nil {
		return
	}

	user, err2 := app.users.GetByEmail(input.UserEmail)
	if err2 != nil {
		app.rest.Error(w, err2)
	}

	// Call the method to update the permission
	if err := app.secrets.UpdateUserPermission(input.SecretID, user.ID, input.Permission); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Respond with success
	app.rest.WriteJSON(w, "secret.updateUserPermission", http.StatusOK, rest.Envelope{
		"message": "Permission updated successfully for the user.",
	})
}

func (app *Secret) revokeGroupPermission(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is DELETE
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Set the response type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Define input structure
	var input struct {
		SecretID  int64  `json:"secret_id"`
		GroupName string `json:"group_name"`
	}

	// Parse the request
	if err := app.rest.ReadJSON(w, r, "secrets.revokeGroupPermission", &input); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Validate input
	v := validator.New()
	v.Check(input.SecretID > 0, "secret_id", "must be provided")
	v.Check(len(input.GroupName) > 0, "group_id", "must be provided")
	if err := v.Valid("secrets.revokeGroupPermission"); err != nil {
		app.rest.Error(w, err)
		return
	}

	err := app.validateSecretOwnership(w, r, input.SecretID)
	if err != nil {
		return
	}

	group, err2 := app.group.GetByGroupName(input.GroupName)
	if err2 != nil {
		app.rest.Error(w, err2)
	}

	// Call the method to revoke the permission
	if err := app.secrets.RevokeFromGroup(input.SecretID, group.ID); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Respond with success
	app.rest.WriteJSON(w, "secret.revokeGroupPermission", http.StatusOK, rest.Envelope{
		"message": "Permission revoked successfully for the group.",
	})
}

func (app *Secret) revokeUserPermission(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is DELETE
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Set the response type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Define input structure
	var input struct {
		SecretID  int64  `json:"secret_id"`
		UserEmail string `json:"user_email"`
	}

	// Parse the request
	if err := app.rest.ReadJSON(w, r, "secrets.revokeUserPermission", &input); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Validate input
	v := validator.New()
	v.Check(input.SecretID > 0, "secret_id", "must be provided")
	v.Check(len(input.UserEmail) > 0, "user_email", "must be provided")
	if err := v.Valid("secrets.revokeUserPermission"); err != nil {
		app.rest.Error(w, err)
		return
	}
	err := app.validateSecretOwnership(w, r, input.SecretID)
	if err != nil {
		return
	}
	user, err2 := app.users.GetByEmail(input.UserEmail)
	if err2 != nil {
		app.rest.Error(w, err2)
	}

	// Call the method to revoke the permission
	if err := app.secrets.RevokeFromUser(input.SecretID, user.ID); err != nil {
		app.rest.Error(w, err)
		return
	}

	// Respond with success
	app.rest.WriteJSON(w, "secret.revokeUserPermission", http.StatusOK, rest.Envelope{
		"message": "Permission revoked successfully for the user.",
	})
}

func (app *Secret) validateSecretOwnership(w http.ResponseWriter, r *http.Request, secretID int64) error {

	currUser := middleware.ContextGetUser(r)
	currsecret, err := app.secrets.GetSecretByID(secretID)
	if err != nil {
		app.rest.Error(w, err)
		return fmt.Errorf("error")
	}
	if currUser.ID != currsecret.OwnerID {
		app.rest.WriteJSON(w, "secrets.validateSecretOwnership", http.StatusUnauthorized, rest.Envelope{
			"message": "Only secret owner can manage access",
		})
		return fmt.Errorf("error")
	}
	return nil
}
