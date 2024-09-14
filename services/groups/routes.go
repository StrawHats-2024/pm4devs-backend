package groups

import (
	"net/http"
	"pm4devs-backend/services/auth"
	"pm4devs-backend/types"
	"pm4devs-backend/utils"

	"github.com/gorilla/mux"
)

type Handler struct {
	store types.GroupStore
}

func NewHandler(store types.GroupStore) *Handler {
	return &Handler{store: store}
}

func (s *Handler) RegisterRoutes(router *mux.Router) {

	router.Handle("/groups",
		auth.WithAuth(utils.MakeHTTPHandleFunc(s.handleGroupManagement)))

	router.Handle("/groups/{group_id}/add_user",
		auth.WithAuth(utils.MakeHTTPHandleFunc(s.handleAddUser))).Methods(http.MethodPost)

	router.Handle("/groups/{group_id}",
		auth.WithAuth(utils.MakeHTTPHandleFunc(s.handleGroupManagementWithId))).Methods(
		http.MethodGet, http.MethodDelete, http.MethodPut)

	router.Handle("/groups/{group_id}/remove_user",
		auth.WithAuth(utils.MakeHTTPHandleFunc(s.handleRemoveUser))).Methods(http.MethodDelete)
}
