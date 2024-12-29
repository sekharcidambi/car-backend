package handlers

import (
	"car-backend/pkg/models"
	"car-backend/pkg/repository"
	"encoding/json"
	"log"
	"net/http"
	"github.com/google/uuid"
)

type InviteHandler struct {
	inviteRepo *repository.InviteRepository
}

func NewInviteHandler(repo *repository.InviteRepository) *InviteHandler {
	return &InviteHandler{
		inviteRepo: repo,
	}
}

func (h *InviteHandler) CreateInvite(w http.ResponseWriter, r *http.Request) {
    var req models.CreateInviteRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            log.Printf("{\"severity\":\"ERROR\",\"message\":\"Failed to decode request: %v\"}", err)
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
    }

    // Placeholder: Replace with actual user ID retrieval logic
    userID, _ := uuid.Parse("user-123") 

    // Create Invite object
    invite := &models.Invite{
            FromUser:  userID, 
            ToUser:    req.ToUser,
            CarpoolID: req.CarpoolID,
            Message:   req.Message,
            Status:    0, // Initial status: pending
    }

    if err := h.inviteRepo.CreateInvite(r.Context(), invite); err != nil {
            log.Printf("{\"severity\":\"ERROR\",\"message\":\"Failed to create invite: %v\"}", err)
            http.Error(w, "Failed to create invite", http.StatusInternalServerError)
            return
    }

    w.WriteHeader(http.StatusCreated) 
}

