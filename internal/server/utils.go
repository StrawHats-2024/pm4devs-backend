package server

import (
	"fmt"
	"io"
	"net/http"
)

func MakeRequest(method string, url string, body io.ReadCloser, originalReq *http.Request) (*http.Response, error) {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the content type to application/json
	req.Header.Set("Content-Type", "application/json")

	// Copy the Authorization header from the original request if it exists
	authHeader := originalReq.Header.Get("Authorization")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	fmt.Println("resp status: ", resp.StatusCode)

	return resp, nil
}

func ForwardRequest(w http.ResponseWriter, r *http.Request, url string) {
	// Forward the request with the Authorization header
	resp, err := MakeRequest(r.Method, url, r.Body, r)
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
