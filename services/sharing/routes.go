package sharing

import (
	"net/http"
	"pm4devs-backend/services/auth"
	"pm4devs-backend/types"
	"pm4devs-backend/utils"

	"github.com/gorilla/mux"
)

type Handler struct {
	store types.SecretStore
}

func NewHandler(store types.SecretStore) *Handler {
	return &Handler{store: store}
}

func (s *Handler) RegisterRoutes(router *mux.Router) {

	router.Handle("/secrets/{secret_id}/share",
		auth.WithAuth(utils.MakeHTTPHandleFunc(s.handleShareToUser))).Methods(http.MethodPost)
	router.Handle("/secrets/{secret_id}/share",
		auth.WithAuth(utils.MakeHTTPHandleFunc(s.handleRevokeAccessToUser))).Methods(http.MethodDelete)
	router.Handle("/secrets/{secret_id}/share/group",
		auth.WithAuth(utils.MakeHTTPHandleFunc(s.handleShareToGroup))).Methods(http.MethodPost)
}
