package secrets

import (
	"time"
)

// SecretRecord represents the secrets table in the database.
type SecretRecord struct {
	ID            int64     `db:"id" json:"id"`                       // Bigserial primary key
	Name          string    `db:"name" json:"name"`                   // Name of the secret
	EncryptedData []byte    `db:"encrypted_data" json:"encrypted_data"` // Encrypted credentials (bytea)
	IV            []byte    `db:"iv" json:"iv"`                       // Initialization Vector (bytea)
	OwnerID       int64     `db:"owner_id" json:"owner_id"`           // Foreign key referencing users(id)
	CreatedAt     time.Time `db:"created_at" json:"created_at"`       // Timestamp with time zone
}
