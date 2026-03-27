package repository

import (
	"context"
	"testing"
	"time"

	"ride-sharing/services/trip-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInmemRepository_CreateTrip_Success(t *testing.T) {
	ctx := context.Background()
	repo := NewInmemRepository()

	trip := &domain.TripModel{
		ID:     primitive.NewObjectID(),
		UserID: "user-123",
		Status: "pending",
		RideFare: &domain.RideFareModel{
			ID:                primitive.NewObjectID(),
			UserID:            "user-123",
			TotalPriceInCents: 1000,
			ExpiresAt:         time.Now().Add(time.Hour),
		},
	}

	result, err := repo.CreateTrip(ctx, trip)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if result.ID != trip.ID {
		t.Errorf("expected ID %s, got %s", trip.ID.Hex(), result.ID.Hex())
	}

	storedTrip, exists := repo.trips[trip.ID.Hex()]
	if !exists {
		t.Error("expected trip to be stored in repository")
	}

	if storedTrip.ID != trip.ID {
		t.Error("expected stored trip ID to match")
	}
}

func TestInmemRepository_CreateTrip_MultipleTrips(t *testing.T) {
	ctx := context.Background()
	repo := NewInmemRepository()

	trip1 := &domain.TripModel{
		ID:     primitive.NewObjectID(),
		UserID: "user-1",
		Status: "pending",
	}
	trip2 := &domain.TripModel{
		ID:     primitive.NewObjectID(),
		UserID: "user-2",
		Status: "pending",
	}

	result1, err := repo.CreateTrip(ctx, trip1)
	if err != nil {
		t.Fatalf("expected no error for trip1, got %v", err)
	}

	result2, err := repo.CreateTrip(ctx, trip2)
	if err != nil {
		t.Fatalf("expected no error for trip2, got %v", err)
	}

	if result1.ID == result2.ID {
		t.Error("expected different IDs for different trips")
	}

	if len(repo.trips) != 2 {
		t.Errorf("expected 2 trips in storage, got %d", len(repo.trips))
	}
}
