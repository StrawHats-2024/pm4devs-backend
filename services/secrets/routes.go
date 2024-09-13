package secrets

import (
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

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/secrets",
		auth.WithAuth(utils.MakeHTTPHandleFunc(h.handleSecretsManagement)))
	router.HandleFunc("/secrets/{secret_id}",
		auth.WithAuth(utils.MakeHTTPHandleFunc(h.handleSecretsManagementById)))
}
