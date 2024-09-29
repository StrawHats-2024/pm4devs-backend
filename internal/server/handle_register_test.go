package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	gofakeit "github.com/brianvoe/gofakeit/v7"
)

func TestHandleRegister(t *testing.T) {
	s := newTestServer()
	server := httptest.NewServer(http.HandlerFunc(s.handleRegister))
	defer server.Close()
	t.Run("Return 400 when body not correct", func(t *testing.T) {
		resp := makePostReq(t, server.URL, nil)
		assertStatusCode(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Get 200 on success", func(t *testing.T) {
		reqBody := CreateUserPayload{
			Email:           gofakeit.Email(),
			Password:        "parikshith",
			PasswordConfirm: "parikshith",
			Name:            gofakeit.Name(),
		}
		resp := makePostReq(t, server.URL, getBodyJson(t, reqBody))
		defer resp.Body.Close() // Make sure to close response body

		got := resp.StatusCode
		want := http.StatusOK
		assertStatusCode(t, want, got)
	})
}

func assertStatusCode(t *testing.T, want, got int) {
	assertAny(t, want, got)
}

func getBodyJson(t *testing.T, bodyObj any) *bytes.Buffer {
	t.Helper()
	jsonData, err := json.Marshal(bodyObj)
	if err != nil {
		t.Fatal("Error marshaling JSON:", err)
		return nil
	}
	return bytes.NewBuffer(jsonData)
}

func newTestServer() *Server {

	s := &Server{
		APIEndpoints: NewAPIEndpoints("http://127.0.0.1:8090/api"),
	}
	return s
}

func makePostReq(t *testing.T, url string, body io.Reader) *http.Response {
	t.Helper()
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func makeNewReq(t *testing.T, url string, method string, body io.Reader, token string) *http.Response {
	t.Helper()

	// Create a new HTTP request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Add the Authorization header with the bearer token
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}
