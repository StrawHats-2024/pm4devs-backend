package main

import (
	"fmt"
	"os"
	"pm4devs/internal/server"
)

func main() {

	server := server.NewServer()

	fmt.Printf("Listening at port: %s\n", os.Getenv("PORT"))
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
