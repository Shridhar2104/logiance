package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
    GRPCPort      string
    Couriers map[string]*CourierConfig
}

type CourierConfig struct {
    BaseURL    string
    RateLimit  int
    ApiKey     string
    ApiSecret  string
}

// internal/config/config.go
func NewConfig() *Config {
    return &Config{
        GRPCPort: ":50052",
        Couriers: map[string]*CourierConfig{
            "XPRESSBEES": {
                BaseURL:   getEnvWithDefault("XPRESSBEES_BASE_URL", "https://shipment.xpressbees.com/api"),
                RateLimit: getEnvAsIntWithDefault("XPRESSBEES_RATE_LIMIT", 35),
                ApiKey:    os.Getenv("XPRESSBEES_EMAIL"),    // Don't provide default value for credentials
                ApiSecret: os.Getenv("XPRESSBEES_PASSWORD"), // Don't provide default value for credentials
            },
            "DELHIVERY": {
                BaseURL:   getEnvWithDefault("DELHIVERY_BASE_URL", "https://staging-express.delhivery.com"),
                RateLimit: getEnvAsIntWithDefault("DELHIVERY_RATE_LIMIT", 40),
                ApiKey:    os.Getenv("DELHIVERY_API_KEY"),
            },
            "BLUEDART": {
                BaseURL:   getEnvWithDefault("BLUEDART_BASE_URL", "https://api.bluedart.com"),
                RateLimit: getEnvAsIntWithDefault("BLUEDART_RATE_LIMIT", 30),
                ApiKey:    os.Getenv("BLUEDART_API_KEY"),
            },
        },
    }
}

// Add validation function
func (c *Config) Validate() error {
    if c.Couriers["XPRESSBEES"].ApiKey == "" {
        return fmt.Errorf("XPRESSBEES_EMAIL environment variable is required")
    }
    if c.Couriers["XPRESSBEES"].ApiSecret == "" {
        return fmt.Errorf("XPRESSBEES_PASSWORD environment variable is required")
    }
    return nil
}

func (c *Config) GetCourierConfig(code string) *CourierConfig {
    if cfg, exists := c.Couriers[code]; exists {
        return cfg
    }
    return nil
}


func (c *Config) GetRateLimit(courierCode string) int {
    if config, exists := c.Couriers[courierCode]; exists {
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