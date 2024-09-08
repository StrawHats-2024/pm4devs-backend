package main

import (
	"fmt"
	"log"
	"net/http"
	"pm4devs-backend/internals/auth"
	"pm4devs-backend/pkg/db"
)

// import "net/http"
func main() {
	db, err := db.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("db: ", db)
	http.HandleFunc("/reg", auth.HandleUserReg)
	// http.HandleFunc("/home", Home)
	// http.HandleFunc("/refresh", Refresh)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

// type APIServer struct {
// 	listenAddr string
// }
//
// func NewAPIServer(listenAddr string) *APIServer {
// 	return &APIServer{listenAddr: listenAddr}
// }
//
// func (s *APIServer) handleUsers(w http.ResponseWriter, r http.Request) error {
//   return nil
// }
//
// func (s *APIServer) handleCreateUsers(w http.ResponseWriter, r http.Request) error {
//   return nil
// }
//
// func (s *APIServer) handleGetUsers(w http.ResponseWriter, r http.Request) error {
//   return nil
// }
