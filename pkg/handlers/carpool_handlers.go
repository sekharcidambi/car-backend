package handlers

import (
	"bytes"
	"car-backend/pkg/models"
	"car-backend/pkg/repository"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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
	// Debug: Print request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("{\"severity\":\"ERROR\",\"message\":\"Failed to read request body: %v\"}", err)
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}
	log.Printf("{\"severity\":\"DEBUG\",\"message\":\"Received request body: %s\"}", string(body))
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore body for later use

	var req models.CreateCarPoolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("{\"severity\":\"ERROR\",\"message\":\"Failed to decode request: %v\"}", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse date and time
	scheduleDate, err := time.Parse("2006-01-02", req.ScheduleDate)
	if err != nil {
		log.Printf("{\"severity\":\"ERROR\",\"message\":\"Invalid date format: %v\"}", err)
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	scheduleTime, err := time.Parse("15:04:05", req.ScheduleTime)
	if err != nil {
		log.Printf("{\"severity\":\"ERROR\",\"message\":\"Invalid time format: %v\"}", err)
		http.Error(w, "Invalid time format", http.StatusBadRequest)
		return
	}

	carpool := &models.CarPool{
		ScheduleDate:    scheduleDate,
		ScheduleTime:    scheduleTime,
		RecurringOption: req.RecurringOption,
		StartPointLat:   req.StartPointLat,
		StartPointLng:   req.StartPointLng,
		DestinationLat:  req.DestinationLat,
		DestinationLng:  req.DestinationLng,
		AvailableSeats:  req.AvailableSeats,
		MusicPreference: req.MusicPreference,
		SmokingAllowed:  req.SmokingAllowed,
		PetsAllowed:     req.PetsAllowed,
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
