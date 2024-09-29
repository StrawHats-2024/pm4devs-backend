package server

import (
	"bytes"
	"encoding/json"
	gofakeit "github.com/brianvoe/gofakeit/v7"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRegister(t *testing.T) {
	s := &Server{
		APIEndpoints: NewAPIEndpoints("http://127.0.0.1:8090/api"),
	}
	server := httptest.NewServer(http.HandlerFunc(s.handleRegister))
	defer server.Close()
	t.Run("Return 400 when body not correct", func(t *testing.T) {
		resp, err := http.Post(server.URL, "application/json", nil)
		if err != nil {
			t.Fatal(err)
		}
		assertStatusCode(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Get 200 on success", func(t *testing.T) {
		reqBody := CreateUserPayload{
			Email:           gofakeit.Email(),
			Password:        "parikshith",
			PasswordConfirm: "parikshith",
			Name:            gofakeit.Name(),
		}
		resp, err := http.Post(server.URL, "application/json", getBodyJson(t, reqBody))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close() // Make sure to close response body

		got := resp.StatusCode
		want := http.StatusOK
		assertStatusCode(t, want, got)
	})
}

func assertStatusCode(t *testing.T, want, got int) {
	if got != want {
		t.Errorf("Expected %d, got %d", want, got)
	}
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
