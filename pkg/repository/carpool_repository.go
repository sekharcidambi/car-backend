package repository

import (
	"car-backend/pkg/models"
	"context"
	"database/sql"
	"fmt"
	"log"
)

type CarPoolRepository struct {
	db *sql.DB
}

func NewCarPoolRepository(db *sql.DB) *CarPoolRepository {
	return &CarPoolRepository{db: db}
}

// Implement the CreateCarPool method
func (r *CarPoolRepository) CreateCarPool(ctx context.Context, carpool *models.Carpool) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Insert main carpool record
	query := `
		INSERT INTO carpools (
			creator_id, name, schedule_date, schedule_time, recurring_option,
			available_seats
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	err = tx.QueryRowContext(
		ctx, query,
		carpool.CreatorID, carpool.Name, carpool.ScheduleDate, carpool.ScheduleTime,
		carpool.RecurringOption, carpool.AvailableSeats,
	).Scan(&carpool.ID, &carpool.CreatedAt, &carpool.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert carpool: %v", err)
	}

	// Insert carpool stops
	stopsQuery := `
		INSERT INTO carpool_stops (
			carpool_id, address, stop_order, stop_type
		) VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	for i := range carpool.Stops {
		stop := &carpool.Stops[i]
		stop.CarpoolID = carpool.ID

		err = tx.QueryRowContext(
			ctx, stopsQuery,
			stop.CarpoolID, stop.Address, stop.StopOrder, stop.StopType,
		).Scan(&stop.ID, &stop.CreatedAt, &stop.UpdatedAt)

		if err != nil {
			return fmt.Errorf("failed to insert stop: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *CarPoolRepository) GetCarPool(ctx context.Context, carpoolID string) (*models.Carpool, error) {
	// Create a new Carpool struct to store the retrieved data
	carpool := &models.Carpool{}

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Get carpool details
	carpoolQuery := `
		SELECT id, creator_id, name, schedule_date, schedule_time, recurring_option,
			   available_seats, created_at, updated_at
		FROM carpools
		WHERE id = $1`

	err = tx.QueryRowContext(ctx, carpoolQuery, carpoolID).Scan(
		&carpool.ID,
		&carpool.CreatorID,
		&carpool.Name,
		&carpool.ScheduleDate,
		&carpool.ScheduleTime,
		&carpool.RecurringOption,
		&carpool.AvailableSeats,
		&carpool.CreatedAt,
		&carpool.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get carpool: %v", err)
	}

	// Get stops for the carpool
	stopsQuery := `
		SELECT id, carpool_id, address, stop_order, stop_type, created_at, updated_at
		FROM carpool_stops
		WHERE carpool_id = $1
		ORDER BY stop_order ASC`

	rows, err := tx.QueryContext(ctx, stopsQuery, carpoolID)
	if err != nil {
		return nil, fmt.Errorf("failed to get carpool stops: %v", err)
	}
	defer rows.Close()

	var stops []models.Stop
	for rows.Next() {
		var stop models.Stop
		err := rows.Scan(
			&stop.ID,
			&stop.CarpoolID,
			&stop.Address,
			&stop.StopOrder,
			&stop.StopType,
			&stop.CreatedAt,
			&stop.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stop: %v", err)
		}
		stops = append(stops, stop)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating stops: %v", err)
	}

	carpool.Stops = stops

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return carpool, nil
}

func (r *CarPoolRepository) DeleteCarPool(ctx context.Context, carpoolID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// First verify the carpool exists and get its stops count
	var stopsCount int
	err = tx.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM carpool_stops WHERE carpool_id = $1
	`).Scan(&stopsCount)
	if err != nil {
		return fmt.Errorf("failed to verify carpool stops: %v", err)
	}

	// Delete the carpool (stops will be deleted automatically due to CASCADE)
	result, err := tx.ExecContext(ctx, `DELETE FROM carpools WHERE id = $1`, carpoolID)
	if err != nil {
		return fmt.Errorf("failed to delete carpool: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // No carpool found
	}

	// Log the operation
	log.Printf("{\"severity\":\"INFO\",\"message\":\"Deleted carpool %s with %d stops\"}", carpoolID, stopsCount)

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// Add methods like:
// UpdateCarPool
// SearchCarPools
