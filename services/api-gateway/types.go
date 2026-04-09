package main

import (
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
