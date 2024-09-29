package server

import (
	"fmt"
	"io"
	"net/http"
)

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	resp, err := MakeRequest(http.MethodPost, s.APIEndpoints.UserCollection, r.Body)
	if err != nil {
		fmt.Printf("Error while making request: %v\n", err)
		http.Error(w, "Failed to make request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		fmt.Printf("Error while copying body: %v\n", err)
		http.Error(w, "Failed to copy response body", http.StatusInternalServerError)
		return
	}
}
