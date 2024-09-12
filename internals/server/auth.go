package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req LoginReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("err: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	user, err := s.store.GetUserByEmail(req.Email)
	// TODO: Update error message to when email not found
	if err != nil {
		return err
	}

	if !user.ValidPassword(req.Password) {
		WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Wrong password or email"})
		return nil
	}
	token, err := createJWT(user)
	if err != nil {
		log.Fatal(err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 24),
	})
	err = s.store.UpdateLastLogin(user.UserID)
	if err != nil {
		return err
	}
	err = WriteJSON(w, http.StatusOK, LoginRes{
		Token:  token,
		UserId: int64(user.UserID),
	})
	return nil

}

func (s *APIServer) handleRegister(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req UserRegReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("err: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	user, err := NewUser(req.Email, req.Username, req.Password)
	if err != nil {
		return err
	}
	userId, err := s.store.CreateUser(user)
	if err != nil {
		fmt.Printf("user: %+v", user)
		fmt.Println("err: ", err)
		return err
	}

	token, err := createJWT(user)
	if err != nil {
		log.Fatal(err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 24),
	})
	err = WriteJSON(w, http.StatusOK, UserRegRes{
		Token:   token,
		Message: "User Registered successfully.",
		UserId:  int64(userId),
	})
	return nil

}

func (s *APIServer) handleTokenRefresh(w http.ResponseWriter, r *http.Request) error {
	// TODO: Impliment token refresh
	return nil
}

func (s *APIServer) handleLogout(w http.ResponseWriter, r *http.Request) error {
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
	err := WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Successfully logged out",
	})
	if err != nil {
		return err
	}

	return nil
}

// auth middleware
func withAuth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the token is missing, return an unauthorized status
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}
			// For any other error, return a bad request status
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		tokenString := cookie.Value
		claims, err := validateToken(tokenString)
		fmt.Println("claims: ", claims.UserId)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, err.Error())
			return
		}
		f(w, r)
	}

}
