package handlers

import (
	"encoding/json"
	"net/http"
	"car-backend/pkg/models"
	"car-backend/pkg/repository"
	
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
	// Verify webhook signature here
	var event struct {
		Type string `json:"type"`
		Data struct {
			ID    string `json:"id"`
			Email string `json:"email"`
			Name  string `json:"name"`
			ImageURL string `json:"image_url"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if event.Type == "user.created" {
		user := &models.User{
			ID:       event.Data.ID,
			Email:    event.Data.Email,
			Name:     event.Data.Name,
			PhotoURL: event.Data.ImageURL,
		}
		
		if err := h.userRepo.CreateUser(r.Context(), user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	
	w.WriteHeader(http.StatusOK)
}

// GetProfile returns the user's profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), claims.Subject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// UpdateProfile updates the user's profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var update models.UpdateUserProfile
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.userRepo.UpdateProfile(r.Context(), claims.Subject, &update); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
} 