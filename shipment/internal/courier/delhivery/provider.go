// internal/courier/delhivery/provider.go
package delhivery

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"
	"net/url"

	"time"

	"github.com/Shridhar2104/logilo/shipment/internal/config"
	"github.com/Shridhar2104/logilo/shipment/internal/courier"
	pb "github.com/Shridhar2104/logilo/shipment/proto"
)

type DelhiveryProvider struct {
    cfg        *config.Config
    httpClient *http.Client
}

func NewProvider(cfg *config.Config) courier.Provider {
    return &DelhiveryProvider{
        cfg: cfg,
        httpClient: &http.Client{
            Timeout: time.Second * 10,
        },
    }
}



type DelhiveryResponse struct {
    Success      bool    `json:"success"`
    TotalAmount  float64 `json:"total_amount"`
    GrossAmount  float64 `json:"gross_amount"`
    TaxAmount    float64 `json:"tax_amount"`
    ErrorMessage string  `json:"error,omitempty"`
}

func (d *DelhiveryProvider) CalculateRate(ctx context.Context, req *pb.RateRequest) (*pb.CourierRate, error) {
    // Build query parameters according to Delhivery API requirements
    params := url.Values{}
    
    // Add mandatory parameters
    params.Add("md", determineMode(req.PaymentMode))    // Use payment_mode to determine shipping mode
    params.Add("cgm", fmt.Sprintf("%d", req.Weight))    // Weight in grams
    params.Add("o_pin", fmt.Sprintf("%d", req.OriginPincode))
    params.Add("d_pin", fmt.Sprintf("%d", req.DestinationPincode))
    params.Add("ss", "Delivered")                       // Default to "Delivered" status

    // Build the complete URL
    baseURL := d.cfg.CourierConfig["DELHIVERY"].BaseURL
    if baseURL == "" {
        baseURL = "https://staging-express.delhivery.com" // Default to staging if not configured
    }
    
    apiURL := fmt.Sprintf("%s/api/kinko/v1/invoice/charges/.json?%s",
        baseURL,
        params.Encode())

    // Create new request with context
    httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    // Add authorization header if configured
    if apiKey := d.cfg.CourierConfig["DELHIVERY"].ApiKey; apiKey != "" {
        httpReq.Header.Add("Authorization", apiKey)
    }

    // Make the request
    resp, err := d.httpClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()

    // Handle non-200 responses
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
    }

    // Parse the response
    var delhiveryResp DelhiveryResponse
    if err := json.NewDecoder(resp.Body).Decode(&delhiveryResp); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    // Check for API-level errors
    if !delhiveryResp.Success {
        return nil, fmt.Errorf("API returned error: %s", delhiveryResp.ErrorMessage)
    }

    // Determine service type based on payment mode
    serviceType := determineServiceType(req.PaymentMode)

    // Map to CourierRate response
    return &pb.CourierRate{
        CourierCode: "DELHIVERY",
        CourierName: "Delhivery",
        ServiceType: serviceType,
        RateDetails: &pb.RateData{
            TotalAmount:    delhiveryResp.TotalAmount,
            GrossAmount:    delhiveryResp.GrossAmount,
            TaxAmount:      delhiveryResp.TaxAmount,
            CodCharges:     calculateCODCharges(req.PaymentMode),
            FuelSurcharge:  0, // Set if available in response
        },
        EstimatedDays: estimateDeliveryDays(serviceType),
    }, nil
}

// Helper functions for parameter mapping

func determineMode(paymentMode string) string {
    // For Delhivery, we'll use Surface by default
    // You might want to adjust this logic based on your business rules
    return "S"
}

func determineServiceType(paymentMode string) string {
    switch paymentMode {
    case "COD":
        return "Surface-COD"
    case "Prepaid":
        return "Surface-Prepaid"
    default:
        return "Surface"
    }
}

func calculateCODCharges(paymentMode string) float64 {
    if paymentMode == "COD" {
        return 50.0 // Example COD charge, adjust based on actual Delhivery rates
    }
    return 0.0
}

func estimateDeliveryDays(serviceType string) int32 {
    // You might want to adjust these estimates based on your experience with Delhivery
    return 4 // Default to 4 days for Surface
}
func (d *DelhiveryProvider) IsAvailable(ctx context.Context, originPin, destPin int32) (bool, error) {
    // Implementation for checking serviceability
    return true, nil
}

func (d *DelhiveryProvider) GetProviderInfo() *pb.CourierInfo {
    return &pb.CourierInfo{
        CourierCode: "DELHIVERY",
        CourierName: "Delhivery",
        SupportedServices: []string{"express", "surface"},
    }
}