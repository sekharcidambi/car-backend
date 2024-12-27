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
	"time"

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

	// Debug log the received time
	log.Printf("{\"severity\":\"DEBUG\",\"message\":\"Received time string: %s\"}", req.ScheduleTime)

	// Parse date and time with more flexible time format
	scheduleDate, err := time.Parse("2006-01-02", req.ScheduleDate)
	if err != nil {
		log.Printf("{\"severity\":\"ERROR\",\"message\":\"Invalid date format: %v\"}", err)
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// Try parsing time in HH:MM format first
	scheduleTime, err := time.Parse("15:04", req.ScheduleTime)
	if err != nil {
		log.Printf("{\"severity\":\"ERROR\",\"message\":\"Failed to parse time in HH:MM format: %v\"}", err)
		// Try HH:MM:SS format as fallback
		scheduleTime, err = time.Parse("15:04:05", req.ScheduleTime)
		if err != nil {
			log.Printf("{\"severity\":\"ERROR\",\"message\":\"Failed to parse time in both formats: %v\"}", err)
			http.Error(w, "Invalid time format. Use HH:MM or HH:MM:SS", http.StatusBadRequest)
			return
		}
	}

	// Create carpool object
	carpool := &models.Carpool{
		CreatorID:       "temp-creator", // TODO: Get from auth context
		Name:            req.Name,
		ScheduleDate:    scheduleDate,
		ScheduleTime:    scheduleTime,
		RecurringOption: req.RecurringOption,
		AvailableSeats:  req.AvailableSeats,
		Stops:           make([]models.Stop, 2), // Always 2 stops: start and destination
	}

	// Create start stop
	carpool.Stops[0] = models.Stop{
		Address:   req.StartAddress,
		StopOrder: 0,
		StopType:  "START",
	}

	// Create destination stop
	carpool.Stops[1] = models.Stop{
		Address:   req.DestinationAddress,
		StopOrder: 1,
		StopType:  "DESTINATION",
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
	carpoolID := params["id"]

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
	carpoolID := params["id"]

	err := h.carpoolRepo.DeleteCarPool(context.Background(), carpoolID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Carpool not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to delete carpool: %v", err), http.StatusInternalServerError)
		return
	}
	print("Carpool deleted successfully")
	w.WriteHeader(http.StatusNoContent)
}

func (h *CarPoolHandler) SearchCarPools(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
}
