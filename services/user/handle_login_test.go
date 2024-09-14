package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pm4devs-backend/types"
	"pm4devs-backend/utils"
	"testing"
)

// Handler for user login

// Test function for successful login
func TestLoginUserSuccess(t *testing.T) {
	// Set up a test server
	reqBody := map[string]string{
		"email":    "usernew@example.com",
		"password": "user_password",
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	testHandler := getTestHandler(t)
	passwordHash, err := utils.HashPassword(reqBody["password"])
	if err != nil {
		t.Fatal(err)
	}
	_, err = testHandler.store.CreateUser(&types.User{
		Username:     "test1",
		Email:        reqBody["email"],
		PasswordHash: passwordHash,
	})
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(utils.MakeHTTPHandleFunc(testHandler.handleLogin))

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

	if response["token"] == nil {
		t.Error("Expected token to be present in response")
	}

	if response["user_id"] == nil {
		t.Error("Expected user_id to be present in response")
	}
}

// Test case for bad request
func TestLoginUserBadRequest(t *testing.T) {
	// Set up a test server with an invalid request
	invalidReqBody := map[string]string{
		"email":    "invalid_email",
		"password": "short",
	}
	body, _ := json.Marshal(invalidReqBody)

	req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	testHandler := getTestHandler(t)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(utils.MakeHTTPHandleFunc(testHandler.handleLogin))

	// Call the handler with the invalid request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized && status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

// Test case for invalid credentials
func TestLoginUserInvalidCredentials(t *testing.T) {
	// Set up a test server with invalid login credentials
	invalidReqBody := map[string]string{
		"email":    "user@example.com",
		"password": "wrong_password",
	}
	body, _ := json.Marshal(invalidReqBody)

	req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	testHandler := getTestHandler(t)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(utils.MakeHTTPHandleFunc(testHandler.handleLogin))
	// Call the handler with the invalid credentials
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized && status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}
