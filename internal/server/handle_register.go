package server

import (
	"net/http"
)

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, s.APIEndpoints.UserCollection)
}
