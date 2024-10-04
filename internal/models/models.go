package models

import (
	"database/sql"

	"pm4devs.strawhats/internal/models/group"
	"pm4devs.strawhats/internal/models/permissions"
	"pm4devs.strawhats/internal/models/secrets"
	"pm4devs.strawhats/internal/models/tokens"
	"pm4devs.strawhats/internal/models/users"
)

// Encapsulates all the models
type Models struct {
	Permissions permissions.PermissionsRepository
	Tokens      tokens.TokensRepository
	Users       users.UsersRepository
	Secrets     secrets.SecretsRepository
	Group       group.GroupRepository
}

func New(db *sql.DB) *Models {
	return &Models{
		Permissions: permissions.Repository(db),
		Tokens:      tokens.Repository(db),
		Users:       users.Repository(db),
		Secrets:     secrets.Repository(db),
		Group:       group.Repository(db),
	}
}
