package server

import (
	"fmt"
	"net/http"
)

func ForwardHandler(targetURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the target URL
		fmt.Println("Forward handler hit")

	}
}
