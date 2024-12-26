// internal/database/models.go
package database

import (
    "time"
    "gorm.io/gorm"
)

// ShipmentTracking represents the tracking information for a shipment
type ShipmentTracking struct {
    ID            uint      `gorm:"primaryKey"`
    AccountID     string    `gorm:"index;not null"`
    OrderID       string    `gorm:"uniqueIndex;not null"`
    TrackingID    string    `gorm:"uniqueIndex;not null"`
    AWBNumber     string    `gorm:"uniqueIndex;not null"`
    CourierCode   string    `gorm:"index;not null"`
    Status        string    `gorm:"index"`
    Label         string
    CreatedAt     time.Time
    UpdatedAt     time.Time
    DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// ShipmentEvent represents individual tracking events for a shipment
type ShipmentEvent struct {
    ID                 uint      `gorm:"primaryKey"`
    ShipmentTrackingID uint      `gorm:"index;not null"`
    Status            string    `gorm:"index"`
    Location          string
    Timestamp         time.Time `gorm:"index"`
    Description       string
    CreatedAt         time.Time
}