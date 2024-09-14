package sharing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pm4devs-backend/services/auth"
	"pm4devs-backend/services/secrets"
	"pm4devs-backend/types"
	"pm4devs-backend/utils"
	"testing"
)

// Test case for successfully sharing a secret with another user
func TestShareSecretSuccess(t *testing.T) {
	// Set up a test server with a valid request
	reqBody := map[string]string{
		"user_email":  "user3@example.com",
		"permissions": "read", // or 'write'
	}
	body, _ := json.Marshal(reqBody)

	// Create a new request without an authorization header, simulating cookie-based auth
	req, err := http.NewRequest("POST", "/secrets/1/share", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	token, err := auth.CreateJWT(&types.User{UserID: 1})
	if err != nil {
		t.Fatal(err)
	}
	// Set a cookie with a valid JWT token
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})

	testHandler := getTestHandler(t)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(utils.MakeHTTPHandleFunc(testHandler.handleShareToUser))

	// Call the handler with the valid request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response["message"] == nil || response["message"] != "Secret shared successfully" {
		t.Error("Expected success message to be present in response")
	}
}

// Test case for sharing a secret with a non-existent user
func TestShareSecretUserNotFound(t *testing.T) {
	// Set up a test server with an invalid request
	reqBody := map[string]string{
		"user_email":  "nonexistentuser@example.com",
		"permissions": "read", // or 'write'
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "/secrets/1/share", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set a cookie with a valid JWT token
	req.AddCookie(&http.Cookie{
		Name:  "jwt",
		Value: "valid_jwt_token", // Replace with a method to generate a valid token
	})

	testHandler := getTestHandler(t)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(utils.MakeHTTPHandleFunc(testHandler.handleShareToUser))

	// Call the handler with the request for a non-existent user
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func getTestHandler(t *testing.T) *Handler {
	db := utils.SetupTestPostgres(t)
	defer utils.TeardownTestPostgres(t, db)
	testStore := secrets.NewStore(db)
	testHandler := NewHandler(testStore)
	return testHandler
}
