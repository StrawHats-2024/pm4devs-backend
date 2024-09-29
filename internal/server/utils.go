package server

import (
	"fmt"
	"io"
	"net/http"
)

func MakeRequest(method string, url string, body io.ReadCloser) (*http.Response, error) {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the content type to application/json
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println("resp status: ", resp.StatusCode)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	return resp, nil
}
