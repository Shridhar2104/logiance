package config

type Config struct {
    GRPCPort     string
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
                BaseURL:   "https://staging-express.delhivery.com",
                RateLimit: 40,
                ApiKey:    "your-delhivery-api-key",
            },
            "BLUEDART": {
                BaseURL:   "https://api.bluedart.com",
                RateLimit: 30,
                ApiKey:    "your-bluedart-api-key",
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