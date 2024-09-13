package utils

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
	"pm4devs-backend/types"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var Validate = validator.New()

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type ApiError struct {
	Error string `json:"error"`
}

func NewUser(email, username, password string) (*types.User, error) {
	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	return &types.User{
		Email:        email,
		PasswordHash: string(passwordHash),
		Username:     username,
		CreatedAt:    time.Now().UTC(),
		LastLogin:    time.Now().UTC(),
	}, nil
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func MakeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
