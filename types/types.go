package types

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserId int64  `json:"user_id"`
}

type RegisterUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Username string `json:"username" validate:"required,min=3,max=30"`
}

type RegisterUserResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	UserId  int64  `json:"user_id"`
}

type ShareSecretToUserPayload struct {
	UserEmail   string         `json:"user_email" validate:"required,email"`
	Permissions PermissionType `json:"permissions" validate:"required"`
}

type ShareSecretToGroupPayload struct {
	GroupId     int            `json:"group_id" validate:"required"`
	Permissions PermissionType `json:"permissions" validate:"required"`
}

type GetSharedSecretRes struct {
	SecretID    string         `json:"secret_id"`
	SharedBy    string         `json:"shared_by"`
	Permissions PermissionType `json:"permissions"`
	Description string         `json:"description"`
}
type ShareSecretResponse struct {
	Message string `json:"message"`
}

type UpdateSecretPayload struct {
	EncryptedData string `json:"encrypted_data" db:"encrypted_data" validate:"required"`
	Description   string `json:"description" db:"description"`
}

type UserInGroup struct {
	UserID   int    `json:"user_id" db:"user_id"`
	Username string `json:"username" db:"username"`
}

type AddUserToGroupPayload struct {
	UserEmail string `json:"user_email" validate:"required,email"`
	Role      Role   `json:"role" validate:"required"`
}
type RevokeSecretAccessPayload struct {
	UserEmail string `json:"user_email" validate:"required,email"`
}
type GetUserGroupRes struct {
	GroupID   int    `json:"group_id"`
	GroupName string `json:"group_name"`
	Role      string `json:"role"`
}

type GroupWithUsers struct {
	Group Group           `json:"group"`
	Users []GroupUserItem `json:"users"`
}

type GroupUserItem struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
