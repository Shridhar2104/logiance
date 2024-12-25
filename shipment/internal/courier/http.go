package courier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClient struct {
    client    *http.Client
    baseURL   string
    headers   map[string]string
}

func NewHTTPClient(baseURL string, timeout time.Duration) *HTTPClient {
    return &HTTPClient{
        client: &http.Client{
            Timeout: timeout,
        },
        baseURL: baseURL,
        headers: make(map[string]string),
    }
}

func (c *HTTPClient) SetHeader(key, value string) {
    c.headers[key] = value
}

// courier/http_client.go
func (c *HTTPClient) DoRaw(ctx context.Context, method, path string, headers map[string]string, body []byte) ([]byte, error) {
    url := c.baseURL + path
    req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
    if err != nil {
        return nil, fmt.Errorf("creating request failed: %w", err)
    }

    // Add headers
    for key, value := range headers {
        req.Header.Set(key, value)
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    // Read response body
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("reading response failed: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
    }

    return respBody, nil
}
func (c *HTTPClient) Do(ctx context.Context, method, path string, body interface{}, response interface{}) error {
    var reqBody []byte
    var err error

    if body != nil {
        reqBody, err = json.Marshal(body)
        if err != nil {
            return fmt.Errorf("failed to marshal request body: %w", err)
        }
    }

    req, err := http.NewRequestWithContext(ctx, method, 
        fmt.Sprintf("%s%s", c.baseURL, path), bytes.NewBuffer(reqBody))
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    // Set default headers
    req.Header.Set("Content-Type", "application/json")
    
    // Set custom headers
    for key, value := range c.headers {
        req.Header.Set(key, value)
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        var errorResp struct {
            Message string `json:"message"`
        }
        if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
            return fmt.Errorf("request failed with status %d", resp.StatusCode)
        }
        return fmt.Errorf("request failed: %s", errorResp.Message)
    }

    if response != nil {
        if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
            return fmt.Errorf("failed to decode response: %w", err)
        }
    }

    return nil
}