package main

import (
	"fmt"
	"log"

	// "pm4devs-backend/internals/server"
	"pm4devs-backend/internals/server"
	"pm4devs-backend/pkg/db"
)

func main() {
	store, err := db.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	// user, err := server.NewUser("testing2", "parikshith")
	// if err != nil {
	// 	panic(err)
	// }
	// userid, err := store.CreateUser(user)
	//  fmt.Println("userid: ", userid);
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println("Starting....")
	server := server.NewAPIServer(":3000", store)
	server.Run()
}
