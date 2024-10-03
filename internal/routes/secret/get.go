package secret

import (
	"net/http"

	"pm4devs.strawhats/internal/rest"
)

func (s *Secret) getUserSecrets(w http.ResponseWriter, r *http.Request) {
	s.rest.WriteJSON(w, "secrets.getUserSecrets", http.StatusOK, rest.Envelope{"error": "testing error"})
}

