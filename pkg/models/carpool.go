package models

import (
	"time"
)

type CarPool struct {
	ID              string    `json:"id" db:"id"`
	CreatorID       string    `json:"creator_id" db:"creator_id"`
	ScheduleDate    time.Time `json:"schedule_date" db:"schedule_date"`
	ScheduleTime    time.Time `json:"schedule_time" db:"schedule_time"`
	RecurringOption string    `json:"recurring_option" db:"recurring_option"`
	StartPointLat   float64   `json:"start_point_lat" db:"starting_point_lat"`
	StartPointLng   float64   `json:"start_point_lng" db:"starting_point_lng"`
	DestinationLat  float64   `json:"destination_lat" db:"destination_lat"`
	DestinationLng  float64   `json:"destination_lng" db:"destination_lng"`
	AvailableSeats  int       `json:"available_seats" db:"available_seats"`
	MusicPreference string    `json:"music_preference" db:"music_preference"`
	SmokingAllowed  bool      `json:"smoking_allowed" db:"smoking_allowed"`
	PetsAllowed     bool      `json:"pets_allowed" db:"pets_allowed"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type CreateCarPoolRequest struct {
	ScheduleDate    string  `json:"schedule_date"`
	ScheduleTime    string  `json:"schedule_time"`
	RecurringOption string  `json:"recurring_option"`
	StartPointLat   float64 `json:"start_point_lat"`
	StartPointLng   float64 `json:"start_point_lng"`
	DestinationLat  float64 `json:"destination_lat"`
	DestinationLng  float64 `json:"destination_lng"`
	AvailableSeats  int     `json:"available_seats"`
	MusicPreference string  `json:"music_preference"`
	SmokingAllowed  bool    `json:"smoking_allowed"`
	PetsAllowed     bool    `json:"pets_allowed"`
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
