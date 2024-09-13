package groups

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

func (s *Handler) handleRemoveUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	vars := mux.Vars(r)
	groupIdStr := vars["group_id"]
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		return err
	}

	userId, err := auth.GetUserIdfromRequest(r)
	if err != nil {
		return err
	}

	// Check if user is admin in the group
	isAdmin, err := s.store.IsUserAdminInGroup(userId, groupId)
	if err != nil {
		return err
	}

	// Check if user is the creator of the group
	isCreator, err := s.store.IsGroupCreator(userId, groupId)
	if err != nil {
		return err
	}

	req := struct {
		UserEmail string `json:"user_email"`
	}{}

	if !isAdmin {
		return utils.WriteJSON(w, http.StatusForbidden, utils.ApiError{Error: "You need to an admin to delete users"})
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	if req.UserEmail == "" {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: "User email is required."})
	}

	// Check if the user to be deleted is a member of the group and get their role
	userRole, err := s.store.GetUserRoleInGroup(req.UserEmail, groupId)
	if err != nil {
		return err
	}

	// If user not found in group, return an error
	if userRole == "" {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: "User not found in the group."})
	}

	// If the current user is an admin, they can only delete members
	// If the current user is the creator, they can delete both members and admins
	if !isCreator && userRole == types.Admin {
		return utils.WriteJSON(w, http.StatusForbidden, utils.ApiError{Error: "Admin cannot remove another admin."})
	}

	// Proceed to delete the user from the group
	err = s.store.DeleteUserFromGroup(groupId, req.UserEmail)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{Message: "User removed from the group"})
}
