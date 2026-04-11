package main

import (
	"fmt"
	math "math/rand/v2"
	pb "ride-sharing/shared/proto/driver"
	"ride-sharing/shared/util"
	"slices"
	"sync"

	"github.com/mmcloughlin/geohash"
)

type driverInMap struct {
	Driver *pb.Driver
	// TODO: route
}

type Service struct {
	drivers []*driverInMap
	sync.RWMutex
}

func NewService() *Service {
	return &Service{
		drivers: []*driverInMap{},
	}
}

func (s *Service) FindAvailableDrivers(packageType string) []string {
	var matchingDrivers []string

	for _, driver := range s.drivers {
		if driver.Driver.PackageSlug == packageType {
			matchingDrivers = append(matchingDrivers, driver.Driver.Id)
		}
	}

	if len(matchingDrivers) == 0 {
		return []string{}
	}

	return matchingDrivers
}

func (s *Service) RegisterDriver(driverId string, packageSlug string) (*pb.Driver, error) {
	s.Lock()
	defer s.Unlock()

	if driverId == "" {
		return nil, fmt.Errorf("invalid driver id: %s", driverId)
	}

	if packageSlug == "" {
		return nil, fmt.Errorf("invalid package slug: %s", packageSlug)
	}

	randomIndex := math.IntN(len(PredefinedRoutes))
	randomRoute := PredefinedRoutes[randomIndex]

	// we can ignore this property for now, but it must be sent to the frontend
	geohash := geohash.Encode(randomRoute[0][0], randomRoute[0][1])

	driver := &pb.Driver{
		Id:             driverId,
		Name:           "Lando Norris",
		ProfilePicture: util.GetRandomAvatar(1),
		CarPlate:       GenerateRandomPlate(),
		GeoHash:        geohash,
		PackageSlug:    packageSlug,
		Location:       &pb.Location{Latitude: randomRoute[0][0], Longitude: randomRoute[0][1]},
	}
	s.drivers = append(s.drivers, &driverInMap{Driver: driver})
	return driver, nil
}

func (s *Service) UnregisterDriver(driverId string) (*pb.Driver, error) {
	s.Lock()
	defer s.Unlock()

	if driverId == "" {
		return nil, fmt.Errorf("invalid driver id %s", driverId)
	}

	for idx, driver := range s.drivers {
		if driver.Driver.Id == driverId {
			s.drivers = slices.Delete(s.drivers, idx, idx+1)
			return driver.Driver, nil
		}
	}

	return nil, fmt.Errorf("could not find driver id %s in driver's list", driverId)
}
