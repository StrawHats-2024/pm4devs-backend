package auth

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

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type LoginReq struct {
	Email    string  `json:"email"`
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
	Message string `json:"message"`
	UserId  int64  `json:"user_id"`
}

// TODO: Make env
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

func HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	var req UserRegReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("err: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: get user from db and validate password
	user, err := NewUser("test", "testing")
	if err != nil {
		fmt.Println("err: ", err)
		log.Fatal(err)
	}
	if !user.ValidPassword(user.PasswordHash) {
		fmt.Println("Password no match")
		w.WriteHeader(http.StatusUnauthorized)
		return
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
	err = WriteJSON(w, http.StatusOK, LoginRes{
		Token:  token,
		UserId: int64(user.UserID),
	})
}

func HandleUserReg(w http.ResponseWriter, r *http.Request) {
	var req UserRegReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("err: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := NewUser(string(req.Email), req.Password)
	if err != nil {
		fmt.Println("err: ", err)
		log.Fatal(err)
	}
	// if !user.ValidPassword(user.PasswordHash) {
	// 	fmt.Print("Password did not match")
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
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
		Message: "User registered successfully",
		UserId:  int64(user.UserID),
	})
}

func NewUser(email, password string) (*models.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// TODO: add new user to db

	return &models.User{
		Email:        email,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now().UTC(),
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
