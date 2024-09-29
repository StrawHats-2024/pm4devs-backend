package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port         int
	APIEndpoints *APIEndpoints
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	var BASE_URL = os.Getenv("POCKETBASE_URL")
	if BASE_URL == "" {
		log.Fatal("POCKETBASE_URL not found")
	}
	NewServer := &Server{
		port:         port,
		APIEndpoints: NewAPIEndpoints(BASE_URL),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
