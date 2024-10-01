package server

type CreateUserPayload struct {
	Name            string `json:"name,omitempty"`
	Email           string `json:"email,omitempty"`
	Password        string `json:"password,omitempty"`
	PasswordConfirm string `json:"passwordConfirm,omitempty"`
}

type LoginPayload struct {
	Identity string `json:"identity,omitempty"`
	Password string `json:"password,omitempty"`
}

// Also refresh token response
type LoginResponse struct {
	Token  string     `json:"token,omitempty"`
	Record UserRecord `json:"record"`
}

type UserRecord struct {
	ID              string `json:"id"`
	CollectionID    string `json:"collectionId"`
	CollectionName  string `json:"collectionName"`
	Username        string `json:"username"`
	Verified        bool   `json:"verified"`
	EmailVisibility bool   `json:"emailVisibility"`
	Email           string `json:"email"`
	Created         string `json:"created"`
	Updated         string `json:"updated"`
	Name            string `json:"name"`
	Avatar          string `json:"avatar"`
}

type SecretsResponse struct {
	Page       int            `json:"page"`
	PerPage    int            `json:"perPage"`
	TotalPages int            `json:"totalPages"`
	TotalItems int            `json:"totalItems"`
	Items      []SecretRecord `json:"items"`
}

type SecretRecord struct {
	ID             string `json:"id"`
	CollectionID   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	Created        string `json:"created"`
	Updated        string `json:"updated"`
	Name           string `json:"name"`
	EncryptedData  string `json:"encrypted_data"`
	Owner          string `json:"owner"`
}

type ErrorResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    ValidationData `json:"data"`
}

// ValidationData represents the structure of the data field in the error response
type ValidationData struct {
	Name ValidationError `json:"name"`
}

// ValidationError represents the structure of the validation error
type ValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
