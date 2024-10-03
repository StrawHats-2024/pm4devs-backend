package secret

import (
	"net/http"

	"pm4devs.strawhats/internal/rest"
)

const GetUserSecretsRoute = "/v1/user/secrets"
func (s *Secret) getUserSecrets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	s.rest.WriteJSON(w, "secrets.getUserSecrets", http.StatusOK, rest.Envelope{"error": "testing error"})
}

