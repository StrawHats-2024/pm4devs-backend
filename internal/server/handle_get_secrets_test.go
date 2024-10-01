package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetUserSecrets(t *testing.T) {
	s := newTestServer()
	server := httptest.NewServer(http.HandlerFunc(s.handleGetUserSecrets))
	defer server.Close()

	// NOTE: API will return 200 empty items response in case a request doesn't satisfy a listRule, 400 for unsatisfied createRule and 404 for unsatisfied viewRule, updateRule and deleteRule.

	t.Run("Return 200 with empty items when invalid token", func(t *testing.T) {
		resp := makeNewReq(t, server.URL, http.MethodGet, nil, "")
		want := http.StatusOK
		got := resp.StatusCode
		assertStatusCode(t, want, got)
		resObj := getBodyObj[SecretsResponse](t, resp.Body)
		if len(resObj.Items) > 0 {
			t.Errorf("Expected len 0 but got %d", len(resObj.Items))
		}
	})
	t.Run("Return 200 when token valid", func(t *testing.T) {
		token := getAuthToken(t)
		resp := makeNewReq(t, server.URL, http.MethodGet, nil, token)
		want := http.StatusOK
		got := resp.StatusCode
		assertStatusCode(t, want, got)
		resObj := getBodyObj[SecretsResponse](t, resp.Body)
		if len(resObj.Items) == 0 {
			t.Errorf("Expected to get some items but got empty array")
		}

	})
}

func getBodyObj[T any](t *testing.T, body io.Reader) T {

	var obj T
	err := json.NewDecoder(body).Decode(&obj)
	if err != nil {
		t.Fatal(err)
	}
	return obj
}
