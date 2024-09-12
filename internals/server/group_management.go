package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pm4devs-backend/pkg/db"
	"pm4devs-backend/pkg/models"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *APIServer) handleGroupManagement(w http.ResponseWriter, r *http.Request) error {
	userId, err := getUserIdfromCookie(r)
	if err != nil {
		return err
	}
	switch r.Method {

	case http.MethodGet:
		groups, err := s.store.GetUserGroups(userId)
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, groups)

	case http.MethodPost:
		req := struct {
			GroupName string `json:"group_name"`
		}{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			return err
		}
		if req.GroupName == "" {
			return WriteJSON(w, http.StatusBadRequest, ApiError{Error: "Group name is required"})
		}
		groupId, err := s.store.CreateGroup(&models.Group{
			GroupName: req.GroupName,
			CreatedBy: userId,
		})
		if err != nil {
			return err
		}
		res := struct {
			GroupId int    `json:"group_id"`
			Message string `json:"message"`
		}{
			GroupId: groupId,
			Message: "Group creted successfully",
		}
		return WriteJSON(w, http.StatusOK, res)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed %s", r.Method)
	}
}

func (s *APIServer) handleRemoveUser(w http.ResponseWriter, r *http.Request) error {
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

	userId, err := getUserIdfromCookie(r)
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
		return WriteJSON(w, http.StatusForbidden, ApiError{Error: "You need to an admin to delete users"})
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	if req.UserEmail == "" {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: "User email is required."})
	}

	// Check if the user to be deleted is a member of the group and get their role
	userRole, err := s.store.GetUserRoleInGroup(req.UserEmail, groupId)
	if err != nil {
		return err
	}

	// If user not found in group, return an error
	if userRole == "" {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: "User not found in the group."})
	}

	// If the current user is an admin, they can only delete members
	// If the current user is the creator, they can delete both members and admins
	if !isCreator && userRole == models.Admin {
		return WriteJSON(w, http.StatusForbidden, ApiError{Error: "Admin cannot remove another admin."})
	}

	// Proceed to delete the user from the group
	err = s.store.DeleteUserFromGroup(groupId, req.UserEmail)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{Message: "User removed from the group"})
}

func (s *APIServer) handleGroupManagementWithId(w http.ResponseWriter, r *http.Request) error {
	userId, err := getUserIdfromCookie(r)
	vars := mux.Vars(r)
	groupIdStr := vars["group_id"]
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		return err
	}
	switch r.Method {
	case http.MethodGet:
		group, err := s.store.GetGroupById(groupId)
		fmt.Println("group with users: ", group)
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, group)

	case http.MethodPut:
		req := struct {
			NewGroupName string `json:"new_group_name"`
		}{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			return err
		}
		if req.NewGroupName == "" {
			return WriteJSON(w, http.StatusBadRequest, ApiError{Error: "Group new name is required"})
		}
		isAdmin, err := s.store.IsUserAdminInGroup(userId, groupId)
		if err != nil {
			return err
		}
		if !isAdmin {
			return WriteJSON(w, http.StatusForbidden, ApiError{Error: "You need to an admin to rename group"})
		}
		err = s.store.UpdateGroupName(groupId, req.NewGroupName)
		if err != nil {
			return err
		}
		res := struct {
			GroupId int    `json:"group_id"`
			Message string `json:"message"`
		}{
			Message: "Group name updated successfully",
		}
		return WriteJSON(w, http.StatusOK, res)

	case http.MethodDelete:

		// Check if the current user is the creator of the group
		isCreator, err := s.store.IsGroupCreator(userId, groupId)
		if err != nil {
			return WriteJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
		}
		if !isCreator {
			return WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Only the group creator can delete the group"})
		}

		err = s.store.DeleteGroup(groupId)
		if err != nil {
			return WriteJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
		}

		res := struct {
			Message string `json:"message"`
		}{
			Message: "Group deleted successfully",
		}

		return WriteJSON(w, http.StatusOK, res)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed %s", r.Method)
	}

}

func (s *APIServer) handleAddUser(w http.ResponseWriter, r *http.Request) error {
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

	userId, err := getUserIdfromCookie(r)
	if err != nil {
		return err
	}

	// Check if user is admin in the group
	isAdmin, err := s.store.IsUserAdminInGroup(userId, groupId)
	if err != nil {
		return err
	}

	if !isAdmin {
		return WriteJSON(w, http.StatusForbidden, ApiError{Error: "You need to an admin to add users"})
	}

	// Check if user is the creator of the group
	isCreator, err := s.store.IsGroupCreator(userId, groupId)
	if err != nil {
		return err
	}

	var req db.AddUserToGroupReq
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	if req.UserEmail == "" {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: "User email is required."})
	}
	if req.Role != models.Admin && req.Role != models.Member {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: "Invalid Role: Only 'member' & 'admin' are allow roles"})
	}
	if !isCreator && req.Role == models.Admin {
		return WriteJSON(w, http.StatusForbidden, ApiError{Error: "You need to be creator of group to add admin's"})
	}
	err = s.store.AddUserToGroup(groupId, req)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{Message: "User added to the group"})
}
