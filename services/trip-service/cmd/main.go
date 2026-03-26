package main

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"time"
)

func main() {
	ctx := context.Background()
	inmemRepo := repository.NewInmemRepository()
	srv := service.NewService(inmemRepo)
	fare := &domain.RideFareModel{
		UserID: "42",
	}
	t, err := srv.CreateTrip(ctx, fare)
	if err != nil {
		log.Println(err)
	}
	log.Println("Created trip:", t)
	// keep the program running for now
	for {
		time.Sleep(time.Second)
	}
}
