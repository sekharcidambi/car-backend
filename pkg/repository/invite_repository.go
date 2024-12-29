package repository

import (
    "car-backend/pkg/models"
	"context"
	"database/sql"
	"fmt"

	
)

type InviteRepository struct {
    db *sql.DB
}

func NewInviteRepository(db *sql.DB) *InviteRepository {
    return &InviteRepository{db: db}
}

func (r *InviteRepository) CreateInvite(ctx context.Context, invite *models.Invite) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
            return fmt.Errorf("failed to begin transaction: %v", err)
    }
    defer tx.Rollback()

    query := `
            INSERT INTO invites (
                    from_user, to_user, carpool_id, message, status
            ) VALUES ($1, $2, $3, $4, $5)
            RETURNING id, created_at, updated_at
    `

    err = tx.QueryRowContext(
            ctx, query,
            invite.FromUser, invite.ToUser, invite.CarpoolID, invite.Message, invite.Status,
    ).Scan(&invite.ID, &invite.CreatedAt, &invite.UpdatedAt)

    if err != nil {
            return fmt.Errorf("failed to create invite: %v", err)
    }

    if err := tx.Commit(); err != nil {
            return fmt.Errorf("failed to commit transaction: %v", err)
    }

    return nil
}


