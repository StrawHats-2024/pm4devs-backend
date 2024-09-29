package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleLogin(t *testing.T) {
	s := newTestServer()
	server := httptest.NewServer(http.HandlerFunc(s.handleLogin))
	defer server.Close()

	t.Run("Return 400 when body not correct", func(t *testing.T) {
		resp := makePostReq(t, server.URL, nil)
		want := http.StatusBadRequest
		got := resp.StatusCode
		assertStatusCode(t, want, got)
	})

	t.Run("Return 200 when request body is correct", func(t *testing.T) {
		payload := LoginPayload{Identity: "javonconnelly@sanford.org", Password: "parikshith"}
		resp := makePostReq(t, server.URL, getBodyJson(t, payload))

		got := resp.StatusCode
		want := http.StatusOK
		assertStatusCode(t, want, got)
		var resBody LoginResponse
		err := json.NewDecoder(resp.Body).Decode(&resBody)
		if err != nil {
			t.Fatal(err)
		}
		assertAny(t, payload.Identity, resBody.Record.Email)
		if len(resBody.Token) == 0 {
			t.Errorf("No auth token in response")
		}
	})
}

func assertAny(t *testing.T, want, got any) {
	if got != want {
		t.Errorf("Expected %v, got %v", want, got)
	}
}
