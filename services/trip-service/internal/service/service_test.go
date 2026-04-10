package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"ride-sharing/services/trip-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockTripRepository struct {
	trips map[string]*domain.TripModel
	err   error
}

func (m *mockTripRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.trips[trip.ID.Hex()] = trip
	return trip, nil
}

func TestCreateTrip_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mockTripRepository{trips: make(map[string]*domain.TripModel)}
	svc := NewService(mockRepo, "")

	fare := &domain.RideFareModel{
		ID:                primitive.NewObjectID(),
		UserID:            "user-123",
		PackageSlug:       "basic",
		TotalPriceInCents: 1000,
		ExpiresAt:         time.Now().Add(time.Hour),
	}

	trip, err := svc.CreateTrip(ctx, fare)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if trip == nil {
		t.Fatal("expected trip, got nil")
	}

	if trip.UserID != fare.UserID {
		t.Errorf("expected UserID %s, got %s", fare.UserID, trip.UserID)
	}

	if trip.Status != "pending" {
		t.Errorf("expected status 'pending', got %s", trip.Status)
	}

	if trip.RideFare != fare {
		t.Error("expected RideFare to be the same reference")
	}
}

func TestCreateTrip_RepositoryError(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mockTripRepository{
		trips: make(map[string]*domain.TripModel),
		err:   errors.New("database error"),
	}
	svc := NewService(mockRepo, "")

	fare := &domain.RideFareModel{
		ID:        primitive.NewObjectID(),
		UserID:    "user-123",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	_, err := svc.CreateTrip(ctx, fare)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("expected error 'database error', got %s", err.Error())
	}
}
