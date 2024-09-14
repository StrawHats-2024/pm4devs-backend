package main

import (
	"database/sql"
	"log"
	"pm4devs-backend/cmd/api"
	"pm4devs-backend/db"
)

func main() {
	log.Println("Starting up server")
	database, err := db.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	initStorage(database)

	server := api.NewAPIServer(":3000", database)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}
