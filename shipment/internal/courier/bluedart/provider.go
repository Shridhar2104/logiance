// internal/courier/bluedart/provider.go
package bluedart

import (
	"context"
	"net/http"
	"time"

	"github.com/Shridhar2104/logilo/shipment/internal/config"
	"github.com/Shridhar2104/logilo/shipment/internal/courier"
	pb "github.com/Shridhar2104/logilo/shipment/proto"
)

type BlueDartProvider struct {
    cfg        *config.Config
    httpClient *http.Client
}

func NewProvider(cfg *config.Config) courier.Provider {
    return &BlueDartProvider{
        cfg: cfg,
        httpClient: &http.Client{
            Timeout: time.Second * 10,
        },
    }
}

func (d *BlueDartProvider) CalculateRate(ctx context.Context, req *pb.RateRequest) (*pb.CourierRate, error) {
    // Implementation similar to previous example but returns CourierRate
	return nil, nil
}


func (d *BlueDartProvider) IsAvailable(ctx context.Context, originPin, destPin int32) (bool, error) {
    // Implementation for checking serviceability
	return true, nil
}

func (d *BlueDartProvider) GetProviderInfo() *pb.CourierInfo {
    return &pb.CourierInfo{
        CourierCode: "DELHIVERY",
        CourierName: "Delhivery",
        SupportedServices: []string{"express", "surface"},
    }
}