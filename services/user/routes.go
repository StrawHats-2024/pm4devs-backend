package user

import (
	"pm4devs-backend/services/auth"
	"pm4devs-backend/types"
	"pm4devs-backend/utils"

	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/login", utils.MakeHTTPHandleFunc(h.handleLogin))
	router.HandleFunc("/auth/register", utils.MakeHTTPHandleFunc(h.handleRegister))
	router.HandleFunc("/auth/logout", auth.WithAuth(utils.MakeHTTPHandleFunc(h.handleLogout)))
}
