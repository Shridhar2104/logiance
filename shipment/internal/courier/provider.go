// internal/courier/provider.go
package courier

import (
    "context"
    pb "github.com/Shridhar2104/logilo/shipment/proto"
)

// Provider defines the interface that all courier services must implement
type Provider interface {
    CalculateRate(ctx context.Context, req *pb.RateRequest) (*pb.CourierRate, error)
    IsAvailable(ctx context.Context, originPin, destPin int32) (bool, error)
    GetProviderInfo() *pb.CourierInfo
}
