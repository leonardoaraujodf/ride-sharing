package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
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

	jsonBody, _ := json.Marshal(reqBody)
	reader := bytes.NewReader(jsonBody)

	resp, err := http.Post("http://trip-service:8083/preview", "application/json", reader)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "TRIP_SERVICE_ERROR", "failed to call trip service")
		return
	}
	defer resp.Body.Close()

	var respBody any
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "DECODE_ERROR", "failed to parse trip service response")
		return
	}

	// TODO: Call trip service
	response := contracts.APIResponse{Data: respBody}
	writeJSON(w, http.StatusCreated, response)
}
