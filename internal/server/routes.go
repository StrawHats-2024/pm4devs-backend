package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/api/health", s.HelloWorldHandler)
	r.HandleFunc("/api/auth/register", s.handleRegister).Methods(http.MethodPost)
	r.HandleFunc("/api/auth/login", s.handleLogin).Methods(http.MethodPost)
	r.HandleFunc("/api/auth/refresh-token", s.handleRefreshToken).Methods(http.MethodPost)
	r.HandleFunc("/api/secrets/user", s.handleGetUserSecrets).Methods(http.MethodGet)

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	return loggedRouter
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
