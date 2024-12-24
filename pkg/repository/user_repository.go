package repository

import (
    "context"
    "database/sql"
    "car-backend/pkg/models"
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
    _, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Name, user.PhotoURL)
    return err
}

func (r *UserRepository) UpdateProfile(ctx context.Context, userID string, update *models.UpdateUserProfile) error {
    query := `
        UPDATE users 
        SET 
            name = COALESCE($1, name),
            music_pref = COALESCE($2, music_pref),
            smoking_pref = COALESCE($3, smoking_pref),
            pet_friendly = COALESCE($4, pet_friendly),
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $5
    `
    _, err := r.db.ExecContext(ctx, query, 
        update.Name, 
        update.MusicPref, 
        update.SmokingPref, 
        update.PetFriendly, 
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
        &user.PhotoURL,
        &user.MusicPref,
        &user.SmokingPref,
        &user.PetFriendly,
        &user.Rating,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &user, nil
} 