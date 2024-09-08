package main

import (
	// "fmt"
	"log"
	// "pm4devs-backend/internals/server"
	"pm4devs-backend/pkg/db"
	"pm4devs-backend/pkg/models"
	"time"
)

// import "net/http"
func main() {
	store, err := db.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	// server := server.NewAPIServer(":5000", store)
	err = store.CreateUser(models.User{Email: "testing", PasswordHash: "slfjsdlfjsldkfjslfjsdl", CreatedAt: time.Now()})
	if err != nil {
		panic(err)
	}
	// fmt.Println("server: ", server)
	// http.HandleFunc("/reg", auth.HandleUserReg)
	// http.HandleFunc("/home", Home)
	// http.HandleFunc("/refresh", Refresh)

	// log.Fatal(http.ListenAndServe(":8080", nil))

}
