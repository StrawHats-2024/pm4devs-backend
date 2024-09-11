package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pm4devs-backend/pkg/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRes struct {
	Token  string `json:"token"`
	UserId int64  `json:"user_id"`
}

type UserRegReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}
type UserRegRes struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	UserId  int64  `json:"user_id"`
}

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
	user, err := NewUser(req.Email, req.Password, req.Username)
	if err != nil {
		return err
	}
	userId, err := s.store.CreateUser(user)
	if err != nil {
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

func NewUser(email, username, password string) (*models.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &models.User{
		Email:        email,
		PasswordHash: string(passwordHash),
		Username:     username,
		CreatedAt:    time.Now().UTC(),
		LastLogin:    time.Now().UTC(),
	}, nil
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
