// Implement Provider interface methods...

// internal/service/shipping.go
package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/Shridhar2104/logilo/shipment/internal/config"
    "github.com/Shridhar2104/logilo/shipment/internal/database"
	"github.com/Shridhar2104/logilo/shipment/internal/courier"
	"github.com/Shridhar2104/logilo/shipment/internal/courier/bluedart"
	"github.com/Shridhar2104/logilo/shipment/internal/courier/delhivery"
	"github.com/Shridhar2104/logilo/shipment/internal/courier/xpressbees"
	"github.com/Shridhar2104/logilo/shipment/internal/ratelimit"
	pb "github.com/Shridhar2104/logilo/shipment/proto/proto"
)

type ShipmentService struct {
    pb.UnimplementedShipmentServiceServer
    cfg          *config.Config
    rateLimiters map[string]*ratelimit.RateLimiter
    providers    map[string]courier.CourierProvider
    mu          sync.RWMutex
    trackingService *ShipmentTrackingService
}
// internal/service/shipping.go
func NewShipmentService(cfg *config.Config,db *database.ShipmentDB) (*ShipmentService, error) {
    s := &ShipmentService{
        cfg:          cfg,
        rateLimiters: make(map[string]*ratelimit.RateLimiter),
        providers:    make(map[string]courier.CourierProvider),
        trackingService: NewShipmentTrackingService(db),

    }

    // Register Xpressbees
    xbProvider, err := xpressbees.NewProvider(cfg.GetCourierConfig("XPRESSBEES"))
    if err != nil {
        return nil, fmt.Errorf("failed to initialize Xpressbees provider: %w", err)
    }
    s.registerProvider("XPRESSBEES", xbProvider)

    // Register Delhivery
    dlProvider := delhivery.NewProvider(cfg.GetCourierConfig("DELHIVERY"))
    if err != nil {
        return nil, fmt.Errorf("failed to initialize Delhivery provider: %w", err)
    }
    s.registerProvider("DELHIVERY", dlProvider)

    // Register Bluedart
    bdProvider := bluedart.NewProvider(cfg.GetCourierConfig("BLUEDART"))
    s.registerProvider("BLUEDART", bdProvider)

    return s, nil
}
func (s *ShipmentService) registerProvider(code string, provider courier.CourierProvider) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.providers[code] = provider
    s.rateLimiters[code] = ratelimit.NewRateLimiter(s.cfg.GetRateLimit(code))
}

func (s *ShipmentService) CalculateRates(ctx context.Context, req *pb.RateRequest) (*pb.MultiRateResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var wg sync.WaitGroup
    ratesChan := make(chan *pb.CourierRate, len(s.providers))
    errorsChan := make(chan error, len(s.providers))

    // If specific couriers requested, use only those
    providers := s.providers
    if len(req.CourierCodes) > 0 {
        providers = make(map[string]courier.CourierProvider)
        for _, code := range req.CourierCodes {
            if provider, exists := s.providers[code]; exists {
                providers[code] = provider
            }
        }
    }

    // Calculate rates from all providers concurrently
    for code, provider := range providers {
        wg.Add(1)
        go func(code string, provider courier.CourierProvider) {
            defer wg.Done()

            // Check rate limit
            if !s.rateLimiters[code].Allow() {
                errorsChan <- fmt.Errorf("rate limit exceeded for %s", code)
                return
            }

            rate, err := provider.CalculateRate(ctx, &courier.RateRequest{
                OriginPincode:      req.OriginPincode,
                DestinationPincode: req.DestinationPincode,
                Weight:             req.Weight,
                Length:             req.Length,
                Width:              req.Width,
                Height:             req.Height,
                PaymentMode:        req.PaymentMode,
                CollectableAmount:  req.CollectableAmount,
            })
            
            if err != nil {
                errorsChan <- fmt.Errorf("%s: %v", code, err)
                return
            }

            ratesChan <- &pb.CourierRate{
                CourierCode:    code,
                BaseCharge:     rate.BaseCharge,
                FuelSurcharge:  rate.FuelSurcharge,
                CodCharge:      rate.CODCharge,
                HandlingCharge: rate.HandlingCharge,
                TotalCharge:    rate.TotalCharge,
                ExpectedDays:   int32(rate.ExpectedDays),
            }
        }(code, provider)
    }

    // Wait for all goroutines to complete
    wg.Wait()
    close(ratesChan)
    close(errorsChan)

    // Collect results
    var rates []*pb.CourierRate
    var errors []string

    for rate := range ratesChan {
        rates = append(rates, rate)
    }

    for err := range errorsChan {
        errors = append(errors, err.Error())
    }

    // Return results based on success/failure
    if len(rates) > 0 {
        return &pb.MultiRateResponse{
            Success: true,
            Rates:   rates,
            Error:   strings.Join(errors, "; "),
        }, nil
    }

    if len(errors) > 0 {
        return &pb.MultiRateResponse{
            Success: false,
            Error:   strings.Join(errors, "; "),
        }, nil
    }

    return &pb.MultiRateResponse{
        Success: true,
        Rates:   rates,
    }, nil
}

func (s *ShipmentService) GetAvailableCouriers(ctx context.Context, req *pb.AvailabilityRequest) (*pb.CourierListResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var wg sync.WaitGroup
    availableChan := make(chan *pb.CourierInfo, len(s.providers))
    errorsChan := make(chan error, len(s.providers))

    for code, provider := range s.providers {
        wg.Add(1)
        go func(code string, provider courier.CourierProvider) {
            defer wg.Done()

            if !s.rateLimiters[code].Allow() {
                errorsChan <- fmt.Errorf("rate limit exceeded for %s", code)
                return
            }

            available, err := provider.CheckServiceability(ctx, req.OriginPincode, req.DestinationPincode, req.Weight)
            if err != nil {
                errorsChan <- fmt.Errorf("%s: %v", code, err)
                return
            }

            if available {
                info := provider.GetProviderInfo()
                availableChan <- &pb.CourierInfo{
                    Code:        info.Code,
                    Name:        info.Name,
                    Description: info.Description,
                }
            }
        }(code, provider)
    }

    wg.Wait()
    close(availableChan)
    close(errorsChan)

    var availableCouriers []*pb.CourierInfo
    var errors []string

    for courier := range availableChan {
        availableCouriers = append(availableCouriers, courier)
    }

    for err := range errorsChan {
        errors = append(errors, err.Error())
    }

    return &pb.CourierListResponse{
        Success:           len(errors) == 0,
        AvailableCouriers: availableCouriers,
        Error:            strings.Join(errors, "; "),
    }, nil
}
func (s *ShipmentService) CreateShipment(ctx context.Context, req *pb.CreateShipmentRequest) (*pb.ShipmentResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
 
    provider, exists := s.providers[req.CourierCode]
    if !exists {
        return nil, fmt.Errorf("unknown courier code: %s", req.CourierCode)
    }
 
    if !s.rateLimiters[req.CourierCode].Allow() {
        return nil, fmt.Errorf("rate limit exceeded for %s", req.CourierCode)
    }

    log.Printf("Processing CreateShipment request: %+v", req)
 
    shipment, err := provider.CreateShipment(ctx, &courier.ShipmentRequest{
        OrderNumber:       req.OrderNumber,
        PaymentType:       req.PaymentType,
        PackageWeight:     req.PackageWeight,
        PackageLength:     req.PackageLength,
        PackageBreadth:    req.PackageBreadth,
        PackageHeight:     req.PackageHeight,
        OrderAmount:       req.OrderAmount,
        CollectableAmount: req.CollectableAmount,
        CODCharges:        req.CollectableAmount,
        AutoPickup:        false,

        Consignee: courier.Address{  // No pointer
            Name:         req.Consignee.Name,
            CompanyName:  req.Consignee.CompanyName,
            AddressLine1: req.Consignee.AddressLine1,
            AddressLine2: req.Consignee.AddressLine2,
            City:         req.Consignee.City,
            State:        req.Consignee.State,
            Pincode:      req.Consignee.Pincode,
            Phone:        req.Consignee.Phone,
            Email:        req.Consignee.Email,
        },

        Pickup: courier.Address{  // No pointer
            Name:         req.Pickup.Name,
            CompanyName:  req.Pickup.CompanyName,
            AddressLine1: req.Pickup.AddressLine1,
            AddressLine2: req.Pickup.AddressLine2,
            City:         req.Pickup.City,
            State:        req.Pickup.State,
            Pincode:      req.Pickup.Pincode,
            Phone:        req.Pickup.Phone,
            GSTIN:        req.Pickup.Gstin,
        },

        Items: mapOrderItems(req.Items),
    })
 
    if err != nil {
        return &pb.ShipmentResponse{
            Success: false,
            Error:   err.Error(),
        }, nil
    }
        // Save tracking information
        err = s.trackingService.SaveNewShipment(
            ctx,
            req.AccountId,  // Make sure to add this field to your proto definition
            req.OrderNumber,
            shipment.TrackingID,
            shipment.AWBNumber,
            req.CourierCode,
            shipment.Label,
        )
        if err != nil {
            log.Printf("Error saving tracking info: %v", err)
            // Continue with the response even if tracking save fails
            // You might want to implement a retry mechanism or handle this differently
        }
    
    return &pb.ShipmentResponse{
        Success:    shipment.Success,
        TrackingId: shipment.TrackingID,
        CourierAwb: shipment.AWBNumber,
        Label:      shipment.Label,
        Error:      shipment.Error,
    }, nil
}
 // Helper function to map order items
 func mapOrderItems(items []*pb.OrderItem) []courier.OrderItem {
    if len(items) == 0 {
        return nil
    }
 
    result := make([]courier.OrderItem, len(items))
    for i, item := range items {
        result[i] = courier.OrderItem{
            Name:     item.Name,
            SKU:      item.Sku,
            Quantity: int(item.Quantity),
            Price:    item.Price,
        }
    }
    return result
 }
func (s *ShipmentService) TrackShipment(ctx context.Context, req *pb.TrackingRequest) (*pb.TrackingResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    provider, exists := s.providers[req.CourierCode]
    if !exists {
        return nil, fmt.Errorf("unknown courier code: %s", req.CourierCode)
    }

    if !s.rateLimiters[req.CourierCode].Allow() {
        return nil, fmt.Errorf("rate limit exceeded for %s", req.CourierCode)
    }

    tracking, err := provider.TrackShipment(ctx, req.TrackingId)
    if err != nil {
        return &pb.TrackingResponse{
            Success: false,
            Error:   err.Error(),
        }, nil
    }

    var events []*pb.TrackingEvent
    for _, event := range tracking {
        events = append(events, &pb.TrackingEvent{
            Status:      event.Status,
            Location:    event.Location,
            Timestamp:   event.ActivityTime,
            Description: event.Description,
        })
    }

    return &pb.TrackingResponse{
        Success: true,
        Events:  events,
    }, nil
}

// Add method to get tracking by order ID
func (s *ShipmentService) GetShipmentByOrder(ctx context.Context, req *pb.OrderTrackingRequest) (*pb.ShipmentResponse, error) {
    tracking, err := s.trackingService.GetShipmentByOrder(ctx, req.OrderId)
    if err != nil {
        return nil, fmt.Errorf("error getting shipment by order: %w", err)
    }

    if tracking == nil {
        return &pb.ShipmentResponse{
            Success: false,
            Error:   "no shipment found for order",
        }, nil
    }

    return &pb.ShipmentResponse{
        Success:    true,
        TrackingId: tracking.TrackingID,
        CourierAwb: tracking.AWBNumber,
        Label:      tracking.Label,
    }, nil
}

// Add new tracking-related methods
func (s *ShipmentService) GetShipmentTracking(ctx context.Context, req *pb.TrackingRequest) (*pb.TrackingResponse, error) {
    // First check the database for existing tracking info
    tracking, _, err := s.trackingService.GetShipmentDetails(ctx, req.TrackingId)
    if err != nil {
        return nil, fmt.Errorf("error getting tracking details: %w", err)
    }

    // If found in database, use the courier code from there
    courierCode := req.CourierCode
    if tracking != nil {
        courierCode = tracking.CourierCode
    }

    // Get latest tracking info from courier
    provider, exists := s.providers[courierCode]
    if !exists {
        return nil, fmt.Errorf("unknown courier code: %s", courierCode)
    }

    if !s.rateLimiters[courierCode].Allow() {
        return nil, fmt.Errorf("rate limit exceeded for %s", courierCode)
    }

    // Get tracking events from courier
    courierEvents, err := provider.TrackShipment(ctx, req.TrackingId)
    if err != nil {
        return &pb.TrackingResponse{
            Success: false,
            Error:   err.Error(),
        }, nil
    }

    // Update database with new events
    for _, event := range courierEvents {
        err = s.trackingService.UpdateShipmentStatus(
            ctx,
            req.TrackingId,
            event.Status,
            event.Location,
            event.Description,
        )
        if err != nil {
            log.Printf("Error updating tracking status: %v", err)
        }
    }

    // Convert tracking events to proto format
    pbEvents := make([]*pb.TrackingEvent, len(courierEvents))
    for i, event := range courierEvents {
        pbEvents[i] = &pb.TrackingEvent{
            Status:      event.Status,
            Location:    event.Location,
            Timestamp:   event.ActivityTime,
            Description: event.Description,
        }
    }

    return &pb.TrackingResponse{
        Success: true,
        Events:  pbEvents,
    }, nil
}


// Add method to get all shipments for an account
func (s *ShipmentService) GetAccountShipments(ctx context.Context, req *pb.AccountShipmentsRequest) (*pb.AccountShipmentsResponse, error) {
    shipments, err := s.trackingService.GetShipmentsByAccount(ctx, req.AccountId, int(req.Page), int(req.PageSize))
    if err != nil {
        return nil, fmt.Errorf("error getting account shipments: %w", err)
    }

    pbShipments := make([]*pb.ShipmentInfo, len(shipments))
    for i, shipment := range shipments {
        pbShipments[i] = &pb.ShipmentInfo{
            OrderNumber: shipment.OrderID,
            TrackingId: shipment.TrackingID,
            CourierAwb: shipment.AWBNumber,
            Status:     shipment.Status,
            Label:      shipment.Label,
        }
    }

    return &pb.AccountShipmentsResponse{
        Success:   true,
        Shipments: pbShipments,
    }, nil
}