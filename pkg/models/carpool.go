package models

import (
	"time"

	"github.com/google/uuid"
)

type Carpool struct {
	ID              uuid.UUID `json:"id" db:"id"`
	CreatorID       string    `json:"creator_id" db:"creator_id"`
	Name            string    `json:"name" db:"name"`
	ScheduleDate    time.Time `json:"schedule_date" db:"schedule_date"`
	ScheduleTime    time.Time `json:"schedule_time" db:"schedule_time"`
	RecurringOption string    `json:"recurring_option" db:"recurring_option"`
	AvailableSeats  int       `json:"available_seats" db:"available_seats"`
	Stops           []Stop    `json:"stops"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type Stop struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CarpoolID uuid.UUID `json:"carpool_id" db:"carpool_id"`
	Address   string    `json:"address" db:"address"`
	StopOrder int       `json:"stop_order" db:"stop_order"`
	StopType  string    `json:"stop_type" db:"stop_type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateCarPoolRequest struct {
	Name               string `json:"name"`
	ScheduleDate       string `json:"schedule_date"`
	ScheduleTime       string `json:"schedule_time"`
	RecurringOption    string `json:"recurring_option"`
	AvailableSeats     int    `json:"available_seats"`
	StartAddress       string `json:"start_address"`
	DestinationAddress string `json:"destination_address"`
}

type StopRequest struct {
	Address   string `json:"address"`
	StopType  string `json:"stop_type"`
	StopOrder int    `json:"stop_order"`
}

type SearchFilters struct {
	StartDate       *time.Time `json:"start_date,omitempty"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	MinSeats        *int       `json:"min_seats,omitempty"`
	MaxDistance     *float64   `json:"max_distance,omitempty"`
	MusicPreference *string    `json:"music_preference,omitempty"`
	SmokingAllowed  *bool      `json:"smoking_allowed,omitempty"`
	PetsAllowed     *bool      `json:"pets_allowed,omitempty"`
}
