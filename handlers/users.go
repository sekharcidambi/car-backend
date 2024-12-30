package handlers

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	ClerkID string `json:"clerk_id"`
	Email   string `json:"email"`
}

func CreateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user := models.User{
			ID:      uuid.New().String(),
			ClerkID: req.ClerkID,
			Email:   req.Email,
		}

		if result := db.Create(&user); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
} 