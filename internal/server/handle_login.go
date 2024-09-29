package server

import (
	"net/http"
)

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, s.APIEndpoints.AuthWithPassword)
}
