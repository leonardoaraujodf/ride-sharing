package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ride-sharing/shared/contracts"
)

func TestHandleTripPreview_Success(t *testing.T) {
	mockTripService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
	}))
	defer mockTripService.Close()

	originalTripServiceURL := tripServiceURL
	tripServiceURL = mockTripService.URL
	defer func() { tripServiceURL = originalTripServiceURL }()

	body := []byte(`{"userID": "123"}`)
	req := httptest.NewRequest(http.MethodPost, "/trip/preview", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handleTripPreview(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var resp contracts.APIResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data == nil {
		t.Error("expected data in response, got nil")
	}
}

func TestHandleTripPreview_EmptyUser(t *testing.T) {
	body := []byte(`{"userID": ""}`)
	req := httptest.NewRequest(http.MethodPost, "/trip/preview", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handleTripPreview(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var resp contracts.APIResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error == nil {
		t.Error("expected error in response")
	}
}

func TestHandleTripPreview_DecodeError(t *testing.T) {
	body := []byte(`{invalid json}`)
	req := httptest.NewRequest(http.MethodPost, "/trip/preview", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handleTripPreview(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var resp contracts.APIResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error == nil {
		t.Error("expected error in response")
	}
}
