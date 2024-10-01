package server

import (
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

