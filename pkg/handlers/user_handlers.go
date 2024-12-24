package handlers

import (
	"car-backend/pkg/models"
	"car-backend/pkg/repository"
	"encoding/json"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// HandleWebhook processes Clerk webhooks for user events
func (h *UserHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Implement webhook handling for Clerk events
	// This is where you'll handle user creation/updates from Clerk
	// Example: user.created, user.updated, etc.
}

// GetProfile returns the user's profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := clerk.SessionClaimsFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user from our database
	user, err := h.userRepo.GetByID(ctx, claims.Subject)
	if err != nil {
		http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateProfile updates the user's profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := clerk.SessionClaimsFromContext(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var update models.UpdateUserProfile
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.userRepo.UpdateProfile(ctx, claims.Subject, &update); err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
