package repository

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/internal/domain"
	pbd "ride-sharing/shared/proto/driver"
	pb "ride-sharing/shared/proto/trip"
	"sync"
)

type inmemRepository struct {
	trips     map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
	sync.RWMutex
}

func NewInmemRepository() *inmemRepository {
	return &inmemRepository{
		trips:     make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}

func (r *inmemRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.Lock()
	defer r.Unlock()
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}

func (r *inmemRepository) SaveRideFare(ctx context.Context, f *domain.RideFareModel) error {
	r.Lock()
	defer r.Unlock()
	r.rideFares[f.ID.Hex()] = f
	return nil
}

func (r *inmemRepository) GetRideFareByID(ctx context.Context, id string) (*domain.RideFareModel, error) {
	r.RLock()
	defer r.RUnlock()
	fare, exist := r.rideFares[id]
	if !exist {
		return nil, fmt.Errorf("fare does not exist with ID: %s", id)
	}
	return fare, nil
}

func (r *inmemRepository) GetTripByID(ctx context.Context, id string) (*domain.TripModel, error) {
	r.RLock()
	defer r.RUnlock()
	trip, exist := r.trips[id]
	if !exist {
		return nil, fmt.Errorf("trip does not exist with ID: %s", id)
	}
	return trip, nil
}

func (r *inmemRepository) UpdateTrip(ctx context.Context, tripID string, status string, driver *pbd.Driver) error {
	r.Lock()
	defer r.Unlock()
	trip, exist := r.trips[tripID]
	if !exist {
		return fmt.Errorf("trip does not exist with ID: %s", tripID)
	}
	trip.Status = status
	if driver != nil {
		trip.Driver = &pb.TripDriver{
			Id:             driver.Id,
			Name:           driver.Name,
			ProfilePicture: driver.ProfilePicture,
			CarPlate:       driver.CarPlate,
		}
	}
	return nil
}
