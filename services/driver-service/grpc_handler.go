package main

import (
	"context"
	pb "ride-sharing/shared/proto/driver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type driverGrpcHandler struct {
	pb.UnimplementedDriverServiceServer

	service *Service
}

func NewGrpcHandler(s *grpc.Server, service *Service) {
	handler := &driverGrpcHandler{
		service: service,
	}

	pb.RegisterDriverServiceServer(s, handler)
}

func (h *driverGrpcHandler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	driverId := req.GetDriverId()
	packageSlug := req.GetPackageSlug()
	driver, err := h.service.RegisterDriver(driverId, packageSlug)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register driver %s: %v", driverId, err)
	}
	return &pb.RegisterDriverResponse{
		Driver: driver,
	}, nil
}

func (h *driverGrpcHandler) UnregisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	driverId := req.GetDriverId()
	driver, err := h.service.UnregisterDriver(driverId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unregister driver %s, %v", driverId, err)
	}
	return &pb.RegisterDriverResponse{
		Driver: driver,
	}, nil
}
