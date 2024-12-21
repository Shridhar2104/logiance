// shipment/internal/repository/repository.go
package repository

import (
    "context"
    "fmt"
    "github.com/Shridhar2104/logilo/shipment/internal/config"
    "github.com/Shridhar2104/logilo/shipment/internal/courier/delhivery"
    "github.com/Shridhar2104/logilo/shipment/proto"
)

type ShipmentRepository interface {
    GetShippingRates(ctx context.Context, req *proto.ShippingRateRequest) ([]*proto.Rate, error)
}

type shipmentRepository struct {
    delhiveryClient *delhivery.Client
}

func NewShipmentRepository(cfg *config.Config) ShipmentRepository {
    return &shipmentRepository{
        delhiveryClient: delhivery.NewClient(&cfg.DelhiveryConfig),
    }
}

func (r *shipmentRepository) GetShippingRates(ctx context.Context, req *proto.ShippingRateRequest) ([]*proto.Rate, error) {
    // Convert weight from kg to grams
    weightInGrams := int(req.Package.WeightKg * 1000)

    // Create Delhivery request
    delhiveryReq := &delhivery.ShippingChargeRequest{
        MD:    "E", // Express delivery
        CGM:   weightInGrams,
        OPin:  parseInt(req.Origin.PostalCode),
        DPin:  parseInt(req.Destination.PostalCode),
        SS:    "Delivered", // Default status
    }

    // Get rates from Delhivery
    delhiveryResp, err := r.delhiveryClient.CalculateShippingCharges(ctx, delhiveryReq)
    if err != nil {
        return nil, fmt.Errorf("delhivery rate calculation failed: %w", err)
    }

    // Convert response to our rate format
    rates := []*proto.Rate{
        {
            CourierName:    "Delhivery",
            ServiceName:    "Express",
            Amount:        float32(delhiveryResp.TotalAmount),
            Currency:      "INR",
            EstimatedDays: 3, // You might want to make this dynamic based on pin codes
        },
    }

    return rates, nil
}

func parseInt(s string) int {
    var i int
    fmt.Sscanf(s, "%d", &i)
    return i
}