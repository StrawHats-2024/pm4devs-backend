package server

import "net/http"

func (s *Server) handleGetUserSecrets(w http.ResponseWriter, r *http.Request) {
	ForwardRequest(w, r, s.APIEndpoints.SecretsCollection)
}

func (s *Server) handleGetSharedSecrets(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) handleGetGroupSecrets(w http.ResponseWriter, r *http.Request) {

}
