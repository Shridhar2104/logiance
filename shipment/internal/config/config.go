// shipment/internal/config/config.go
package config

type Config struct {
    ServerAddress string
    DelhiveryConfig DelhiveryConfig
}

type DelhiveryConfig struct {
    BaseURL     string
    APIKey      string
    RateLimit   int           // 40 requests per minute
    Environment string        // "staging" or "production"
}

func Load() *Config {
    // In production, this should load from environment variables
    return &Config{
        ServerAddress: ":50051",
        DelhiveryConfig: DelhiveryConfig{
            BaseURL:     "https://track.delhivery.com", // or staging URL
            APIKey:      "your-api-key",
            RateLimit:   40,
            Environment: "production",
        },
    }
}