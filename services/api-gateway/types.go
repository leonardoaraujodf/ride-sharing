package main

import (
	pbDriver "ride-sharing/shared/proto/driver"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"
)

// Response DTOs with camelCase JSON tags so the frontend receives consistent field names.

type coordinateResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type geometryResponse struct {
	Coordinates []coordinateResponse `json:"coordinates"`
}

type routeResponse struct {
	Geometry []geometryResponse `json:"geometry"`
	Distance float64            `json:"distance"`
	Duration float64            `json:"duration"`
}

type rideFareResponse struct {
	ID                string  `json:"id"`
	PackageSlug       string  `json:"packageSlug"`
	TotalPriceInCents float64 `json:"totalPriceInCents"`
}

type previewTripResponse struct {
	TripID    string             `json:"tripID"`
	Route     routeResponse      `json:"route"`
	RideFares []rideFareResponse `json:"rideFares"`
}

type createTripResponse struct {
	TripID string `json:"tripID"`
}

func toPreviewTripResponse(r *pb.PreviewTripResponse) previewTripResponse {
	rideFares := make([]rideFareResponse, len(r.RideFares))
	for i, f := range r.RideFares {
		rideFares[i] = rideFareResponse{
			ID:                f.Id,
			PackageSlug:       f.PackageSlug,
			TotalPriceInCents: f.TotalPriceInCents,
		}
	}

	geometries := make([]geometryResponse, len(r.Route.Geometry))
	for i, g := range r.Route.Geometry {
		coords := make([]coordinateResponse, len(g.Coordinates))
		for j, c := range g.Coordinates {
			coords[j] = coordinateResponse{Latitude: c.Latitude, Longitude: c.Longitude}
		}
		geometries[i] = geometryResponse{Coordinates: coords}
	}

	return previewTripResponse{
		TripID: r.TripId,
		Route: routeResponse{
			Geometry: geometries,
			Distance: r.Route.Distance,
			Duration: r.Route.Duration,
		},
		RideFares: rideFares,
	}
}

func toCreateTripResponse(r *pb.CreateTripResponse) createTripResponse {
	return createTripResponse{
		TripID: r.TripId,
	}
}

type previewTripRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func (p *previewTripRequest) toProto() *pb.PreviewTripRequest {
	return &pb.PreviewTripRequest{
		UserId: p.UserID,
		StartLocation: &pb.Coordinate{
			Latitude:  p.Pickup.Latitude,
			Longitude: p.Pickup.Longitude,
		},
		EndLocation: &pb.Coordinate{
			Latitude:  p.Destination.Latitude,
			Longitude: p.Destination.Longitude,
		},
	}
}

type createTripRequest struct {
	UserID     string `json:"userID"`
	RideFareID string `json:"rideFareID"`
}

func (s *createTripRequest) toProto() *pb.CreateTripRequest {
	return &pb.CreateTripRequest{
		UserId:     s.UserID,
		RideFareId: s.RideFareID,
	}
}

// driverWsResponse is the camelCase DTO sent to the driver WebSocket on registration.
type driverWsResponse struct {
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	ProfilePicture string             `json:"profilePicture"`
	CarPlate       string             `json:"carPlate"`
	GeoHash        string             `json:"geohash"`
	Location       coordinateResponse `json:"location"`
}

func toDriverWsResponse(d *pbDriver.Driver) driverWsResponse {
	loc := coordinateResponse{}
	if d.Location != nil {
		loc = coordinateResponse{Latitude: d.Location.Latitude, Longitude: d.Location.Longitude}
	}
	return driverWsResponse{
		ID:             d.Id,
		Name:           d.Name,
		ProfilePicture: d.ProfilePicture,
		CarPlate:       d.CarPlate,
		GeoHash:        d.GeoHash,
		Location:       loc,
	}
}

// tripDriverWsResponse is the camelCase DTO for the driver nested inside a trip.
type tripDriverWsResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profilePicture"`
	CarPlate       string `json:"carPlate"`
}

// tripWsResponse is the camelCase DTO for a trip sent via WebSocket.
type tripWsResponse struct {
	ID           string                `json:"id"`
	UserID       string                `json:"userID"`
	Status       string                `json:"status"`
	SelectedFare rideFareResponse      `json:"selectedFare"`
	Route        routeResponse         `json:"route"`
	Driver       *tripDriverWsResponse `json:"driver,omitempty"`
}

// tripWsPayload wraps the trip in the envelope the frontend expects.
type tripWsPayload struct {
	Trip tripWsResponse `json:"trip"`
}

func toTripWsPayload(t *pb.Trip) tripWsPayload {
	fare := rideFareResponse{}
	if t.SelectedFare != nil {
		fare = rideFareResponse{
			ID:                t.SelectedFare.Id,
			PackageSlug:       t.SelectedFare.PackageSlug,
			TotalPriceInCents: t.SelectedFare.TotalPriceInCents,
		}
	}

	geometries := []geometryResponse{}
	if t.Route != nil {
		geometries = make([]geometryResponse, len(t.Route.Geometry))
		for i, g := range t.Route.Geometry {
			coords := make([]coordinateResponse, len(g.Coordinates))
			for j, c := range g.Coordinates {
				coords[j] = coordinateResponse{Latitude: c.Latitude, Longitude: c.Longitude}
			}
			geometries[i] = geometryResponse{Coordinates: coords}
		}
	}

	route := routeResponse{}
	if t.Route != nil {
		route = routeResponse{
			Geometry: geometries,
			Distance: t.Route.Distance,
			Duration: t.Route.Duration,
		}
	}

	var driver *tripDriverWsResponse
	if t.Driver != nil {
		driver = &tripDriverWsResponse{
			ID:             t.Driver.Id,
			Name:           t.Driver.Name,
			ProfilePicture: t.Driver.ProfilePicture,
			CarPlate:       t.Driver.CarPlate,
		}
	}

	return tripWsPayload{
		Trip: tripWsResponse{
			ID:           t.Id,
			UserID:       t.UserId,
			Status:       t.Status,
			SelectedFare: fare,
			Route:        route,
			Driver:       driver,
		},
	}
}
