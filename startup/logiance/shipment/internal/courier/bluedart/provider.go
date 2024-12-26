package bluedart

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"

	"sync"
	"time"

	"github.com/Shridhar2104/logilo/shipment/internal/config"
	"github.com/Shridhar2104/logilo/shipment/internal/courier"
)

type Provider struct {
    config     *config.CourierConfig
    client     *courier.HTTPClient
    licenseKey string
    mu         sync.RWMutex
}

func NewProvider(cfg *config.CourierConfig) courier.CourierProvider {
    provider := &Provider{
        config: cfg,
        client: courier.NewHTTPClient(cfg.BaseURL, 30*time.Second),
        licenseKey: cfg.ApiKey,
    }

    // Set default headers for Bluedart
    provider.client.SetHeader("LicenseKey", cfg.ApiKey)
    provider.client.SetHeader("LoginID", cfg.ApiSecret)
    
    return provider
}

func (p *Provider) GetProviderInfo() *courier.ProviderInfo {
    return &courier.ProviderInfo{
        Code:        "BLUEDART",
        Name:        "Bluedart Express",
        Description: "Bluedart Express Shipping Services",
    }
}



func (p *Provider) CalculateRate(ctx context.Context, req *courier.RateRequest) (*courier.RateResponse, error) {
    // Create SOAP envelope
    soapEnvelope := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
        <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"
                      xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                      xmlns:xsd="http://www.w3.org/2001/XMLSchema">
            <soap:Body>
                <CalculateRate xmlns="http://tempuri.org/">
                    <Request>
                        <LoginID>%s</LoginID>
                        <LicenseKey>%s</LicenseKey>
                        <SourcePincode>%s</SourcePincode>
                        <DestinationPincode>%s</DestinationPincode>
                        <Weight>%.2f</Weight>
                        <Length>%.2f</Length>
                        <Width>%.2f</Width>
                        <Height>%.2f</Height>
                        <ProductCode>A</ProductCode>
                        <PaymentMode>%s</PaymentMode>
                        <CollectableAmount>%.2f</CollectableAmount>
                    </Request>
                </CalculateRate>
            </soap:Body>
        </soap:Envelope>`,
        p.config.ApiKey,
        p.config.ApiSecret,
        req.OriginPincode,
        req.DestinationPincode,
        req.Weight,
        req.Length,
        req.Width,
        req.Height,
        req.PaymentMode,
        req.CollectableAmount,
    )

    // Set SOAP action header
    headers := map[string]string{
        "Content-Type": "text/xml; charset=utf-8",
        "SOAPAction":   "http://tempuri.org/CalculateRate",
    }

    // Make the SOAP request using the HTTP client
    response, err := p.client.DoRaw(ctx, "POST", "/CalculateRate", headers, []byte(soapEnvelope))
    if err != nil {
        return nil, fmt.Errorf("SOAP request failed: %w", err)
    }

    // Parse SOAP response
    var soapResponse struct {
        XMLName xml.Name `xml:"Envelope"`
        Body    struct {
            CalculateRateResponse struct {
                Result struct {
                    BaseRate    float64 `xml:"BaseRate"`
                    FuelCharge  float64 `xml:"FuelSurcharge"`
                    CODCharge   float64 `xml:"CodCharge"`
                    OtherCharge float64 `xml:"OtherCharge"`
                    TotalAmount float64 `xml:"TotalAmount"`
                    ETADays     int     `xml:"EstimatedDays"`
                    Status      struct {
                        Code        string `xml:"Code"`
                        Description string `xml:"Description"`
                    } `xml:"Status"`
                } `xml:"CalculateRateResult"`
            } `xml:"CalculateRateResponse"`
        } `xml:"Body"`
    }

    if err := xml.NewDecoder(bytes.NewReader(response)).Decode(&soapResponse); err != nil {
        return nil, fmt.Errorf("failed to parse SOAP response: %w", err)
    }

    result := soapResponse.Body.CalculateRateResponse.Result
    
    // Check if the calculation was successful
    if result.Status.Code != "200" {
        return nil, fmt.Errorf("rate calculation failed: %s", result.Status.Description)
    }

    return &courier.RateResponse{
        BaseCharge:     result.BaseRate,
        FuelSurcharge:  result.FuelCharge,
        CODCharge:      result.CODCharge,
        HandlingCharge: result.OtherCharge,
        TotalCharge:    result.TotalAmount,
        ExpectedDays:   result.ETADays,
    }, nil
}

// HTTPClient interface extension (in courier/http_client.go)
type HTTPClient interface {
    Do(ctx context.Context, method, path string, request, response interface{}) error
    // Add DoRaw method for raw HTTP requests
    DoRaw(ctx context.Context, method, path string, headers map[string]string, body []byte) ([]byte, error)
}
func (p *Provider) CreateShipment(ctx context.Context, req *courier.ShipmentRequest) (*courier.ShipmentResponse, error) {
    // Bluedart specific shipment request format
    shipmentReq := map[string]interface{}{
        "profile": map[string]interface{}{
            "api_type": "S",
            "area": req.Consignee.City,
            "customerCode": p.config.ApiSecret,
            "licenseKey": p.licenseKey,
        },
        "services": map[string]interface{}{
            "product_code": getProductCode(req.PaymentType),
            "order_number": req.OrderNumber,
            "weight": req.PackageWeight,
            "dimensions": map[string]interface{}{
                "length": req.PackageLength,
                "width": req.PackageBreadth,
                "height": req.PackageHeight,
            },
            "payment_type": getPaymentType(req.PaymentType),
            "collect_amount": getCODAmount(req.PaymentType, req.OrderAmount),
        },
        "consignee": map[string]interface{}{
            "name": req.Consignee.Name,
            "company_name": req.Consignee.CompanyName,
            "address1": req.Consignee.AddressLine1,
            "address2": req.Consignee.AddressLine2,
            "pincode": req.Consignee.Pincode,
            "city": req.Consignee.City,
            "state": req.Consignee.State,
            "phone1": req.Consignee.Phone,
            "email": req.Consignee.Email,
        },
        "shipper": map[string]interface{}{
            "name": req.Pickup.Name,
            "company_name": req.Pickup.CompanyName,
            "address1": req.Pickup.AddressLine1,
            "address2": req.Pickup.AddressLine2,
            "pincode": req.Pickup.Pincode,
            "city": req.Pickup.City,
            "state": req.Pickup.State,
            "phone1": req.Pickup.Phone,
            "email": req.Pickup.Email,
        },
    }

    var response struct {
        Status struct {
            Code    int    `json:"code"`
            Message string `json:"message"`
        } `json:"status"`
        AWB struct {
            AWBNo     string `json:"awbno"`
            Status    string `json:"status"`
            Reference string `json:"reference"`
            LabelURL  string `json:"label_url"`
        } `json:"awb"`
    }

    err := p.client.Do(ctx, "POST", "/shipment/create", shipmentReq, &response)
    if err != nil {
        return nil, fmt.Errorf("create shipment failed: %w", err)
    }

    success := response.Status.Code == 200
    return &courier.ShipmentResponse{
        Success:     success,
        TrackingID:  response.AWB.AWBNo,
        AWBNumber:   response.AWB.AWBNo,
        Label:       response.AWB.LabelURL,
        Error:       response.Status.Message,
    }, nil
}

func (p *Provider) TrackShipment(ctx context.Context, trackingID string) ([]courier.TrackingEvent, error) {
    req := map[string]interface{}{
        "awbno": trackingID,
        "licenseKey": p.licenseKey,
    }

    var response struct {
        Status struct {
            Code    int    `json:"code"`
            Message string `json:"message"`
        } `json:"status"`
        Scans []struct {
            ScanCode    string `json:"scan_code"`
            ScanStatus  string `json:"scan_status"`
            ScanType    string `json:"scan_type"`
            ScanDate    string `json:"scan_date"`
            ScanTime    string `json:"scan_time"`
            Location    string `json:"location"`
            Remarks     string `json:"remarks"`
        } `json:"scans"`
    }

    err := p.client.Do(ctx, "POST", "/tracking/awb", req, &response)
    if err != nil {
        return nil, fmt.Errorf("tracking failed: %w", err)
    }

    if response.Status.Code != 200 {
        return nil, fmt.Errorf("tracking failed: %s", response.Status.Message)
    }

    var events []courier.TrackingEvent
    for _, scan := range response.Scans {
        activityTime := fmt.Sprintf("%s %s", scan.ScanDate, scan.ScanTime)
        events = append(events, courier.TrackingEvent{
            Status:       getMappedStatus(scan.ScanCode),
            StatusCode:   scan.ScanCode,
            Location:     scan.Location,
            ActivityTime: activityTime,
            Description: scan.Remarks,
        })
    }

    return events, nil
}

func (p *Provider) CheckServiceability(ctx context.Context, originPin, destinationPin string, weight float64) (bool, error) {
    req := map[string]interface{}{
        "profile": map[string]interface{}{
            "licenseKey": p.licenseKey,
        },
        "pincode_details": map[string]interface{}{
            "origin_pin": originPin,
            "destination_pin": destinationPin,
            "weight": weight,
        },
    }

    var response struct {
        Status struct {
            Code    int    `json:"code"`
            Message string `json:"message"`
        } `json:"status"`
        Serviceable bool `json:"serviceable"`
        Details struct {
            TransitDays int     `json:"transit_days"`
            BaseRate    float64 `json:"base_rate"`
        } `json:"details"`
    }

    err := p.client.Do(ctx, "POST", "/pincode/serviceability", req, &response)
    if err != nil {
        return false, fmt.Errorf("serviceability check failed: %w", err)
    }

    if response.Status.Code != 200 {
        return false, fmt.Errorf("serviceability check failed: %s", response.Status.Message)
    }

    return response.Serviceable, nil
}

func (p *Provider) CancelShipment(ctx context.Context, trackingID string) error {
    req := map[string]interface{}{
        "awbno": trackingID,
        "licenseKey": p.licenseKey,
        "reason": "Cancelled by customer",
    }

    var response struct {
        Status struct {
            Code    int    `json:"code"`
            Message string `json:"message"`
        } `json:"status"`
        Cancelled bool `json:"cancelled"`
    }

    err := p.client.Do(ctx, "POST", "/shipment/cancel", req, &response)
    if err != nil {
        return fmt.Errorf("cancel shipment failed: %w", err)
    }

    if response.Status.Code != 200 || !response.Cancelled {
        return fmt.Errorf("cancel shipment failed: %s", response.Status.Message)
    }

    return nil
}

// Helper functions

func getProductCode(paymentType string) string {
    switch paymentType {
    case "COD":
        return "A2B_COD"
    case "PREPAID":
        return "A2B_PPD"
    default:
        return "A2B_PPD"
    }
}

func getPaymentType(paymentType string) string {
    switch paymentType {
    case "COD":
        return "COD"
    case "PREPAID":
        return "PREPAID"
    default:
        return "PREPAID"
    }
}

func getCODAmount(paymentType string, orderAmount float64) float64 {
    if paymentType == "COD" {
        return orderAmount
    }
    return 0
}

func getMappedStatus(statusCode string) string {
    statusMap := map[string]string{
        "PDUP": "PICKUP_DONE",
        "PDNG": "PENDING_PICKUP",
        "INTC": "IN_TRANSIT",
        "DLVD": "DELIVERED",
        "CRTA": "CREATED",
        "CNCL": "CANCELLED",
        "RDEL": "RTO_DELIVERED",
        "RHLD": "RTO_HOLD",
        "DLVG": "OUT_FOR_DELIVERY",
        "EXC1": "EXCEPTION",
        "LOST": "LOST",
    }

    if status, exists := statusMap[statusCode]; exists {
        return status
    }
    return statusCode
}

// Calculate volumetric weight
func calculateVolumetricWeight(length, breadth, height float64) float64 {
    return (length * breadth * height) / 5000
}