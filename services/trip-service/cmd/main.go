package main

import (
	"log"
	"net/http"
	h "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
)

func main() {
	inmemRepo := repository.NewInmemRepository()
	srv := service.NewService(inmemRepo)
	mux := http.NewServeMux()
	httphandler := &h.HttpHandler{Service: srv}
	mux.HandleFunc("/preview", httphandler.HandleTripPreview)

	server := &http.Server{
		Addr:    ":8083",
		Handler: mux,
	}

	log.Printf("Trip Service listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
