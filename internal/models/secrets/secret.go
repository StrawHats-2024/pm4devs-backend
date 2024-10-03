package secrets

import (
	"time"
)

// SecretRecord represents the secrets table in the database.
type SecretRecord struct {
	ID            int64     `db:"id"`             // Bigserial primary key
	Name          string    `db:"name"`           // Name of the secret
	EncryptedData []byte    `db:"encrypted_data"` // Encrypted credentials (bytea)
	OwnerID       int64     `db:"owner_id"`       // Foreign key referencing users(id)
	CreatedAt     time.Time `db:"created_at"`     // Timestamp with time zone
}
