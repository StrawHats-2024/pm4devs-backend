package server

import (
	"fmt"
	"pm4devs-backend/pkg/models"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

// TODO: Make secretKey env
var secretKey = []byte("secret-key")

func validateToken(tokenString string) (*Claims, error) {
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

func createJWT(user *models.User) (string, error) {
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

func refreshToken(tokenString string) (string, error) {
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
