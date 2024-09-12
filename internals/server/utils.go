package server

import (
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


func getUserIdfromCookie(r *http.Request) (int, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return -1, err
	}
	claims, err := validateToken(cookie.Value)
	if err != nil {
		return -1, err
	}
	return claims.UserId, nil
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
