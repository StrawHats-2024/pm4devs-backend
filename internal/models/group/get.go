package group

import (
	"context"
	"time"

	"pm4devs.strawhats/internal/xerrors"
)

func (g *Group) GetGroupsByUserID(userID int64) ([]GroupRecord, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Query to get groups the user is part of
	queryGroups := `
		SELECT gr.id, gr.name, gr.creator_id, gr.created_at
		FROM groups gr
		JOIN group_members gm ON gm.group_id = gr.id
		WHERE gm.user_id = $1;
	`

	rows, err := g.DB.QueryContext(ctx, queryGroups, userID)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "group.GetGroupsByUserID")
	}
	defer rows.Close()

	var groups []GroupRecord

	for rows.Next() {
		var group GroupRecord
		if err := rows.Scan(&group.ID, &group.Name, &group.CreatorID, &group.CreatedAt); err != nil {
			return nil, xerrors.DatabaseError(err, "group.GetGroupsByUserID")
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, xerrors.DatabaseError(err, "group.GetGroupsByUserID")
	}

	return groups, nil
}
