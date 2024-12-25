package xpressbees

import (
	"context"
	"strings"

	"fmt"
	"log"
	"strconv"

	"sync"
	"time"

	"github.com/Shridhar2104/logilo/shipment/internal/config"
	"github.com/Shridhar2104/logilo/shipment/internal/courier"
)
type Provider struct {
    client     *courier.HTTPClient
    config     *config.CourierConfig
    authToken  string
    tokenExp   time.Time
    mu         sync.RWMutex
}

func NewProvider(cfg *config.CourierConfig) (courier.CourierProvider, error) {
    if cfg == nil {
        return nil, fmt.Errorf("courier config cannot be nil")
    }
    
    if cfg.BaseURL == "" {
        return nil, fmt.Errorf("base URL is required")
    }
    
    if cfg.ApiKey == "" {
        return nil, fmt.Errorf("API key (email) is required")
    }
    
    if cfg.ApiSecret == "" {
        return nil, fmt.Errorf("API secret (password) is required")
    }

    // Create HTTP client with base URL and timeout
    client := courier.NewHTTPClient(cfg.BaseURL, 30*time.Second)
    
    provider := &Provider{
        client: client,
        config: cfg,
    }

    // Initial authentication with context
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := provider.authenticate(ctx); err != nil {
        return nil, fmt.Errorf("failed to authenticate with Xpressbees: %w", err)
    }

    return provider, nil
}

func (p *Provider) GetProviderInfo() *courier.ProviderInfo {
    return &courier.ProviderInfo{
        Code:        "XPRESSBEES",
        Name:        "Xpressbees",
        Description: "Xpressbees Shipping Services",
    }
}
func (p *Provider) authenticate(ctx context.Context) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.authToken != "" && time.Now().Before(p.tokenExp) {
        return nil
    }

    log.Printf("Attempting Xpressbees authentication with email: %s", p.config.ApiKey)

    // Use proper response structure
    var authResp struct {
        Status  bool   `json:"status"`
        Data    string `json:"data"`    // JWT token comes directly as string
        Message string `json:"message"`
        Error   string `json:"error"`
    }

    // Make the authentication request
    err := p.client.Do(ctx, "POST", "/users/login", 
        map[string]string{
            "email": p.config.ApiKey,
            "password": p.config.ApiSecret,
        }, 
        &authResp)

    if err != nil {
        return fmt.Errorf("authentication request failed: %w", err)
    }

    if !authResp.Status {
        if authResp.Message != "" {
            return fmt.Errorf("authentication failed: %s", authResp.Message)
        }
        if authResp.Error != "" {
            return fmt.Errorf("authentication failed: %s", authResp.Error)
        }
        return fmt.Errorf("authentication failed with unknown error")
    }

    if authResp.Data == "" {
        return fmt.Errorf("authentication failed: no token received")
    }

    p.authToken = authResp.Data
    p.tokenExp = time.Now().Add(23 * time.Hour)
    p.client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", p.authToken))

    log.Printf("Xpressbees authentication successful")
    return nil
}

// parseExpectedDeliveryDays calculates the number of days until expected delivery
func parseExpectedDeliveryDays(deliveryDate string) int {
    if deliveryDate == "" {
        return 0
    }

    // Parse the expected delivery date string
    // Assuming the date format is "2006-01-02"
    expectedDate, err := time.Parse("2006-01-02", deliveryDate)
    if err != nil {
        return 0
    }

    // Calculate days between now and expected delivery date
    days := int(expectedDate.Sub(time.Now()).Hours() / 24)
    if days < 0 {
        return 0
    }

    return days
}
func (p *Provider) CalculateRate(ctx context.Context, req *courier.RateRequest) (*courier.RateResponse, error) {
    if err := p.authenticate(ctx); err != nil {
        return nil, err
    }

    // Request structure following Xpressbees API format
    rateReq := map[string]interface{}{
        "origin":          req.OriginPincode,      // Changed from origin_pin
        "destination":     req.DestinationPincode, // Changed from destination_pin
        "payment_type":    strings.ToLower(req.PaymentMode), // Must be "cod" or "prepaid"
        "order_amount":    req.CollectableAmount,  // Required for COD orders
        "weight":          req.Weight,             // Weight in grams
        "length":          req.Length,             // Length in cm
        "breadth":         req.Width,              // Width/Breadth in cm
        "height":          req.Height,             // Height in cm
    }

    var response struct {
        Status bool `json:"status"`
        Data   []struct {
            ID              string  `json:"id"`
            Name            string  `json:"name"`
            FreightCharges  float64 `json:"freight_charges"`
            CODCharges      float64 `json:"cod_charges"`
            TotalCharges    float64 `json:"total_charges"`
            MinWeight       int     `json:"min_weight"`
            ChargeableWeight int    `json:"chargeable_weight"`
        } `json:"data"`
        Message string `json:"message"`
    }

    // Changed endpoint to match API documentation
    err := p.client.Do(ctx, "POST", "/courier/serviceability", rateReq, &response)
    if err != nil {
        return nil, fmt.Errorf("rate calculation failed: %w", err)
    }

    // Check if response is successful
    if !response.Status {
        if response.Message != "" {
            return nil, fmt.Errorf("rate calculation failed: %s", response.Message)
        }
        return nil, fmt.Errorf("rate calculation failed with unknown error")
    }

    // If no service available
    if len(response.Data) == 0 {
        return nil, fmt.Errorf("no service available for given parameters")
    }

    // Use the first available rate
    rate := response.Data[0]
    return &courier.RateResponse{
        BaseCharge:     rate.FreightCharges,
        CODCharge:      rate.CODCharges,
        TotalCharge:    rate.TotalCharges,
        ExpectedDays:   2, // Default value as API doesn't provide this
    }, nil
}
func (p *Provider) CreateShipment(ctx context.Context, req *courier.ShipmentRequest) (*courier.ShipmentResponse, error) {
    if err := p.authenticate(ctx); err != nil {
        return nil, err
    }

    log.Printf("Creating shipment with request: %+v", req)

    // Build shipment request
    xbReq := map[string]interface{}{
        "order_number":        req.OrderNumber,
        "payment_type":        strings.ToLower(req.PaymentType),
        "package_weight":      strconv.Itoa(int(req.PackageWeight)),
        "package_length":      strconv.Itoa(int(req.PackageLength)),
        "package_breadth":     strconv.Itoa(int(req.PackageBreadth)),
        "package_height":      strconv.Itoa(int(req.PackageHeight)),
        "shipping_charges":    fmt.Sprintf("%.2f", req.ShippingCharges),
        "cod_charges":         fmt.Sprintf("%.2f", req.CODCharges),
        "discount":            "0",
        "order_amount":        fmt.Sprintf("%.2f", req.OrderAmount),
        "collectable_amount":  fmt.Sprintf("%.2f", req.CollectableAmount),
        "request_auto_pickup": "no",

        // Consignee details
        "consignee": map[string]interface{}{
            "name":         req.Consignee.Name,
            "company_name": req.Consignee.CompanyName,
            "address":      req.Consignee.AddressLine1,
            "address_2":    req.Consignee.AddressLine2,
            "city":         req.Consignee.City,
            "state":        req.Consignee.State,
            "pincode":      req.Consignee.Pincode,
            "phone":        req.Consignee.Phone,
            "email":        req.Consignee.Email,
        },

        // Pickup details
        "pickup": map[string]interface{}{
            "warehouse_name": req.Pickup.CompanyName,
            "name":          req.Pickup.Name,
            "address":       req.Pickup.AddressLine1,
            "address_2":     req.Pickup.AddressLine2,
            "city":          req.Pickup.City,
            "state":         req.Pickup.State,
            "pincode":       req.Pickup.Pincode,
            "phone":         req.Pickup.Phone,
            "gst_number":    req.Pickup.GSTIN,
        },

        // Order items
        "order_items": []map[string]interface{}{{
            "name":  req.Items[0].Name,
            "qty":   strconv.Itoa(req.Items[0].Quantity),
            "price": fmt.Sprintf("%.2f", req.Items[0].Price),
            "sku":   req.Items[0].SKU,
        }},
    }

    log.Printf("Sending request to Xpressbees: %+v", xbReq)
    
    // ... rest of the function remains the same
    var response struct {
        Status  bool   `json:"status"`
        Message string `json:"message"`
        Data    struct {
            OrderID     string `json:"order_id"`
            ShipmentID  string `json:"shipment_id"`
            AWBNumber   string `json:"awb_number"`
            CourierID   string `json:"courier_id"`
            CourierName string `json:"courier_name"`
            Status      string `json:"status"`
            Label       string `json:"label"`
        } `json:"data"`
    }

    err := p.client.Do(ctx, "POST", "/shipments2", xbReq, &response)
    if err != nil {
        log.Printf("API error: %v", err)
        return nil, fmt.Errorf("create shipment failed: %w", err)
    }

    log.Printf("API Response: %+v", response)

    return &courier.ShipmentResponse{
        Success:     response.Status,
        OrderID:     response.Data.OrderID,
        ShipmentID:  response.Data.ShipmentID,
        TrackingID:  response.Data.AWBNumber,
        AWBNumber:   response.Data.AWBNumber,
        CourierName: response.Data.CourierName,
        Label:       response.Data.Label,
        Error:       response.Message,
    }, nil
}

func (p *Provider) TrackShipment(ctx context.Context, awbNumber string) ([]courier.TrackingEvent, error) {
    if err := p.authenticate(ctx); err != nil {
        return nil, err
    }

    var response struct {
        Status bool `json:"status"`
        Data   struct {
            History []struct {
                StatusCode string `json:"status_code"`
                Location   string `json:"location"`
                EventTime  string `json:"event_time"`
                Message    string `json:"message"`
            } `json:"history"`
        } `json:"data"`
    }

    err := p.client.Do(ctx, "GET", fmt.Sprintf("/shipments2/track/%s", awbNumber), nil, &response)
    if err != nil {
        return nil, fmt.Errorf("tracking failed: %w", err)
    }

    details := make([]courier.TrackingEvent, len(response.Data.History))
    for i, event := range response.Data.History {
        details[i] = courier.TrackingEvent{
            StatusCode:   event.StatusCode,
            Location:     event.Location,
            ActivityTime: event.EventTime,
            Description: event.Message,
        }
    }

    return details, nil
}

func (p *Provider) CancelShipment(ctx context.Context, awbNumber string) error {
    if err := p.authenticate(ctx); err != nil {
        return err
    }

    cancelReq := map[string]string{
        "awb": awbNumber,
    }

    var response struct {
        Status  bool   `json:"status"`
        Message string `json:"message"`
    }

    err := p.client.Do(ctx, "POST", "/shipments2/cancel", cancelReq, &response)
    if err != nil {
        return fmt.Errorf("cancel shipment failed: %w", err)
    }

    if !response.Status {
        return fmt.Errorf("cancel shipment failed: %s", response.Message)
    }

    return nil
}

func (p *Provider) GetNDRList(ctx context.Context, page int, limit int) ([]courier.NDRDetails, error) {
    if err := p.authenticate(ctx); err != nil {
        return nil, err
    }

    var response struct {
        Status bool `json:"status"`
        Data   []struct {
            AWBNumber      string `json:"awb_number"`
            EventDate      string `json:"event_date"`
            CourierRemarks string `json:"courier_remarks"`
            TotalAttempts  string `json:"total_attempts"`
        } `json:"data"`
    }

    queryParams := fmt.Sprintf("?per_page=%d&page=%d", limit, page)
    err := p.client.Do(ctx, "GET", "/ndr"+queryParams, nil, &response)
    if err != nil {
        return nil, fmt.Errorf("get NDR list failed: %w", err)
    }

    ndrList := make([]courier.NDRDetails, len(response.Data))
    for i, ndr := range response.Data {
        attempts, _ := strconv.Atoi(ndr.TotalAttempts)
        ndrList[i] = courier.NDRDetails{
            AWBNumber:      ndr.AWBNumber,
            EventDate:      ndr.EventDate,
            CourierRemarks: ndr.CourierRemarks,
            TotalAttempts:  attempts,
        }
    }

    return ndrList, nil
}

func (p *Provider) UpdateNDR(ctx context.Context, actions []courier.NDRAction) error {
    if err := p.authenticate(ctx); err != nil {
        return err
    }

    var response struct {
        Status  bool   `json:"status"`
        Message string `json:"message"`
    }

    err := p.client.Do(ctx, "POST", "/ndr/create", actions, &response)
    if err != nil {
        return fmt.Errorf("update NDR failed: %w", err)
    }

    if !response.Status {
        return fmt.Errorf("update NDR failed: %s", response.Message)
    }

    return nil
}

func (p *Provider) CheckServiceability(ctx context.Context, originPin, destinationPin string, weight float64) (bool, error) {
    if err := p.authenticate(ctx); err != nil {
        return false, err
    }

    serviceReq := map[string]interface{}{
        "origin":      originPin,
        "destination": destinationPin,
        "weight":      weight,
    }

    var response struct {
        Status bool `json:"status"`
        Data   []struct {
            ID string `json:"id"`
        } `json:"data"`
    }

    err := p.client.Do(ctx, "POST", "/courier/serviceability", serviceReq, &response)
    if err != nil {
        return false, fmt.Errorf("serviceability check failed: %w", err)
    }

    return response.Status && len(response.Data) > 0, nil
}

func (p *Provider) GetCourierType() courier.CourierType {
    return courier.CourierXpressbees
}