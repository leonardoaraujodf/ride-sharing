package main

import (
	"encoding/json"
	"net/http"

	"ride-sharing/shared/contracts"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) error {
	return writeJSON(w, status, contracts.APIResponse{
		Error: &contracts.APIError{Code: code, Message: message},
	})
}
