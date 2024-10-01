package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func assertAny(t *testing.T, want, got any) {
	if got != want {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func getAuthToken(t *testing.T) string {
	s := newTestServer()
	server := httptest.NewServer(http.HandlerFunc(s.handleLogin))
	defer server.Close()
	payload := LoginPayload{Identity: "olenharris@mclaughlin.name", Password: "parikshith"}
	resp := makePostReq(t, server.URL, getBodyJson(t, payload))
	got := resp.StatusCode
	want := http.StatusOK
	assertStatusCode(t, want, got)
	var resBody LoginResponse
	err := json.NewDecoder(resp.Body).Decode(&resBody)
	if err != nil {
		t.Fatal(err)
	}
	return resBody.Token
}

func TestGetAuthTokenHelper(t *testing.T) {
	token := getAuthToken(t)
	if token == "" {
		t.Errorf("No token")
	}
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

func TestAuthHeader(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the Authorization header
		authHeader := r.Header.Get("Authorization")
		expectedToken := "Bearer test-token"

		if authHeader != expectedToken {
			t.Errorf("Expected Authorization header to be '%s', got '%s'", expectedToken, authHeader)
		}

		// Respond with a success message
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	// Make the request to the test server
	makeNewReq(t, server.URL, "GET", nil, "test-token")
}

func makeNewReq(t *testing.T, url string, method string, body io.Reader, token string) *http.Response {
	t.Helper()

	// Create a new HTTP request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}
