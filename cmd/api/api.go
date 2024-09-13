package api

import (
	"database/sql"
	"log"
	"net/http"
	"pm4devs-backend/services/groups"
	"pm4devs-backend/services/secrets"
	"pm4devs-backend/services/user"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(listenAddr string, store *sql.DB) *APIServer {
	return &APIServer{
		addr: listenAddr,
		db:   store,
	}
}
func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	secretsStore := secrets.NewStore(s.db)
	secretsHandler := secrets.NewHandler(secretsStore)
	secretsHandler.RegisterRoutes(subrouter)

	groupsStore := groups.NewStore(s.db)
	groupsHandler := groups.NewHandler(groupsStore)
	groupsHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
