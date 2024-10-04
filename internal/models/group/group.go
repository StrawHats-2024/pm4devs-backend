package group

import (
	"time"
)

type GroupRecord struct {
	ID        int64     `db:"id"`         // Primary key
	Name      string    `db:"name"`       // Group name (unique, not null)
	CreatorID int64     `db:"creator_id"` // Foreign key referencing users (creator)
	CreatedAt time.Time `db:"created_at"` // Timestamp when the group was created
}

type GroupMemberRecord struct {
	GroupID int64 `db:"group_id"` // Foreign key referencing groups
	UserID  int64 `db:"user_id"`  // Foreign key referencing users
}
