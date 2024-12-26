// internal/database/operations.go
package database

import (
    "context"
    "fmt"
    "gorm.io/gorm"
)

type ShipmentDB struct {
    db *gorm.DB
}

func NewShipmentDB(db *gorm.DB) *ShipmentDB {
    return &ShipmentDB{db: db}
}

// Database operations
// SaveShipmentTracking saves a new shipment tracking record
func (s *ShipmentDB) SaveShipmentTracking(ctx context.Context, tracking *ShipmentTracking) error {
    result := s.db.WithContext(ctx).Create(tracking)
    if result.Error != nil {
        return fmt.Errorf("failed to save shipment tracking: %w", result.Error)
    }
    return nil
}

// GetShipmentByTrackingID retrieves shipment tracking info by tracking ID
func (s *ShipmentDB) GetShipmentByTrackingID(ctx context.Context, trackingID string) (*ShipmentTracking, error) {
    var tracking ShipmentTracking
    result := s.db.WithContext(ctx).Where("tracking_id = ?", trackingID).First(&tracking)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get shipment tracking: %w", result.Error)
    }
    return &tracking, nil
}

// GetShipmentByOrderID retrieves shipment tracking info by order ID
func (s *ShipmentDB) GetShipmentByOrderID(ctx context.Context, orderID string) (*ShipmentTracking, error) {
    var tracking ShipmentTracking
    result := s.db.WithContext(ctx).Where("order_id = ?", orderID).First(&tracking)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get shipment tracking: %w", result.Error)
    }
    return &tracking, nil
}

// GetShipmentsByAccountID retrieves all shipments for an account
func (s *ShipmentDB) GetShipmentsByAccountID(ctx context.Context, accountID string, limit, offset int) ([]ShipmentTracking, error) {
    var trackings []ShipmentTracking
    result := s.db.WithContext(ctx).
        Where("account_id = ?", accountID).
        Order("created_at DESC").
        Limit(limit).
        Offset(offset).
        Find(&trackings)
    
    if result.Error != nil {
        return nil, fmt.Errorf("failed to get shipments: %w", result.Error)
    }
    return trackings, nil
}

// UpdateShipmentStatus updates the status of a shipment
func (s *ShipmentDB) UpdateShipmentStatus(ctx context.Context, trackingID, status string) error {
    result := s.db.WithContext(ctx).
        Model(&ShipmentTracking{}).
        Where("tracking_id = ?", trackingID).
        Update("status", status)
    
    if result.Error != nil {
        return fmt.Errorf("failed to update shipment status: %w", result.Error)
    }
    return nil
}

// SaveShipmentEvent saves a new tracking event for a shipment
func (s *ShipmentDB) SaveShipmentEvent(ctx context.Context, event *ShipmentEvent) error {
    result := s.db.WithContext(ctx).Create(event)
    if result.Error != nil {
        return fmt.Errorf("failed to save shipment event: %w", result.Error)
    }
    return nil
}

// GetShipmentEvents retrieves all tracking events for a shipment
func (s *ShipmentDB) GetShipmentEvents(ctx context.Context, shipmentTrackingID uint) ([]ShipmentEvent, error) {
    var events []ShipmentEvent
    result := s.db.WithContext(ctx).
        Where("shipment_tracking_id = ?", shipmentTrackingID).
        Order("timestamp DESC").
        Find(&events)
    
    if result.Error != nil {
        return nil, fmt.Errorf("failed to get shipment events: %w", result.Error)
    }
    return events, nil
}