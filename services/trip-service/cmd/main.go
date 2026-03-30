package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	h "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"syscall"
	"time"
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

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Trip Service listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("Error starting the server: %v", err)
	case sig := <-shutdown:
		log.Printf("Server is shutting down due to signal: %v", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Could not stop the server gracefully: %v", err)
			server.Close()
		} else {
			log.Println("Server shutdown gracefully")
		}
	}
}
