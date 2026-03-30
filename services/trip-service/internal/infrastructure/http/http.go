package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"
	"time"
)

type HttpHandler struct {
	Service domain.TripService
}

type previewTripRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func (h *HttpHandler) HandleTripPreview(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENABLE_DELAY") == "true" {
		time.Sleep(3 * time.Second)
	}

	fmt.Println("Received request for trip preview")
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	fmt.Printf("pickup, latitude: %f", reqBody.Pickup.Latitude)
	fmt.Printf("pickup, longitude: %f", reqBody.Pickup.Longitude)
	fmt.Printf("destination, latitude: %f", reqBody.Destination.Latitude)
	fmt.Printf("destination, longitude: %f", reqBody.Destination.Longitude)

	ctx := r.Context()
	trip, err := h.Service.GetRoute(ctx, &reqBody.Pickup, &reqBody.Destination)
	if err != nil {
		http.Error(w, "failed to create trip", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, trip)
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
