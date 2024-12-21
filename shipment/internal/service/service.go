// internal/service/service.go
package service

import (
    "context"

    "github.com/google/uuid"
    "github.com/Shridhar2104/logilo/shipment/internal/repository"
    "github.com/Shridhar2104/logilo/shipment/proto"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type ShipmentService struct {
    proto.UnimplementedShipmentServiceServer
    repo repository.ShipmentRepository
}

func NewShipmentService(repo repository.ShipmentRepository) *ShipmentService {
    return &ShipmentService{
        repo: repo,
    }
}

func (s *ShipmentService) CalculateShippingRate(ctx context.Context, req *proto.ShippingRateRequest) (*proto.ShippingRateResponse, error) {
    // Validate request
    if err := validateRequest(req); err != nil {
        return nil, status.Error(codes.InvalidArgument, err.Error())
    }

    // Generate unique request ID
    requestID := uuid.New().String()

    // Get rates from repository (which will interact with courier APIs)
    rates, err := s.repo.GetShippingRates(ctx, req)
    if err != nil {
        return &proto.ShippingRateResponse{
            RequestId:    requestID,
            ErrorMessage: err.Error(),
        }, nil
    }

    return &proto.ShippingRateResponse{
        RequestId: requestID,
        Rates:    rates,
    }, nil
}

func validateRequest(req *proto.ShippingRateRequest) error {
    // Add validation logic here
    // For example: check if addresses are complete, package dimensions are valid, etc.
    return nil
}