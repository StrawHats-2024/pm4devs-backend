package groups

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pm4devs-backend/services/auth"
	"pm4devs-backend/utils"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) handleGroupManagementWithId(w http.ResponseWriter, r *http.Request) error {
	userId, err := auth.GetUserIdfromRequest(r)
	vars := mux.Vars(r)
	groupIdStr := vars["group_id"]
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		return err
	}
	switch r.Method {
	case http.MethodGet:
		group, err := h.store.GetGroupById(groupId)
		fmt.Println("group with users: ", group)
		if err != nil {
			return err
		}
		return utils.WriteJSON(w, http.StatusOK, group)

	case http.MethodPut:
		req := struct {
			NewGroupName string `json:"new_group_name"`
		}{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			return err
		}
		if req.NewGroupName == "" {
			return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: "Group new name is required"})
		}
		isAdmin, err := h.store.IsUserAdminInGroup(userId, groupId)
		if err != nil {
			return err
		}
		if !isAdmin {
			return utils.WriteJSON(w, http.StatusForbidden, utils.ApiError{Error: "You need to an admin to rename group"})
		}
		err = h.store.UpdateGroupName(groupId, req.NewGroupName)
		if err != nil {
			return err
		}
		res := struct {
			Message string `json:"message"`
		}{
			Message: "Group name updated successfully",
		}
		return utils.WriteJSON(w, http.StatusOK, res)

	case http.MethodDelete:

		// Check if the current user is the creator of the group
		isCreator, err := h.store.IsGroupCreator(userId, groupId)
		if err != nil {
			return utils.WriteJSON(w, http.StatusInternalServerError, utils.ApiError{Error: err.Error()})
		}
		if !isCreator {
			return utils.WriteJSON(w, http.StatusUnauthorized, utils.ApiError{Error: "Only the group creator can delete the group"})
		}

		err = h.store.DeleteGroup(groupId)
		if err != nil {
			return utils.WriteJSON(w, http.StatusInternalServerError, utils.ApiError{Error: err.Error()})
		}

		res := struct {
			Message string `json:"message"`
		}{
			Message: "Group deleted successfully",
		}

		return utils.WriteJSON(w, http.StatusOK, res)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed %s", r.Method)
	}

}
