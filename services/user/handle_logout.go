package user

import (
	"fmt"
	"net/http"
	"pm4devs-backend/utils"
	"time"
)

func (s *Handler) handleLogout(w http.ResponseWriter, r *http.Request) error {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	// Clear the token by setting the cookie expiration date to a past date
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",                             // Empty token value
		Path:    "/",                            // Ensure it applies to the whole site
		Expires: time.Now().Add(-1 * time.Hour), // Set to a past time to invalidate the cookie
	})

	// Respond with a successful logout message
	err := utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Successfully logged out",
	})
	if err != nil {
		return err
	}

	return nil
}

