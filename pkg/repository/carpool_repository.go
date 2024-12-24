package repository

import (
	"car-backend/pkg/models"
	"context"
	"database/sql"
)

type CarPoolRepository struct {
	db *sql.DB
}

func NewCarPoolRepository(db *sql.DB) *CarPoolRepository {
	return &CarPoolRepository{db: db}
}

// Implement the CreateCarPool method
func (r *CarPoolRepository) CreateCarPool(ctx context.Context, carpool *models.CarPool) error {
	query := `
		INSERT INTO carpools (
			creator_id, schedule_date, schedule_time, recurring_option,
			starting_point_lat, starting_point_lng, destination_lat, destination_lng,
			available_seats, music_preference, smoking_allowed, pets_allowed
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`

	return r.db.QueryRowContext(
		ctx, query,
		carpool.CreatorID, carpool.ScheduleDate, carpool.ScheduleTime,
		carpool.RecurringOption, carpool.StartPointLat, carpool.StartPointLng,
		carpool.DestinationLat, carpool.DestinationLng, carpool.AvailableSeats,
		carpool.MusicPreference, carpool.SmokingAllowed, carpool.PetsAllowed,
	).Scan(&carpool.ID)
}

// Add methods like:
// CreateCarPool
// GetCarPool
// UpdateCarPool
// DeleteCarPool
// SearchCarPools
