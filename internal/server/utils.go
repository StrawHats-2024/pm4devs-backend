package server

import (
	"fmt"
	"io"
	"net/http"
)

func ForwardRequest(w http.ResponseWriter, r *http.Request, url string) {
	// Create a new request with the same method, URL, and body as the original request
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		fmt.Printf("Failed to create new request: %v\n", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Copy all headers from the original request to the new request
	for header, values := range r.Header {
		for _, value := range values {
			req.Header.Add(header, value)
		}
	}

	// Create an HTTP client and send the new request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response headers from the new request's response to the original response
	for header, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}

	// Set the response status code
	w.WriteHeader(resp.StatusCode)

	// Copy the response body to the original response
	if _, err := io.Copy(w, resp.Body); err != nil {
		fmt.Printf("Error while copying response body: %v\n", err)
		http.Error(w, "Failed to copy response body", http.StatusInternalServerError)
		return
	}
}
