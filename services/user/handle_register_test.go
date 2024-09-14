package user

import (
	// "encoding/json"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pm4devs-backend/utils"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	// Set up a test server
	reqBody := map[string]string{
		"email":    "user@example.com",
		"username": "parikshith",
		"password": "user_password",
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	testHandler := getTestHandler(t)
	handler := http.HandlerFunc(utils.MakeHTTPHandleFunc(testHandler.handleRegister))

	// Call the handler with the request and recorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Check the response body
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	expectedMessage := "User registered successfully"
	if response["message"] != expectedMessage {
		t.Errorf("handler returned unexpected body: got %v want %v", response["message"], expectedMessage)
	}

	if response["token"] == nil {
		t.Error("Expected token to be present in response")
	}

	if response["user_id"] == nil {
		t.Error("Expected user_id to be present in response")
	}
}

// Test case for bad request
func TestRegisterUserBadRequest(t *testing.T) {
	// Set up a test server with an invalid request
	invalidReqBody := map[string]string{
		"email":    "invalid_email",
		"username": "parikshith",
		"password": "short",
	}
	body, _ := json.Marshal(invalidReqBody)

	req, err := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testHandler := getTestHandler(t)
	handler := http.HandlerFunc(utils.MakeHTTPHandleFunc(testHandler.handleRegister))

	// Call the handler with the invalid request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func getTestHandler(t *testing.T) *Handler {
	db := utils.SetupTestPostgres(t)
	// defer utils.TeardownTestPostgres(t, db)
	testStore := NewStore(db)
	testHandler := NewHandler(testStore)
	return testHandler
}
