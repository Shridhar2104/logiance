package main

import (
	"context"
	"fmt"

	"github.com/Shridhar2104/logilo/graphql/models"
	pb "github.com/Shridhar2104/logilo/shipment/proto"
)

type ShippingResolver struct {
    shipmentClient pb.ShipmentServiceClient
}

func NewShippingResolver(client pb.ShipmentServiceClient) *ShippingResolver {
    return &ShippingResolver{
        shipmentClient: client,
    }
}

func (r *ShippingResolver) CalculateShippingRates(ctx context.Context, input ShippingRateInput) (*models.ShippingRateResponse, error) {
    req := &pb.RateRequest{
        OriginPincode:      fmt.Sprintf("%d", input.OriginPincode),
        DestinationPincode: fmt.Sprintf("%d", input.DestinationPincode),
        Weight:             input.Weight,
        Length:             float64(*input.Length),
        Width:              float64(*input.Width),
        Height:             float64(*input.Height),
        PaymentMode:        string(input.PaymentMode),
        CollectableAmount:  input.CollectableAmount,
        CourierCodes:       input.CourierCodes,
    }

    resp, err := r.shipmentClient.CalculateRates(ctx, req)
    if err != nil {
        return &models.ShippingRateResponse{
            Success: false,
            Error:   err.Error(),
        }, nil
    }

    rates := make([]*models.CourierRate, len(resp.Rates))
    for i, rate := range resp.Rates {
        rates[i] = &models.CourierRate{
            CourierCode:    rate.CourierCode,
            BaseCharge:     rate.BaseCharge,
            FuelSurcharge:  rate.FuelSurcharge,
            CodCharge:      rate.CodCharge,
            HandlingCharge: rate.HandlingCharge,
            TotalCharge:    rate.TotalCharge,
            ExpectedDays:   int(rate.ExpectedDays),
        }
    }

    return &models.ShippingRateResponse{
        Success: resp.Success,
        Rates:   rates,
        Error:   resp.Error,
    }, nil
}

func (r *ShippingResolver) GetAvailableCouriers(ctx context.Context, input AvailabilityInput) (*models.CourierAvailabilityResponse, error) {
    req := &pb.AvailabilityRequest{
        OriginPincode:      fmt.Sprintf("%d", input.OriginPincode),
        DestinationPincode: fmt.Sprintf("%d", input.DestinationPincode),
        Weight:             input.Weight,
        PaymentMode:        string(*input.PaymentMode),
    }

    resp, err := r.shipmentClient.GetAvailableCouriers(ctx, req)
    if err != nil {
        return &models.CourierAvailabilityResponse{
            Success: false,
            Error:   err.Error(),
        }, nil
    }

    couriers := make([]*models.CourierInfo, len(resp.AvailableCouriers))
    for i, c := range resp.AvailableCouriers {
        couriers[i] = &models.CourierInfo{
            CourierCode:        c.Code,
            CourierName:        c.Name,
            Description: c.Description,
        }
    }

    return &models.CourierAvailabilityResponse{
        Success:           resp.Success,
        AvailableCouriers: couriers,
        Error:            resp.Error,
    }, nil
}

func (r *ShippingResolver) CreateShipment(ctx context.Context, input CreateShipmentInput) (*ShipmentResponse, error) {
    // Convert addresses to protobuf format
    consignee := &pb.Address{
        Name:         input.Consignee.Name,
        CompanyName:  *input.Consignee.CompanyName,
        Phone:        input.Consignee.Phone,
        Email:        *input.Consignee.Email,
        AddressLine1: input.Consignee.AddressLine1,
        AddressLine2: *input.Consignee.AddressLine2,
        City:         input.Consignee.City,
        State:        input.Consignee.State,
        Pincode:      input.Consignee.Pincode,
        Gstin:        *input.Consignee.Gstin,
    }

    pickup := &pb.Address{
        Name:         input.Pickup.Name,
        CompanyName:  *input.Pickup.CompanyName,
        Phone:        input.Pickup.Phone,
        Email:        *input.Pickup.Email,
        AddressLine1: input.Pickup.AddressLine1,
        AddressLine2: *input.Pickup.AddressLine2,
        City:         input.Pickup.City,
        State:        input.Pickup.State,
        Pincode:      input.Pickup.Pincode,
        Gstin:        *input.Pickup.Gstin,
    }

    // Convert order items
    items := make([]*pb.OrderItem, len(input.Items))
    for i, item := range input.Items {
        items[i] = &pb.OrderItem{
            Sku:          item.Sku,
            Name:         item.Name,
            Quantity:     int32(item.Quantity),
            Price:        item.Price,
            HsnCode:      *item.HsnCode,
            Category:     *item.Category,
            ActualWeight: *item.ActualWeight,
        }
    }

    req := &pb.CreateShipmentRequest{
        CourierCode:       input.CourierCode,
        OrderNumber:       input.OrderNumber,
        PaymentType:       string(input.PaymentType),
        PackageWeight:     input.PackageWeight,
        PackageLength:     input.PackageLength,
        PackageBreadth:    input.PackageBreadth,
        PackageHeight:     input.PackageHeight,
        OrderAmount:       input.OrderAmount,
        CollectableAmount: input.CollectableAmount,
        Consignee:        consignee,
        Pickup:           pickup,
        Items:            items,
        AutoPickup:       *input.AutoPickup,
    }

    resp, err := r.shipmentClient.CreateShipment(ctx, req)
    if err != nil {
        return &ShipmentResponse{
            Success: false,
            Error:  &resp.Error,
        }, nil
    }

    return &ShipmentResponse{
        Success:    resp.Success,
        TrackingID: &resp.TrackingId,
        CourierAwb: &resp.CourierAwb,
        Label:      &resp.Label,
        Error:      &resp.Error,
    }, nil
}

func (r *ShippingResolver) TrackShipment(ctx context.Context, input TrackingInput) (*TrackingResponse, error) {
    req := &pb.TrackingRequest{
        CourierCode: input.CourierCode,
        TrackingId:  input.TrackingID,
    }

    resp, err := r.shipmentClient.TrackShipment(ctx, req)
    if err != nil {
        return &TrackingResponse{
            Success: false,
            Error:   &resp.Error,
        }, nil
    }

    events := make([]*TrackingEvent, len(resp.Events))
    for i, event := range resp.Events {
        events[i] = &TrackingEvent{
            Status:      event.Status,
            Location:    event.Location,
            Timestamp:   event.Timestamp,
            Description: event.Description,
        }
    }

    return &TrackingResponse{
        Success: resp.Success,
        Events:  events,
        Error:   &resp.Error,
    }, nil
}