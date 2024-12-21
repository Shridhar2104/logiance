// shipment/internal/courier/delhivery/client.go
package delhivery

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "github.com/Shridhar2104/logilo/shipment/internal/config"
    "golang.org/x/time/rate"
)

type Client struct {
    httpClient  *http.Client
    config      *config.DelhiveryConfig
    rateLimiter *rate.Limiter
}

type ShippingChargeResponse struct {
    TotalAmount  float64 `json:"total_amount"`
    GrossAmount  float64 `json:"gross_amount"`
    Tax          float64 `json:"tax"`
    ErrorMessage string  `json:"error,omitempty"`
}

type ShippingChargeRequest struct {
    MD    string `json:"md"`     // Billing Mode (E/S)
    CGM   int    `json:"cgm"`    // Chargeable weight in grams
    OPin  int    `json:"o_pin"`  // Origin pincode
    DPin  int    `json:"d_pin"`  // Destination pincode
    SS    string `json:"ss"`     // Shipment status
}

func NewClient(cfg *config.DelhiveryConfig) *Client {
    // Create rate limiter: 40 requests per minute
    limiter := rate.NewLimiter(rate.Every(time.Minute/40), 40)
    
    return &Client{
        httpClient:  &http.Client{Timeout: 10 * time.Second},
        config:      cfg,
        rateLimiter: limiter,
    }
}

func (c *Client) CalculateShippingCharges(ctx context.Context, req *ShippingChargeRequest) (*ShippingChargeResponse, error) {
    // Wait for rate limiter
    err := c.rateLimiter.Wait(ctx)
    if err != nil {
        return nil, fmt.Errorf("rate limit error: %w", err)
    }

    // Construct URL
    baseURL := c.config.BaseURL
    if c.config.Environment == "staging" {
        baseURL = "https://staging-express.delhivery.com"
    }
    
    url := fmt.Sprintf("%s/api/kinko/v1/invoice/charges/.json?md=%s&cgm=%d&o_pin=%d&d_pin=%d&ss=%s",
        baseURL,
        req.MD,
        req.CGM,
        req.OPin,
        req.DPin,
        req.SS,
    )

    // Create request
    httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("creating request: %w", err)
    }

    // Add headers
    httpReq.Header.Set("Authorization", c.config.APIKey)
    httpReq.Header.Set("Content-Type", "application/json")

    // Make request
    resp, err := c.httpClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("making request: %w", err)
    }
    defer resp.Body.Close()

    // Handle non-200 responses
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    // Parse response
    var result ShippingChargeResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("decoding response: %w", err)
    }

    if result.ErrorMessage != "" {
        return nil, fmt.Errorf("API error: %s", result.ErrorMessage)
    }

    return &result, nil
}