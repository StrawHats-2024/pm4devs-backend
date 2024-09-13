package groups

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pm4devs-backend/services/auth"
	"pm4devs-backend/types"
	"pm4devs-backend/utils"
)

func (h *Handler) handleGroupManagement(w http.ResponseWriter, r *http.Request) error {
	userId, err := auth.GetUserIdfromRequest(r)
	if err != nil {
		return err
	}
	switch r.Method {

	case http.MethodGet:
		groups, err := h.store.GetUserGroups(userId)
		if err != nil {
			return err
		}
		return utils.WriteJSON(w, http.StatusOK, groups)

	case http.MethodPost:
		req := struct {
			GroupName string `json:"group_name"`
		}{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			return err
		}
		if req.GroupName == "" {
			return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: "Group name is required"})
		}
		groupId, err := h.store.CreateGroup(&types.Group{
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
		return utils.WriteJSON(w, http.StatusOK, res)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed %s", r.Method)
	}
}
