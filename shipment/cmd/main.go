// cmd/main.go
package main

import (
    "log"
    "net"

    "github.com/Shridhar2104/logilo/shipment/internal/config"
    "github.com/Shridhar2104/logilo/shipment/internal/repository"
    "github.com/Shridhar2104/logilo/shipment/internal/service"
    "github.com/Shridhar2104/logilo/shipment/proto"
    "google.golang.org/grpc"
)

func main() {
    // Load configuration
    cfg := config.Load()

    // Initialize repository
    repo := repository.NewShipmentRepository(cfg)

    // Initialize service
    svc := service.NewShipmentService(repo)

    // Initialize gRPC server
    lis, err := net.Listen("tcp", cfg.ServerAddress)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    proto.RegisterShipmentServiceServer(grpcServer, svc)

    log.Printf("Starting gRPC server on %s", cfg.ServerAddress)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}