package api

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"pm4devs-backend/services/groups"
	"pm4devs-backend/services/secrets"
	"pm4devs-backend/services/sharing"
	"pm4devs-backend/services/user"

	"github.com/gorilla/handlers"
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

	sharingHandler := sharing.NewHandler(secretsStore)
	sharingHandler.RegisterRoutes(subrouter)

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins, modify as necessary
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	logger := getLogger("api_server.log")
	loggingRouter := handlers.LoggingHandler(logger, router)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, corsHandler(loggingRouter))
}

func getLogger(fileName string) io.Writer {
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	log.SetOutput(multiWriter)
	return multiWriter
}
