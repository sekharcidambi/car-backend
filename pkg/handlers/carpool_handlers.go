package handlers

import (
	"car-backend/pkg/models"
	"car-backend/pkg/repository"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CarPoolHandler struct {
	carpoolRepo *repository.CarPoolRepository
}


func NewCarPoolHandler(repo *repository.CarPoolRepository) *CarPoolHandler {
	return &CarPoolHandler{
		carpoolRepo: repo,
	}
}

func (h *CarPoolHandler) CreateCarPool(w http.ResponseWriter, r *http.Request) {
    var req models.CreateCarPoolRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("{\"severity\":\"ERROR\",\"message\":\"Failed to decode request: %v\"}", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Create carpool object
    carpool := &models.Carpool{
        CreatorID:      "temp-creator", // TODO: Get from auth context
        CarpoolName:    req.CarpoolName,  // Use CarpoolName from request
        Status:         false,           // Default status (can be modified later)
        RecurringOption: req.RecurringOption,
        AvailableSeats:  req.AvailableSeats,
        DestinationAddress: req.DestinationAddress,
        Seats:           req.Seats,        // Use Seats from request
    }



    if err := h.carpoolRepo.CreateCarPool(r.Context(), carpool); err != nil {
        log.Printf("{\"severity\":\"ERROR\",\"message\":\"Failed to create carpool: %v\"}", err)
        http.Error(w, "Failed to create carpool", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(carpool)
}

func (h *CarPoolHandler) GetCarPool(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    params := mux.Vars(r)
    carpoolID, err := uuid.Parse(params["id"])
    if err != nil {
        http.Error(w, "Invalid carpool ID", http.StatusBadRequest)
        return
    }

    carpool, err := h.carpoolRepo.GetCarPool(context.Background(), carpoolID)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Carpool not found", http.StatusNotFound)
            return
        }
        http.Error(w, fmt.Sprintf("Failed to get carpool: %v", err), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(carpool)
}

func (h *CarPoolHandler) UpdateCarPool(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
}

func (h *CarPoolHandler) DeleteCarPool(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    params := mux.Vars(r)
    carpoolIDStr := params["id"] 

    carpoolID, err := uuid.Parse(carpoolIDStr)
    if err != nil {
        http.Error(w, "Invalid carpool ID", http.StatusBadRequest)
        return
    }

    err = h.carpoolRepo.DeleteCarPool(context.Background(), carpoolID)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Carpool not found", http.StatusNotFound)
            return
        }
        http.Error(w, fmt.Sprintf("Failed to delete carpool: %v", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent) 
}


func (h *CarPoolHandler) SearchCarPools(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
}