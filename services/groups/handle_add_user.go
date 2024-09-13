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

func (h *Handler) handleAddUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed %s", r.Method)
	}
	vars := mux.Vars(r)
	groupIdStr := vars["group_id"]
	fmt.Println("groupIdStr: ", groupIdStr)
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		return err
	}

	userId, err := auth.GetUserIdfromRequest(r)
	if err != nil {
		return err
	}

	// Check if user is admin in the group
	isAdmin, err := h.store.IsUserAdminInGroup(userId, groupId)
	if err != nil {
		return err
	}

	if !isAdmin {
		return utils.WriteJSON(w, http.StatusForbidden, utils.ApiError{Error: "You need to an admin to add users"})
	}

	// Check if user is the creator of the group
	isCreator, err := h.store.IsGroupCreator(userId, groupId)
	if err != nil {
		return err
	}

	var req types.AddUserToGroupPayload
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	if req.UserEmail == "" {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: "User email is required."})
	}
	if req.Role != types.Admin && req.Role != types.Member {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: "Invalid Role: Only 'member' & 'admin' are allow roles"})
	}
	if !isCreator && req.Role == types.Admin {
		return utils.WriteJSON(w, http.StatusForbidden, utils.ApiError{Error: "You need to be creator of group to add admin's"})
	}
	err = h.store.AddUserToGroup(groupId, req)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{Message: "User added to the group"})
}
