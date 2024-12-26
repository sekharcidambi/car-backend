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

func (r *CarPoolRepository) GetCarPool(ctx context.Context, carpoolID string) (*models.CarPool, error) {
	// Create a new CarPool struct to store the retrieved data
	carpool := &models.CarPool{}

	// Define the select query with placeholder for carpool ID
	query := `
        SELECT id, creator_id, schedule_date, schedule_time, recurring_option,
               starting_point_lat, starting_point_lng, destination_lat, destination_lng,
               available_seats, music_preference, smoking_allowed, pets_allowed
        FROM carpools
        WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, carpoolID)

	err := row.Scan(
		&carpool.ID,
		&carpool.CreatorID,
		&carpool.ScheduleDate,
		&carpool.ScheduleTime,
		&carpool.RecurringOption,
		&carpool.StartPointLat,
		&carpool.StartPointLng,
		&carpool.DestinationLat,
		&carpool.DestinationLng,
		&carpool.AvailableSeats,
		&carpool.MusicPreference,
		&carpool.SmokingAllowed,
		&carpool.PetsAllowed,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return carpool, nil
}

func (r *CarPoolRepository) DeleteCarPool(ctx context.Context, carpoolID string) error {
	// Define the delete query with placeholder for carpool ID
	query := `DELETE FROM carpools WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, carpoolID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // No rows affected, likely carpool not found
	}

	return nil
}

// Add methods like:
// UpdateCarPool
// SearchCarPools
