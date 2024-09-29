package server

import (
	"encoding/json"
	"fmt"
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

func TestHandleRefreshToken(t *testing.T) {
	s := newTestServer()
	server := httptest.NewServer(http.HandlerFunc(s.handleRefreshToken))
	defer server.Close()

	t.Run("Return 200 when valid bearer", func(t *testing.T) {
		token := getAuthToken(t)
		resp := makeNewReq(t, server.URL, http.MethodPost, nil, token)
		assertStatusCode(t, http.StatusOK, resp.StatusCode)
		var resBody LoginResponse
		err := json.NewDecoder(resp.Body).Decode(&resBody)
		if err != nil {
			t.Error(err)
		}
		if resBody.Token == "" {
			t.Errorf("Token not in refresh token response response")
		}
	})
	t.Run("Return 401 when no valid bearer", func(t *testing.T) {
		resp := makeNewReq(t, server.URL, http.MethodPost, nil, "")
		assertStatusCode(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
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
	fmt.Println("Token: ", resBody.Token)
	return resBody.Token
}
