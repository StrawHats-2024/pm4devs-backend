package group

import (
	"context"
	"database/sql"
	"time"

	"pm4devs.strawhats/internal/models/core"
	"pm4devs.strawhats/internal/xerrors"
)

type GroupRepository interface {
	GetByGroupID(id int64) (*GroupRecord, *xerrors.AppError)
	UpdateGroupName(newName string, groupID int64) (*GroupRecord, *xerrors.AppError)
	DeleteByGroupID(groupID int64) *xerrors.AppError
	NewRecord(name string, ownerID int64) (*GroupRecord, *xerrors.AppError)
}

type Group struct {
	DB core.Queryable
}

func Repository(db core.Queryable) GroupRepository {
	return &Group{DB: db}
}

func (g *Group) NewRecord(name string, ownerID int64) (*GroupRecord, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO groups (name, creator_id, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, name, creator_id, created_at;
	`

	var newGroup GroupRecord
	err := g.DB.QueryRowContext(ctx, query, name, ownerID).
		Scan(&newGroup.ID, &newGroup.Name, &newGroup.CreatorID, &newGroup.CreatedAt)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "group.NewRecord")
	}

	return &newGroup, nil
}

func (g *Group) GetByGroupID(id int64) (*GroupRecord, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, name, creator_id, created_at
		FROM groups
		WHERE id = $1;
	`

	var group GroupRecord
	err := g.DB.QueryRowContext(ctx, query, id).Scan(&group.ID, &group.Name, &group.CreatorID, &group.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, xerrors.DatabaseError(err, "group.GetByGroupID")
		}
		return nil, xerrors.DatabaseError(err, "group.GetByGroupID")
	}

	return &group, nil
}

func (g *Group) UpdateGroupName(newName string, groupID int64) (*GroupRecord, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE groups
		SET name = $1
		WHERE id = $2
		RETURNING id, name, creator_id, created_at;
	`

	var updatedGroup GroupRecord
	err := g.DB.QueryRowContext(ctx, query, newName, groupID).
		Scan(&updatedGroup.ID, &updatedGroup.Name, &updatedGroup.CreatorID, &updatedGroup.CreatedAt)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "group.UpdateByGroupID")
	}

	return &updatedGroup, nil
}

func (g *Group) DeleteByGroupID(groupID int64) *xerrors.AppError {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM groups WHERE id = $1;`

	_, err := g.DB.ExecContext(ctx, query, groupID)
	if err != nil {
		return xerrors.DatabaseError(err, "group.DeleteByGroupID")
	}

	return nil
}
