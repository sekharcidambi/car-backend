package models

import (
    "time"
)

type User struct {
    ID            string    `json:"id" db:"id"`                     // Clerk User ID
    Email         string    `json:"email" db:"email"`
    Name          string    `json:"name" db:"name"`
    PhotoURL      string    `json:"photo_url" db:"photo_url"`
    MusicPref     string    `json:"music_pref" db:"music_pref"`
    SmokingPref   bool      `json:"smoking_pref" db:"smoking_pref"`
    PetFriendly   bool      `json:"pet_friendly" db:"pet_friendly"`
    Rating        float64   `json:"rating" db:"rating"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type UpdateUserProfile struct {
    Name        *string `json:"name,omitempty"`
    MusicPref   *string `json:"music_pref,omitempty"`
    SmokingPref *bool   `json:"smoking_pref,omitempty"`
    PetFriendly *bool   `json:"pet_friendly,omitempty"`
} 