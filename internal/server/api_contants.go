package server

import (
	"fmt"
)

type APIEndpoints struct {
	AuthWithPassword string
	AuthTokenRefresh string
	UserVerification string
	UserCollection   string
	GroupCollection  string
}

// Function to create the struct with BASE_URL
func NewAPIEndpoints(baseURL string) *APIEndpoints {
	return &APIEndpoints{
		AuthWithPassword: fmt.Sprintf("%s/%s/", baseURL, "collections/users/auth-with-password"),
		AuthTokenRefresh: fmt.Sprintf("%s/%s/", baseURL, "collections/users/auth-refresh"),
		UserVerification: fmt.Sprintf("%s/%s/", baseURL, "collections/users/request-verification"),
		UserCollection:   fmt.Sprintf("%s/%s/", baseURL, "collections/users/records"),
		GroupCollection:  fmt.Sprintf("%s/%s/", baseURL, "collections/group/records"),
	}
}
