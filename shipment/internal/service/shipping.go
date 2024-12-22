// Implement Provider interface methods...

// internal/service/shipping.go
package service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/Shridhar2104/logilo/shipment/internal/config"
	"github.com/Shridhar2104/logilo/shipment/internal/courier"
	"github.com/Shridhar2104/logilo/shipment/internal/courier/bluedart"
	"github.com/Shridhar2104/logilo/shipment/internal/courier/delhivery"
	"github.com/Shridhar2104/logilo/shipment/internal/ratelimit"
	pb "github.com/Shridhar2104/logilo/shipment/proto"
)

type ShipmentService struct {
    pb.UnimplementedShipmentServiceServer
    cfg          *config.Config
    rateLimiters map[string]*ratelimit.RateLimiter
    providers    map[string]courier.Provider
}

func NewShipmentService(cfg *config.Config) *ShipmentService {
    s := &ShipmentService{
        cfg:          cfg,
        rateLimiters: make(map[string]*ratelimit.RateLimiter),
        providers:    make(map[string]courier.Provider),
    }

    // Register providers
    s.registerProvider("DELHIVERY", delhivery.NewProvider(cfg))
    s.registerProvider("BLUEDART", bluedart.NewProvider(cfg))

    return s
}

func (s *ShipmentService) registerProvider(code string, provider courier.Provider) {
    s.providers[code] = provider
    s.rateLimiters[code] = ratelimit.NewRateLimiter(s.cfg.GetRateLimit(code))
}

func (s *ShipmentService) CalculateRates(ctx context.Context, req *pb.RateRequest) (*pb.MultiRateResponse, error) {
    var wg sync.WaitGroup
    ratesChan := make(chan *pb.CourierRate, len(s.providers))
    errorsChan := make(chan error, len(s.providers))

    // If specific couriers requested, use only those
    providers := s.providers
    if len(req.CourierCodes) > 0 {
        providers = make(map[string]courier.Provider)
        for _, code := range req.CourierCodes {
            if provider, exists := s.providers[code]; exists {
                providers[code] = provider
            }
        }
    }

    // Calculate rates from all providers concurrently
    for code, provider := range providers {
        wg.Add(1)
        go func(code string, provider courier.Provider) {
            defer wg.Done()

            if !s.rateLimiters[code].Allow() {
                errorsChan <- fmt.Errorf("rate limit exceeded for %s", code)
                return
            }

            rate, err := provider.CalculateRate(ctx, req)
            if err != nil {
                errorsChan <- err
                return
            }
            ratesChan <- rate
        }(code, provider)
    }

    // Wait for all goroutines to complete
    wg.Wait()
    close(ratesChan)
    close(errorsChan)

    // Collect results
    var rates []*pb.CourierRate
    var errors []string

    for rate := range ratesChan {
        rates = append(rates, rate)
    }

    for err := range errorsChan {
        errors = append(errors, err.Error())
    }

    // If we got some rates but some failed, still return success with available rates
    if len(rates) > 0 {
        return &pb.MultiRateResponse{
            Success: true,
            Rates:   rates,
            Error:   strings.Join(errors, "; "),
        }, nil
    }

    // If all failed, return error
    if len(errors) > 0 {
        return &pb.MultiRateResponse{
            Success: false,
            Error:   strings.Join(errors, "; "),
        }, nil
    }

    return &pb.MultiRateResponse{
        Success: true,
        Rates:   rates,
    }, nil
}

func (s *ShipmentService) GetAvailableCouriers(ctx context.Context, req *pb.AvailabilityRequest) (*pb.CourierListResponse, error) {
    var wg sync.WaitGroup
    availableChan := make(chan *pb.CourierInfo, len(s.providers))
    
    for _, provider := range s.providers {
        wg.Add(1)
        go func(provider courier.Provider) {
            defer wg.Done()
            
            available, err := provider.IsAvailable(ctx, req.OriginPincode, req.DestinationPincode)
            if err != nil || !available {
                return
            }
            
            availableChan <- provider.GetProviderInfo()
        }(provider)
    }

    wg.Wait()
    close(availableChan)

    var availableCouriers []*pb.CourierInfo
    for courier := range availableChan {
        availableCouriers = append(availableCouriers, courier)
    }

    return &pb.CourierListResponse{
        Success: true,
        AvailableCouriers: availableCouriers,
    }, nil
}