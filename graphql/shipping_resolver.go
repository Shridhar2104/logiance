// shipping_resolver.go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Shridhar2104/logilo/graphql/models"
	pb "github.com/Shridhar2104/logilo/shipment/proto"
)

type ShippingResolver interface {
    CalculateShippingRates(ctx context.Context, input models.ShippingRateInput) (*models.ShippingRateResponse, error)
    GetAvailableCouriers(ctx context.Context, input models.AvailabilityInput) (*models.CourierAvailabilityResponse, error)
}


type shippingResolver struct {
    server *Server
}

func (r *shippingResolver) CalculateShippingRates(ctx context.Context, input models.ShippingRateInput) (*models.ShippingRateResponse, error) {
    log.Printf("Starting CalculateShippingRates with input: %+v", input)
    
    if r.server == nil {
        log.Printf("ERROR: Server is nil in shipping resolver")
        return nil, fmt.Errorf("server not initialized")
    }
    
    if r.server.shipmentClient == nil {
        log.Printf("ERROR: ShipmentClient is nil in shipping resolver")
        return nil, fmt.Errorf("shipment client not initialized")
    }

    log.Printf("Creating request with dimensions: %+v", input.Dimensions)
    
    dimensions := make([]*pb.Package, len(input.Dimensions))
    for i, dim := range input.Dimensions {
        dimensions[i] = &pb.Package{
            Length: float32(dim.Length),
            Width:  float32(dim.Width),
            Height: float32(dim.Height),
            Weight: float32(dim.Weight),
        }
    }

    req := &pb.RateRequest{
        OriginPincode:      int32(input.OriginPincode),
        DestinationPincode: int32(input.DestinationPincode),
        Weight:             int32(input.Weight),
        CourierCodes:       input.CourierCodes,
        PaymentMode:        string(input.PaymentMode),
        Dimensions:         dimensions,
    }

    log.Printf("Calling shipment service with request: %+v", req)
    
    resp, err := r.server.shipmentClient.CalculateRates(ctx, req)
    if err != nil {
        log.Printf("ERROR calling shipment service: %v", err)
        errMsg := fmt.Sprintf("failed to calculate rates: %v", err)
        return &models.ShippingRateResponse{
            Success: false,
            Error:   errMsg,
        }, nil
    }

    // Convert protobuf response to GraphQL model
    rates := make([]*models.CourierRate, len(resp.Rates))
    for i, rate := range resp.Rates {
        rates[i] = &models.CourierRate{
            CourierCode:    rate.CourierCode,
            CourierName:    rate.CourierName,
            ServiceType:    rate.ServiceType,
            EstimatedDays: int(rate.EstimatedDays),
            RateDetails: &models.RateDetails{
                TotalAmount:   rate.RateDetails.TotalAmount,
                GrossAmount:   rate.RateDetails.GrossAmount,
                TaxAmount:     rate.RateDetails.TaxAmount,
                CodCharges:    rate.RateDetails.CodCharges,
                FuelSurcharge: rate.RateDetails.FuelSurcharge,
            },
        }
    }

    return &models.ShippingRateResponse{
        Success: resp.Success,
        Rates:   rates,
        Error:   resp.Error,
    }, nil
}

func (r *shippingResolver) GetAvailableCouriers(ctx context.Context, input models.AvailabilityInput) (*models.CourierAvailabilityResponse, error) {
    req := &pb.AvailabilityRequest{
        OriginPincode:      int32(input.OriginPincode),
        DestinationPincode: int32(input.DestinationPincode),
    }

    resp, err := r.server.shipmentClient.GetAvailableCouriers(ctx, req)
    if err != nil {
        return &models.CourierAvailabilityResponse{
            Success: false,
            Error:   err.Error(),
        }, nil
    }

    couriers := make([]*models.CourierInfo, len(resp.AvailableCouriers))
    for i, courier := range resp.AvailableCouriers {
        couriers[i] = &models.CourierInfo{
            CourierCode:       courier.CourierCode,
            CourierName:       courier.CourierName,
            SupportedServices: courier.SupportedServices,
        }
    }

    return &models.CourierAvailabilityResponse{
        Success:           resp.Success,
        AvailableCouriers: couriers,
        Error:            resp.Error,
    }, nil
}

func validateShippingInput(input models.ShippingRateInput) error {
    if input.OriginPincode < 100000 || input.OriginPincode > 999999 {
        return fmt.Errorf("invalid origin pincode")
    }
    if input.DestinationPincode < 100000 || input.DestinationPincode > 999999 {
        return fmt.Errorf("invalid destination pincode")
    }
    if input.Weight <= 0 {
        return fmt.Errorf("weight must be positive")
    }
    return nil
}