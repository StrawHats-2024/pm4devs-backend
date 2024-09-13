package types

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserId int64  `json:"userId"`
}

type RegisterUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Username string `json:"username" validate:"required,min=3,max=30"`
}

type RegisterUserResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	UserId  int64  `json:"userId"`
}

type UpdateSecretPayload struct {
	EncryptedData string `json:"encrypted_data" db:"encrypted_data"`
	Description   string `json:"description" db:"description"`
}

type UserInGroup struct {
	UserID   int    `json:"user_id" db:"user_id"`
	Username string `json:"username" db:"username"`
}

type AddUserToGroupPayload struct {
	UserEmail string `json:"user_email"`
	Role      Role   `json:"role"`
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
