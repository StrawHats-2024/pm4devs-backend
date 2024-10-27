package group

import (
	"time"
)

type GroupRecord struct {
	ID        int64     `db:"id" json:"id"`                   // Primary key
	Name      string    `db:"name" json:"name"`               // Group name (unique, not null)
	CreatorID int64     `db:"creator_id" json:"creator_id"`   // Foreign key referencing users (creator)
	CreatedAt time.Time `db:"created_at" json:"created_at"`   // Timestamp when the group was created
}

type GroupMemberRecord struct {
	GroupID int64 `db:"group_id" json:"group_id"`  // Foreign key referencing groups
	UserID  int64 `db:"user_id" json:"user_id"`    // Foreign key referencing users
}
