package repository

import (
	"car-backend/pkg/models"
	"context"
	"database/sql"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (id, email, name, photo_url)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id) DO NOTHING
    `
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Name)
	return err
}

func (r *UserRepository) UpdateProfile(ctx context.Context, userID string, update *models.UpdateUserProfile) error {
	query := `
        UPDATE users 
        SET 
            display_name = COALESCE($1, display_name),
            city = COALESCE($2, city),
            state = COALESCE($3, state),
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $4
    `
	_, err := r.db.ExecContext(ctx, query,
		update.DisplayName,
		update.City,
		update.State,
		userID,
	)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.DisplayName,
		&user.City,
		&user.State,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUserIfNotExists(ctx context.Context, user *models.User) error {
	// Check if user exists
	var exists bool
	err := r.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)",
		user.Email,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("failed to check user existence: %v", err)
	}

	if exists {
		// User already exists, no need to create
		return nil
	}

	// User doesn't exist, create new user
	query := `
        INSERT INTO users (
            id, email, name, display_name, city, state
        ) VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING created_at, updated_at`

	err = r.db.QueryRowContext(ctx, query,
		user.ID,
		user.Email,
		user.Name,
		user.DisplayName,
		user.City,
		user.State,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}
