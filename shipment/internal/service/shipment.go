package service

import (
    "context"
    "fmt"
    "time"
    
    "github.com/Shridhar2104/logilo/shipment/internal/database"
)

type ShipmentTrackingService struct {
    db *database.ShipmentDB
}

func NewShipmentTrackingService(db *database.ShipmentDB) *ShipmentTrackingService {
    return &ShipmentTrackingService{db: db}
}

// SaveNewShipment saves tracking information for a newly created shipment
func (s *ShipmentTrackingService) SaveNewShipment(ctx context.Context, accountID, orderID, trackingID, awbNumber, courierCode, label string) error {
    tracking := &database.ShipmentTracking{
        AccountID:   accountID,
        OrderID:     orderID,
        TrackingID:  trackingID,
        AWBNumber:   awbNumber,
        CourierCode: courierCode,
        Status:      "CREATED",
        Label:       label,
        CreatedAt:   time.Now(),
    }

    err := s.db.SaveShipmentTracking(ctx, tracking)
    if err != nil {
        return fmt.Errorf("failed to save new shipment: %w", err)
    }
    
    // Create initial tracking event
    event := &database.ShipmentEvent{
        ShipmentTrackingID: tracking.ID,
        Status:            "CREATED",
        Description:       "Shipment created successfully",
        Timestamp:         time.Now(),
    }
    
    err = s.db.SaveShipmentEvent(ctx, event)
    if err != nil {
        return fmt.Errorf("failed to save initial tracking event: %w", err)
    }
    
    return nil
}

// UpdateShipmentStatus updates the status and adds a new tracking event
func (s *ShipmentTrackingService) UpdateShipmentStatus(ctx context.Context, trackingID, status, location, description string) error {
    // Update main tracking status
    err := s.db.UpdateShipmentStatus(ctx, trackingID, status)
    if err != nil {
        return fmt.Errorf("failed to update shipment status: %w", err)
    }
    
    // Get shipment tracking record to get the ID
    tracking, err := s.db.GetShipmentByTrackingID(ctx, trackingID)
    if err != nil {
        return fmt.Errorf("failed to get shipment tracking: %w", err)
    }
    if tracking == nil {
        return fmt.Errorf("shipment not found for tracking ID: %s", trackingID)
    }
    
    // Create new tracking event
    event := &database.ShipmentEvent{
        ShipmentTrackingID: tracking.ID,
        Status:            status,
        Location:          location,
        Description:       description,
        Timestamp:         time.Now(),
    }
    
    err = s.db.SaveShipmentEvent(ctx, event)
    if err != nil {
        return fmt.Errorf("failed to save tracking event: %w", err)
    }
    
    return nil
}

// GetShipmentDetails retrieves shipment tracking details and events
func (s *ShipmentTrackingService) GetShipmentDetails(ctx context.Context, trackingID string) (*database.ShipmentTracking, []database.ShipmentEvent, error) {
    tracking, err := s.db.GetShipmentByTrackingID(ctx, trackingID)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to get shipment tracking: %w", err)
    }
    if tracking == nil {
        return nil, nil, nil
    }
    
    events, err := s.db.GetShipmentEvents(ctx, tracking.ID)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to get shipment events: %w", err)
    }
    
    return tracking, events, nil
}

// GetShipmentsByAccount retrieves all shipments for an account with pagination
func (s *ShipmentTrackingService) GetShipmentsByAccount(ctx context.Context, accountID string, page, pageSize int) ([]database.ShipmentTracking, error) {
    offset := (page - 1) * pageSize
    return s.db.GetShipmentsByAccountID(ctx, accountID, pageSize, offset)
}

// GetShipmentByOrder retrieves shipment tracking info for an order
func (s *ShipmentTrackingService) GetShipmentByOrder(ctx context.Context, orderID string) (*database.ShipmentTracking, error) {
    return s.db.GetShipmentByOrderID(ctx, orderID)
}