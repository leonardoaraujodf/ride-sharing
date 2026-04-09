package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
	"time"
)

var tripServiceURL = getTripServiceURL()

func getTripServiceURL() string {
	if url := os.Getenv("TRIP_SERVICE_URL"); url != "" {
		return url
	}
	return "http://trip-service:8083/preview"
}

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENABLE_DELAY") == "true" {
		time.Sleep(2 * time.Second)
	}

	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeJSONError(w, http.StatusBadRequest, "DECODE_ERROR", "failed to parse JSON data")
		return
	}

	defer r.Body.Close()

	if reqBody.UserID == "" {
		writeJSONError(w, http.StatusBadRequest, "VALIDATION_ERROR", "user ID is required")
		return
	}

	// Why we need to create a new client for each connection?
	// Because if a service is down, we don't want to block the whole application
	// so we createa a new client for each connection
	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer tripService.Close()

	tripPreview, err := tripService.Client.PreviewTrip(r.Context(), reqBody.toProto())
	if err != nil {
		log.Printf("Failed to preview a trip: %v", err)
		http.Error(w, "Failed to preview a trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: toPreviewTripResponse(tripPreview)}
	writeJSON(w, http.StatusCreated, response)
}

func handleTripStart(w http.ResponseWriter, r *http.Request) {
	var reqBody createTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeJSONError(w, http.StatusBadRequest, "DECODE_ERROR", "failed to parse JSON data")
		return
	}

	defer r.Body.Close()
	if reqBody.UserID == "" {
		writeJSONError(w, http.StatusBadRequest, "VALIDATION_ERROR", "user ID is required")
		return
	}

	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer tripService.Close()

	trip, err := tripService.Client.CreateTrip(r.Context(), reqBody.toProto())
	if err != nil {
		log.Printf("Failed to create a trip: %v", err)
		http.Error(w, "Failed to create a trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: toCreateTripResponse(trip)}
	writeJSON(w, http.StatusCreated, response)
}
