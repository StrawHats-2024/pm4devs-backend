package auth

import (
	"fmt"
	"net/http"
	"pm4devs-backend/types"
	"pm4devs-backend/utils"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

// TODO: Make secretKey env
var secretKey = []byte("secret-key")

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Check if the token is valid
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func CreateJWT(user *types.User) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24)
	claims := &Claims{
		UserId: user.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func RefreshToken(tokenString string) (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Handle token parsing errors or invalid tokens
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", fmt.Errorf("invalid token signature")
		}
		return "", err
	}

	// Token is expired or invalid
	if !token.Valid && !claims.VerifyExpiresAt(time.Now(), false) {
		return "", fmt.Errorf("token is expired and cannot be refreshed")
	}

	// At this point, the token is valid or within a reasonable refresh window.
	// Create a new token with an extended expiration time
	newExpirationTime := time.Now().Add(15 * time.Minute) // New token with extended expiration

	newClaims := &Claims{
		UserId: claims.UserId, // Keep the same user ID
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(newExpirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create and sign the new token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenString, err = newToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}


// auth middleware
func WithAuth(f http.HandlerFunc) http.HandlerFunc {
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
		claims, err := ValidateToken(tokenString)
		fmt.Println("claims: ", claims.UserId)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, err.Error())
			return
		}
		f(w, r)
	}

}


// Gets userId by reading cookie from the request
func GetUserIdfromRequest(r *http.Request) (int, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return -1, err
	}
	claims, err := ValidateToken(cookie.Value)
	if err != nil {
		return -1, err
	}
	return claims.UserId, nil
}

