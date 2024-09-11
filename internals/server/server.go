package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pm4devs-backend/pkg/db"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      db.Storage
}

func NewAPIServer(listenAddr string, store db.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, "working fine")
	})

	router.HandleFunc("/auth/register", makeHTTPHandleFunc(s.handleRegister))
	router.HandleFunc("/auth/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/auth/logout", makeHTTPHandleFunc(s.handleLogout))
	router.HandleFunc("/auth/refresh", makeHTTPHandleFunc(s.handleTokenRefresh))

  router.HandleFunc("/secrets", makeHTTPHandleFunc(s.handleSecretsManagement))

	router.HandleFunc("/get/users", withAuth(makeHTTPHandleFunc(s.handleGetAllUsers)))
	//
	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func withAuth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the token is missing, return an unauthorized status
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}
			// For any other error, return a bad request status
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		tokenString := cookie.Value
		claims, err := validateToken(tokenString)
		fmt.Println("claims: ", claims.UserId)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, err.Error())
			return
		}
		f(w, r)
	}

}

func (s *APIServer) handleSecretsManagement(w http.ResponseWriter, r *http.Request) error {
  // TODO: Impliment Create secret POST
  // TODO: Impliment Get all secret GET:
  // TODO: Impliment Get One secret by ID GET:
  // TODO: Impliment update One secret by ID GET:
  // TODO: Impliment Delete one secret by ID DEL:
  return nil
}

func (s *APIServer) handleGetAllUsers(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed %s", r.Method)
	}
	//TODO: Exclude password hash from response
	users, err := s.store.GetAllUsers()
	if err != nil {
		return err
	}
	err = WriteJSON(w, http.StatusOK, users)
	if err != nil {
		return err
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type ApiError struct {
	Error string `json:"error"`
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
