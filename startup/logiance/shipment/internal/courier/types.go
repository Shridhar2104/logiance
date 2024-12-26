package courier

import "context"

// Provider interface defines methods that must be implemented by all courier providers
type CourierProvider interface {
    // CalculateRate calculates shipping rates
    CalculateRate(ctx context.Context, req *RateRequest) (*RateResponse, error)
    
    // CreateShipment creates a new shipment
    CreateShipment(ctx context.Context, req *ShipmentRequest) (*ShipmentResponse, error)
    
    // TrackShipment gets tracking information for a shipment
    TrackShipment(ctx context.Context, trackingID string) ([]TrackingEvent, error)
    
    // CheckServiceability checks if delivery is possible
    CheckServiceability(ctx context.Context, originPin, destinationPin string, weight float64) (bool, error)
    
    // GetProviderInfo returns information about the courier provider
    GetProviderInfo() *ProviderInfo
}
type CourierType string
const (
    CourierXpressbees CourierType = "XPRESSBEES"
    CourierDelhivery  CourierType = "DELHIVERY"
    CourierBluedart   CourierType = "BLUEDART"
    // Add other courier types as needed
)


// RateRequest contains parameters for rate calculation
type RateRequest struct {
    OriginPincode      string  `json:"origin_pincode"`
    DestinationPincode string  `json:"destination_pincode"`
    Weight             float64 `json:"weight"`
    Length             float64 `json:"length"`
    Width              float64 `json:"width"`
    Height             float64 `json:"height"`
    PaymentMode        string  `json:"payment_mode"` // COD/Prepaid
    CollectableAmount  float64 `json:"collectable_amount"`
}

// RateResponse contains calculated shipping rates
type RateResponse struct {
    BaseCharge     float64 `json:"base_charge"`
    FuelSurcharge  float64 `json:"fuel_surcharge"`
    CODCharge      float64 `json:"cod_charge"`
    HandlingCharge float64 `json:"handling_charge"`
    TotalCharge    float64 `json:"total_charge"`
    ExpectedDays   int     `json:"expected_days"`
}
// ShipmentRequest contains details for creating a new shipment
type ShipmentRequest struct {
    OrderNumber     string         `json:"order_number"`
    PaymentType     string         `json:"payment_type"` // COD/Prepaid
    PackageWeight   float64        `json:"package_weight"`
    PackageLength   float64        `json:"package_length"`
    PackageBreadth  float64        `json:"package_breadth"`
    PackageHeight   float64        `json:"package_height"`
    OrderAmount     float64        `json:"order_amount"`
    CollectableAmount   float64    `json:collectable_amount`
    ShippingCharges float64        `json:"shipping_charges"`   // Added shipping charges
    Discount        float64        `json:"discount"`           // Added discount
    CODCharges      float64        `json:"cod_charges"`        // Added COD charges
    Consignee       Address        `json:"consignee"`
    Pickup          Address        `json:"pickup"`
    Items           []OrderItem    `json:"items,omitempty"`
    AutoPickup      bool           `json:"auto_pickup,omitempty"`
    ReturnInfo      *ReturnInfo    `json:"return_info,omitempty"`
	RTO RTODetails `json:"rto_details"`
}

// RTODetails represents RTO (Return to Origin) specific information
type RTODetails struct {
    IsDifferent    bool    `json:"is_rto_different"`
    WarehouseName  string  `json:"warehouse_name,omitempty"`
    Name           string  `json:"name,omitempty"`
    AddressLine1   string  `json:"address_1,omitempty"`
    AddressLine2   string  `json:"address_2,omitempty"`
    City           string  `json:"city,omitempty"`
    State          string  `json:"state,omitempty"`
    Pincode        string  `json:"pincode,omitempty"`
    Phone          string  `json:"phone,omitempty"`
    GSTIN          string  `json:"gstin,omitempty"`
}

// ShipmentResponse contains the response from shipment creation
type ShipmentResponse struct {
    Success     bool   `json:"success"`
	OrderID  	string `json:"order_id"`
	ShipmentID	string `json:"shipment_id"`
    TrackingID  string `json:"tracking_id,omitempty"`
    AWBNumber   string `json:"awb_number,omitempty"`
    CourierName string `json:"courier_name,omitempty"`
    Label       string `json:"label,omitempty"`
    Error       string `json:"error,omitempty"`
}

// Address contains shipping address details
type Address struct {
    Name        string `json:"name"`
    CompanyName string `json:"company_name,omitempty"`
    Phone       string `json:"phone"`
    Email       string `json:"email,omitempty"`
    AddressLine1 string `json:"address_line1"`
    AddressLine2 string `json:"address_line2,omitempty"`
    Landmark    string `json:"landmark,omitempty"`
    City        string `json:"city"`
    State       string `json:"state"`
    Country     string `json:"country"`
    Pincode     string `json:"pincode"`
    GSTIN       string `json:"gstin,omitempty"`
}

// OrderItem contains details of items in the shipment
type OrderItem struct {
    SKU           string  `json:"sku"`
    Name          string  `json:"name"`
    Quantity      int     `json:"quantity"`
    Price         float64 `json:"price"`
    HSNCode       string  `json:"hsn_code,omitempty"`
    Category      string  `json:"category,omitempty"`
    ActualWeight  float64 `json:"actual_weight,omitempty"`
}

// ReturnInfo contains details for return shipments
type ReturnInfo struct {
    Address         Address `json:"address"`
    AWBNumber       string  `json:"awb_number,omitempty"`
    ReturnReason    string  `json:"return_reason,omitempty"`
    ReturnComment   string  `json:"return_comment,omitempty"`
}

// TrackingEvent contains shipment tracking information
type TrackingEvent struct {
    Status       string `json:"status"`
    StatusCode   string `json:"status_code"`
    Location     string `json:"location"`
    ActivityTime string `json:"activity_time"`
    Description  string `json:"description"`
}

// ProviderInfo contains information about the courier provider
type ProviderInfo struct {
    Code        string `json:"code"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

// ShipmentStatus represents the current status of a shipment
type ShipmentStatus string

const (
    StatusPending    ShipmentStatus = "PENDING"
    StatusPickup     ShipmentStatus = "PICKUP"
    StatusInTransit  ShipmentStatus = "IN_TRANSIT"
    StatusOutForDelivery ShipmentStatus = "OUT_FOR_DELIVERY"
    StatusDelivered  ShipmentStatus = "DELIVERED"
    StatusException  ShipmentStatus = "EXCEPTION"
    StatusCancelled  ShipmentStatus = "CANCELLED"
    StatusRTO        ShipmentStatus = "RTO"
)

// NDRInfo contains information about NDR (Non-Delivery Report)
type NDRInfo struct {
    AWBNumber     string `json:"awb_number"`
    Attempt       int    `json:"attempt"`
    Reason        string `json:"reason"`
    SubReason     string `json:"sub_reason,omitempty"`
    Comment       string `json:"comment,omitempty"`
    AttemptedAt   string `json:"attempted_at"`
    NextAttemptAt string `json:"next_attempt_at,omitempty"`
}

// NDRAction represents actions that can be taken for NDR
type NDRDetails struct{
	AWBNumber     string `json:"awb_number"`
	EventDate       string    `json:"event_date"`
	CourierRemarks	string `json:"courier_remarks"`
	TotalAttempts	int `json:"total_attempts"`

}
type NDRAction struct {
    AWBNumber    string            `json:"awb_number"`
    ActionType   string            `json:"action_type"` // retry/reschedule/cancel
    ActionData   map[string]string `json:"action_data,omitempty"`
}

// ServiceabilityRequest contains parameters for serviceability check
type ServiceabilityRequest struct {
    OriginPincode      string  `json:"origin_pincode"`
    DestinationPincode string  `json:"destination_pincode"`
    Weight             float64 `json:"weight"`
    Length             float64 `json:"length,omitempty"`
    Width              float64 `json:"width,omitempty"`
    Height             float64 `json:"height,omitempty"`
    PaymentMode        string  `json:"payment_mode,omitempty"`
}

// ServiceabilityResponse contains serviceability check results
type ServiceabilityResponse struct {
    Serviceable     bool    `json:"serviceable"`
    EstimatedDays   int     `json:"estimated_days,omitempty"`
    EstimatedCharge float64 `json:"estimated_charge,omitempty"`
    Error           string  `json:"error,omitempty"`
}

