package config

import (
	"os"
	"strconv"
)

type Config struct {
    GRPCPort      string
    CourierConfig map[string]*CourierConfig
}

type CourierConfig struct {
    BaseURL    string
    RateLimit  int
    ApiKey     string
    ApiSecret  string
}

func NewConfig() *Config {
    return &Config{
        GRPCPort: ":50051",
        CourierConfig: map[string]*CourierConfig{
            "DELHIVERY": {
                BaseURL:   getEnvWithDefault("DELHIVERY_BASE_URL", "https://staging-express.delhivery.com"),
                RateLimit: getEnvAsIntWithDefault("DELHIVERY_RATE_LIMIT", 40),
                ApiKey:    getEnvWithDefault("DELHIVERY_API_KEY", ""),
            },
            "BLUEDART": {
                BaseURL:   getEnvWithDefault("BLUEDART_BASE_URL", "https://api.bluedart.com"),
                RateLimit: getEnvAsIntWithDefault("BLUEDART_RATE_LIMIT", 30),
                ApiKey:    getEnvWithDefault("BLUEDART_API_KEY", ""),
            },
        },
    }
}

func (c *Config) GetRateLimit(courierCode string) int {
    if config, exists := c.CourierConfig[courierCode]; exists {
        return config.RateLimit
    }
    return 30 // default rate limit
}

func getEnvWithDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsIntWithDefault(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}