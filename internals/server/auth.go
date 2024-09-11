package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pm4devs-backend/pkg/models"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
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
	token, err := createToken(user.Email)
	if err != nil {
		log.Fatal(err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
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
	user, err := NewUser(req.Email, req.Password)
	if err != nil {
		return err
	}
	userId, err := s.store.CreateUser(user)
	if err != nil {
		return err
	}

	token, err := createToken(user.Email)
	if err != nil {
		log.Fatal(err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
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

// TODO: Make secretKey env
var secretKey = []byte("secret-key")

func createToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewUser(email, password string) (*models.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &models.User{
		Email:        email,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now().UTC(),
		LastLogin:    time.Now().UTC(),
	}, nil
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
